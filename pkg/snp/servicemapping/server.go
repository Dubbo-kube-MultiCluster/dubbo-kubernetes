/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicemapping

import (
	"context"
	"io"
	"time"
)

import (
	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	"github.com/apache/dubbo-kubernetes/pkg/config/snp"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	core_mesh "github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/mesh"
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	core_store "github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	"github.com/apache/dubbo-kubernetes/pkg/events"
	"github.com/apache/dubbo-kubernetes/pkg/snp/client"
)

var log = core.Log.WithName("snp").WithName("server").WithName("service-name-mapping")

const queueSize = 100

var _ component.Component = &SnpServer{}

type SnpServer struct {
	mesh_proto.ServiceNameMappingServiceServer

	locallyZone string
	config      snp.ServiceMapping
	queue       chan *RegisterRequest

	ctx           context.Context
	eventBus      events.EventBus
	resourceStore core_store.ResourceStore
	transactions  core_store.Transactions
}

func (s *SnpServer) Start(stop <-chan struct{}) error {
	// we start debounce to prevent too many MappingRegisterRequests, we aggregate mapping register information
	go s.debounce(stop, s.register)

	return nil
}

func (s *SnpServer) NeedLeaderElection() bool {
	return false
}

func NewSnpServer(
	ctx context.Context,
	config snp.ServiceMapping,
	resourceStore core_store.ResourceStore,
	transactions core_store.Transactions,
	eventBus events.EventBus,
	locallyZone string,
) *SnpServer {
	return &SnpServer{
		locallyZone:   locallyZone,
		config:        config,
		queue:         make(chan *RegisterRequest, queueSize),
		ctx:           ctx,
		eventBus:      eventBus,
		resourceStore: resourceStore,
		transactions:  transactions,
	}
}

func (s *SnpServer) MappingRegister(ctx context.Context, req *mesh_proto.MappingRegisterRequest) (*mesh_proto.MappingRegisterResponse, error) {
	mesh := core_model.DefaultMesh // todo: mesh
	interfaces := req.GetInterfaceNames()
	applicationName := req.GetApplicationName()

	registerReq := &RegisterRequest{ConfigsUpdated: map[core_model.ResourceKey]map[string]struct{}{}}
	for _, interfaceName := range interfaces {
		key := core_model.ResourceKey{
			Mesh: mesh,
			Name: interfaceName,
		}
		if _, ok := registerReq.ConfigsUpdated[key]; !ok {
			registerReq.ConfigsUpdated[key] = make(map[string]struct{})
		}
		registerReq.ConfigsUpdated[key][applicationName] = struct{}{}
	}

	// push into queue to debounce, register Mapping Resource
	s.queue <- registerReq

	return &mesh_proto.MappingRegisterResponse{
		Success: true,
		Message: "success",
	}, nil
}

