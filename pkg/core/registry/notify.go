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

package registry

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/manager"
	"sync"
)

type NotifyListener struct {
	manager.ResourceManager
	dataplaneCache *sync.Map
}

func NewNotifyListener(manager manager.ResourceManager, cache *sync.Map) *NotifyListener {
	return &NotifyListener{
		manager,
		cache,
	}
}

func (l *NotifyListener) Notify(event *dubboRegistry.ServiceEvent) {
	switch event.Action {
	case remoting.EventTypeAdd, remoting.EventTypeUpdate:
		if err := l.createOrUpdateDataplane(context.Background(), event.Service); err != nil {
			return
		}
	case remoting.EventTypeDel:
		if err := l.deleteDataplane(context.Background(), event.Service); err != nil {
			return
		}
	}
}

func (l *NotifyListener) NotifyAll(events []*dubboRegistry.ServiceEvent, f func()) {
	for _, event := range events {
		l.Notify(event)
	}
}

func (l *NotifyListener) deleteDataplane(ctx context.Context, url *common.URL) error {
	return nil
}

func (l *NotifyListener) createOrUpdateDataplane(ctx context.Context, url *common.URL) error {
	return nil
}

func InboundInterfacesFor(ctx context.Context, url *common.URL) ([]*mesh_proto.Dataplane_Networking_Inbound, error) {
	var ifaces []*mesh_proto.Dataplane_Networking_Inbound

	return ifaces, nil
}

func OutboundInterfacesFor(ctx context.Context, url *common.URL) ([]*mesh_proto.Dataplane_Networking_Outbound, error) {
	var outbounds []*mesh_proto.Dataplane_Networking_Outbound

	return outbounds, nil
}
