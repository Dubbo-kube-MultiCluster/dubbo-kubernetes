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

package traditional

import (
	"context"
	"fmt"
	"github.com/apache/dubbo-kubernetes/pkg/core/logger"
	"github.com/dubbogo/gost/encoding/yaml"
	"sync"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/config_center"
	dubbo_identifier "dubbo.apache.org/dubbo-go/v3/metadata/identifier"
	"dubbo.apache.org/dubbo-go/v3/metadata/report"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/pkg/errors"

	"golang.org/x/exp/maps"
)

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	"github.com/apache/dubbo-kubernetes/pkg/core/consts"
	"github.com/apache/dubbo-kubernetes/pkg/core/governance"
	"github.com/apache/dubbo-kubernetes/pkg/core/reg_client"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/mesh"
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/plugins/util/ccache"
)

const (
	dubboGroup    = "dubbo"
	mappingGroup  = "mapping"
	cpGroup       = "dubbo-cp"
	pathSeparator = "/"
)

type traditionalStore struct {
	configCenter   config_center.DynamicConfiguration
	metadataReport report.MetadataReport
	registryCenter dubboRegistry.Registry
	governance     governance.GovernanceConfig
	dCache         *sync.Map
	regClient      reg_client.RegClient
}

func NewStore(
	configCenter config_center.DynamicConfiguration,
	metadataReport report.MetadataReport,
	registryCenter dubboRegistry.Registry,
	governance governance.GovernanceConfig,
	dCache *sync.Map,
	regClient reg_client.RegClient,
) store.ResourceStore {
	return &traditionalStore{
		configCenter:   configCenter,
		metadataReport: metadataReport,
		registryCenter: registryCenter,
		governance:     governance,
		dCache:         dCache,
		regClient:      regClient,
	}
}

func (t *traditionalStore) Create(_ context.Context, resource core_model.Resource, fs ...store.CreateOptionsFunc) error {
	var err error
	opts := store.NewCreateOptions(fs...)

	switch resource.Descriptor().Name {
	case mesh.MappingType:
		spec := resource.GetSpec()
		mapping := spec.(*mesh_proto.Mapping)
		appNames := mapping.ApplicationNames
		serviceInterface := mapping.InterfaceName
		for _, app := range appNames {
			err = t.metadataReport.RegisterServiceAppMapping(serviceInterface, mappingGroup, app)
			if err != nil {
				return err
			}
		}
	case mesh.MetaDataType:
		spec := resource.GetSpec()
		metadata := spec.(*mesh_proto.MetaData)
		identifier := &dubbo_identifier.SubscriberMetadataIdentifier{
			Revision: metadata.GetRevision(),
			BaseApplicationMetadataIdentifier: dubbo_identifier.BaseApplicationMetadataIdentifier{
				Application: metadata.GetApp(),
				Group:       dubboGroup,
			},
		}
		services := map[string]*common.ServiceInfo{}
		// 把metadata赋值到services中
		for key, serviceInfo := range metadata.GetServices() {
			services[key] = &common.ServiceInfo{
				Name:     serviceInfo.GetName(),
				Group:    serviceInfo.GetGroup(),
				Version:  serviceInfo.GetVersion(),
				Protocol: serviceInfo.GetProtocol(),
				Path:     serviceInfo.GetPath(),
				Params:   serviceInfo.GetParams(),
			}
		}
		info := &common.MetadataInfo{
			App:      metadata.GetApp(),
			Revision: metadata.GetRevision(),
			Services: services,
		}
		err = t.metadataReport.PublishAppMetadata(identifier, info)
		if err != nil {
			return err
		}
	case mesh.DataplaneType:
		// Dataplane无法Create, 只能Get和List
	case mesh.TagRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetRoutePath(id, consts.TagRoute)
		bytes, err := core_model.ToYAML(resource.GetSpec())
		if err != nil {
			return err
		}

		err = t.governance.SetConfig(path, string(bytes))
		if err != nil {
			return err
		}

	case mesh.ConditionRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetRoutePath(id, consts.ConditionRoute)

		bytes, err := core_model.ToYAML(resource.GetSpec())
		if err != nil {
			return err
		}

		err = t.governance.SetConfig(path, string(bytes))
		if err != nil {
			return err
		}
	case mesh.DynamicConfigType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetOverridePath(id)
		bytes, err := core_model.ToYAML(resource.GetSpec())
		if err != nil {
			return err
		}

		err = t.governance.SetConfig(path, string(bytes))
		if err != nil {
			return err
		}
	default:
		bytes, err := core_model.ToYAML(resource.GetSpec())
		if err != nil {
			return err
		}

		path := GenerateCpGroupPath(string(resource.Descriptor().Name), opts.Name)
		// 使用RegClient
		err = t.regClient.SetContent(path, bytes)
		if err != nil {
			return err
		}
	}

	resource.SetMeta(&resourceMetaObject{
		Name:             opts.Name,
		Mesh:             opts.Mesh,
		CreationTime:     opts.CreationTime,
		ModificationTime: opts.CreationTime,
		Labels:           maps.Clone(opts.Labels),
	})
	return nil
}

func (t *traditionalStore) Update(ctx context.Context, resource core_model.Resource, fs ...store.UpdateOptionsFunc) error {
	return nil
}