func (s *SnpServer) MappingSync(stream mesh_proto.ServiceNameMappingService_MappingSyncServer) error {
	mesh := core_model.DefaultMesh // todo: mesh
	errChan := make(chan error)

	mappingSyncStream := client.NewMappingSyncStream(stream)
	// MappingSyncClient is to handle MappingSyncRequest from data plane
	mappingSyncClient := client.NewMappingSyncClient(log.WithName("client"), mappingSyncStream, &client.Callbacks{
		OnResourcesReceived: func(request *mesh_proto.MappingSyncRequest) (core_mesh.MappingResourceList, error) {
			// List Mapping Resources which client subscribed.
			resourceKeys := make([]core_model.ResourceKey, 0, len(mappingSyncStream.SubscribedInterfaceNames()))
			for _, interfaceName := range mappingSyncStream.SubscribedInterfaceNames() {
				resourceKeys = append(resourceKeys, core_model.ResourceKey{
					Mesh: mesh,
					Name: interfaceName,
				})
			}

			var mappingList core_mesh.MappingResourceList
			err := s.resourceStore.List(stream.Context(), &mappingList, core_store.ListByResourceKeys(resourceKeys))
			if err != nil {
				return core_mesh.MappingResourceList{}, err
			}

			return mappingList, err
		},
	})
	go func() {
		// Handle requests from client
		err := mappingSyncClient.HandleReceive()
		if err != nil {
			log.Info("MappingSyncClient finished gracefully")
			errChan <- nil
		}

		log.Error(err, "MappingSyncClient finished with an error")
		errChan <- errors.Wrap(err, "MappingSyncClient finished with an error")
	}()

	// subscribe Mapping changed Event
	mappingChanged := s.eventBus.Subscribe(func(e events.Event) bool {
		resourceChangedEvent, ok := e.(events.ResourceChangedEvent)
		if ok {
			for _, interfaceName := range mappingSyncStream.SubscribedInterfaceNames() {
				if resourceChangedEvent.Key.Name == interfaceName {
					return true
				}
			}
		}

		return false
	})
	defer mappingChanged.Close()

	for {
		select {
		case changedEvent := <-mappingChanged.Recv():
			// when Mapping Changed
			mappingChangedEvent, ok := changedEvent.(events.ResourceChangedEvent)
			if !ok {
				continue
			}

			// get Mapping Resource
			mapping := core_mesh.NewMappingResource()
			err := s.resourceStore.Get(stream.Context(), mapping, core_store.GetBy(
				mappingChangedEvent.Key,
			))
			if err != nil {
				log.Error(err, "MappingSync get mapping resource failed")
				if err := mappingSyncStream.NACK([]string{mapping.Spec.InterfaceName}, err.Error()); err != nil {
					if err == io.EOF {
						errChan <- nil
						continue
					}
					errChan <- errors.Wrap(err, "failed to NACK a MappingSyncRequest")
					continue
				}
				continue
			}

			// sync Mapping Resource to client
			var mappingList core_mesh.MappingResourceList
			_ = mappingList.AddItem(mapping)
			err = mappingSyncStream.ACK(mappingList)
			if err != nil {
				log.Error(err, "MappingSync ack mapping resource failed")
				if err := mappingSyncStream.NACK([]string{mapping.Spec.InterfaceName}, err.Error()); err != nil {
					if err == io.EOF {
						errChan <- nil
						continue
					}
					errChan <- errors.Wrap(err, "failed to ACK a MappingSyncRequest")
					continue
				}
				continue
			}

		case err := <-errChan:
			if err == nil {
				log.Info("MappingSync finished gracefully")
				return nil
			}

			log.Error(err, "MappingSync finished with an error")
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (s *SnpServer) debounce(stopCh <-chan struct{}, pushFn func(m *RegisterRequest)) {
	ch := s.queue
	var timeChan <-chan time.Time
	var startDebounce time.Time
	var lastConfigUpdateTime time.Time

	pushCounter := 0
	debouncedEvents := 0

	var req *RegisterRequest

	free := true
	freeCh := make(chan struct{}, 1)

	push := func(req *RegisterRequest) {
		pushFn(req)
		freeCh <- struct{}{}
	}

	pushWorker := func() {
		eventDelay := time.Since(startDebounce)
		quietTime := time.Since(lastConfigUpdateTime)
		if eventDelay >= s.config.Debounce.Max || quietTime >= s.config.Debounce.After {
			if req != nil {
				pushCounter++

				if req.ConfigsUpdated != nil {
					log.Info("debounce stable[%d] %d for config %s: %v since last change, %v since last push",
						pushCounter, debouncedEvents, configsUpdated(req),
						quietTime, eventDelay)
				}
				free = false
				go push(req)
				req = nil
				debouncedEvents = 0
			}
		} else {
			timeChan = time.After(s.config.Debounce.After - quietTime)
		}
	}

	for {
		select {
		case <-freeCh:
			free = true
			pushWorker()
		case r := <-ch:
			if !s.config.Debounce.Enable {
				go push(r)
				req = nil
				continue
			}

			lastConfigUpdateTime = time.Now()
			if debouncedEvents == 0 {
				timeChan = time.After(200 * time.Millisecond)
				startDebounce = lastConfigUpdateTime
			}
			debouncedEvents++

			req = req.Merge(r)
		case <-timeChan:
			if free {
				pushWorker()
			}
		case <-stopCh:
			return
		}
	}
}

func (s *SnpServer) register(req *RegisterRequest) {
	for key, m := range req.ConfigsUpdated {
		var appNames []string
		for app := range m {
			appNames = append(appNames, app)
		}
		for i := 0; i < 3; i++ {
			if err := s.tryRegister(key.Mesh, key.Name, appNames); err != nil {
				log.Error(err, "register failed", "key", key)
			} else {
				break
			}
		}
	}
}

func (s *SnpServer) tryRegister(mesh, interfaceName string, newApps []string) error {
	err := core_store.InTx(s.ctx, s.transactions, func(ctx context.Context) error {
		key := core_model.ResourceKey{
			Mesh: mesh,
			Name: interfaceName,
		}

		// get Mapping Resource first,
		// if Mapping is not found, create it,
		// else update it.
		mapping := core_mesh.NewMappingResource()
		err := s.resourceStore.Get(s.ctx, mapping, core_store.GetBy(key))
		if err != nil && !core_store.IsResourceNotFound(err) {
			log.Error(err, "get Mapping Resource")
			return err
		}

		if core_store.IsResourceNotFound(err) {
			// create if not found
			mapping.Spec = &mesh_proto.Mapping{
				Zone:             s.locallyZone,
				InterfaceName:    interfaceName,
				ApplicationNames: newApps,
			}
			err = s.resourceStore.Create(s.ctx, mapping, core_store.CreateBy(key), core_store.CreatedAt(time.Now()))
			if err != nil {
				log.Error(err, "create Mapping Resource failed")
				return err
			}

			log.Info("create Mapping Resource success", "key", key, "applicationNames", newApps)
			return nil
		} else {
			// if found, update it
			previousLen := len(mapping.Spec.ApplicationNames)
			previousAppNames := make(map[string]struct{}, previousLen)
			for _, name := range mapping.Spec.ApplicationNames {
				previousAppNames[name] = struct{}{}
			}
			for _, newApp := range newApps {
				previousAppNames[newApp] = struct{}{}
			}
			if len(previousAppNames) == previousLen {
				log.Info("Mapping not need to register", "interfaceName", interfaceName, "applicationNames", newApps)
				return nil
			}

			mergedApps := make([]string, 0, len(previousAppNames))
			for name := range previousAppNames {
				mergedApps = append(mergedApps, name)
			}
			mapping.Spec = &mesh_proto.Mapping{
				Zone:             s.locallyZone,
				InterfaceName:    interfaceName,
				ApplicationNames: mergedApps,
			}

			err = s.resourceStore.Update(s.ctx, mapping, core_store.ModifiedAt(time.Now()))
			if err != nil {
				log.Error(err, "update Mapping Resource failed")
				return err
			}

			log.Info("update Mapping Resource success", "key", key, "applicationNames", newApps)
			return nil
		}
	})
	if err != nil {
		log.Error(err, "transactions failed")
		return err
	}

	return nil
}
