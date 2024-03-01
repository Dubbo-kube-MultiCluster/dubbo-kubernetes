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

package snp

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	core_runtime "github.com/apache/dubbo-kubernetes/pkg/core/runtime"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	snp_server "github.com/apache/dubbo-kubernetes/pkg/snp/server"
	"github.com/apache/dubbo-kubernetes/pkg/snp/servicemapping"
)

var log = core.Log.WithName("snp")

func Setup(rt core_runtime.Runtime) error {
	cfg := rt.Config().ServiceNameMapping

	snpComponent := component.ComponentFunc(func(stop <-chan struct{}) error {
		server := snp_server.NewServer(cfg.Server)

		// register ServiceNameMappingService
		serviceMapping := servicemapping.NewSnpServer(rt.AppContext(), cfg.ServiceMapping, rt.ResourceStore(), rt.Transactions(), rt.EventBus())
		mesh_proto.RegisterServiceNameMappingServiceServer(server.GrpcServer(), serviceMapping)

		return rt.Add(server, serviceMapping)
	})

	return rt.Add(component.NewResilientComponent(
		log,
		snpComponent,
	))
}