func (t *traditionalStore) Delete(ctx context.Context, resource core_model.Resource, fs ...store.DeleteOptionsFunc) error {
	opts := store.NewDeleteOptions(fs...)

	switch resource.Descriptor().Name {
	case mesh.DataplaneType:
		t.dCache.Delete(ccache.GenerateDCacheKey(opts.Name, opts.Mesh))
	case mesh.TagRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		key := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetOverridePath(key)
		err := t.governance.DeleteConfig(path)
		if err != nil {
			return err
		}
	case mesh.ConditionRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		key := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetRoutePath(key, consts.ConditionRoute)
		err := t.governance.DeleteConfig(path)
		if err != nil {
			return err
		}
	case mesh.DynamicConfigType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		key := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetOverridePath(key)
		conf, err := t.governance.GetConfig(path)
		if err != nil {
			logger.Sugar().Error(err.Error())
			return err
		}
		if err := core_model.FromYAML([]byte(conf), resource.GetSpec()); err != nil {
			return err
		}
		override := resource.GetSpec().(*mesh_proto.DynamicConfig)
		if len(override.Configs) > 0 {
			newConfigs := make([]*mesh_proto.OverrideConfig, 0)
			for _, c := range override.Configs {
				if consts.Configs.Contains(c.Type) {
					newConfigs = append(newConfigs, c)
				}
			}
			if len(newConfigs) == 0 {
				err := t.governance.DeleteConfig(path)
				if err != nil {
					return err
				}
			} else {
				override.Configs = newConfigs
				if b, err := yaml.MarshalYML(override); err != nil {
					return err
				} else {
					err := t.governance.SetConfig(path, string(b))
					if err != nil {
						return err
					}
				}
			}
		} else {
			err := t.governance.DeleteConfig(path)
			if err != nil {
				return err
			}
		}
	case mesh.MappingType:
		// 无法删除
	case mesh.MetaDataType:
		// 无法删除
	default:
		path := GenerateCpGroupPath(string(resource.Descriptor().Name), opts.Name)
		err := t.regClient.DeleteContent(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *traditionalStore) Get(_ context.Context, resource core_model.Resource, fs ...store.GetOptionsFunc) error {
	opts := store.NewGetOptions(fs...)

	switch resource.Descriptor().Name {
	case mesh.DataplaneType:
		value, ok := c.dCache.Load(ccache.GenerateDCacheKey(opts.Name, opts.Mesh))
		if !ok {
			return nil
		}
		r := value.(core_model.Resource)
		resource.SetMeta(r.GetMeta())
		err := resource.SetSpec(r.GetSpec())
		if err != nil {
			return err
		}
	case mesh.TagRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetRoutePath(id, consts.TagRoute)
		cfg, err := c.governance.GetConfig(path)
		if err != nil {
			return err
		}
		if cfg != "" {
			if err := core_model.FromYAML([]byte(cfg), resource.GetSpec()); err != nil {
				return errors.Wrap(err, "failed to convert json to spec")
			}
		}
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
	case mesh.ConditionRouteType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetRoutePath(id, consts.ConditionRoute)
		cfg, err := c.governance.GetConfig(path)
		if err != nil {
			return err
		}
		if cfg != "" {
			if err := core_model.FromYAML([]byte(cfg), resource.GetSpec()); err != nil {
				return errors.Wrap(err, "failed to convert json to spec")
			}
		}
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
	case mesh.DynamicConfigType:
		labels := opts.Labels
		base := mesh_proto.Base{
			Application:    labels[mesh_proto.Application],
			Service:        labels[mesh_proto.Service],
			ID:             labels[mesh_proto.ID],
			ServiceVersion: labels[mesh_proto.ServiceVersion],
			ServiceGroup:   labels[mesh_proto.ServiceGroup],
		}
		id := mesh_proto.BuildServiceKey(base)
		path := mesh_proto.GetOverridePath(id)
		cfg, err := c.governance.GetConfig(path)
		if err != nil {
			return err
		}
		if cfg != "" {
			if err := core_model.FromYAML([]byte(cfg), resource.GetSpec()); err != nil {
				return errors.Wrap(err, "failed to convert json to spec")
			}
		}
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
	case mesh.MappingType:
		// Get通过Key获取, 不设置listener
		key := opts.Name
		set, err := c.metadataReport.GetServiceAppMapping(key, mappingGroup, nil)
		if err != nil {
			return err
		}
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
		mapping := resource.GetSpec().(*mesh_proto.Mapping)
		mapping.Zone = "default"
		mapping.InterfaceName = key
		var items []string
		for k := range set.Items {
			items = append(items, fmt.Sprintf("%v", k))
		}
		mapping.ApplicationNames = items
	case mesh.MetaDataType:
		labels := opts.Labels
		id := dubbo_identifier.NewSubscriberMetadataIdentifier(labels[mesh_proto.Application], labels[mesh_proto.Revision])
		appMetadata, err := c.metadataReport.GetAppMetadata(id)
		if err != nil {
			return err
		}
		metaData := resource.GetSpec().(*mesh_proto.MetaData)
		metaData.App = appMetadata.App
		metaData.Revision = appMetadata.Revision
		service := map[string]*mesh_proto.ServiceInfo{}
		for key, serviceInfo := range appMetadata.Services {
			service[key] = &mesh_proto.ServiceInfo{
				Name:     serviceInfo.Name,
				Group:    serviceInfo.Group,
				Version:  serviceInfo.Version,
				Protocol: serviceInfo.Protocol,
				Path:     serviceInfo.Path,
				Params:   serviceInfo.Params,
			}
		}
		metaData.Services = service
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
	default:
		path := GenerateCpGroupPath(string(resource.Descriptor().Name), opts.Name)
		value, err := c.regClient.GetContent(path)
		if err != nil {
			return err
		}
		if err := core_model.FromYAML(value, resource.GetSpec()); err != nil {
			return err
		}
		meta := &resourceMetaObject{
			Name: opts.Name,
			Mesh: opts.Mesh,
		}
		resource.SetMeta(meta)
	}
	return nil
}

func (c *traditionalStore) List(_ context.Context, resource core_model.ResourceList, fs ...store.ListOptionsFunc) error {
	return nil
}
