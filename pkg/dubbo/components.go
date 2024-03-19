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

package dubbo

import (
	"time"
)

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	core_mesh "github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/mesh"
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	core_runtime "github.com/apache/dubbo-kubernetes/pkg/core/runtime"
	"github.com/apache/dubbo-kubernetes/pkg/dubbo/pusher"
	"github.com/apache/dubbo-kubernetes/pkg/dubbo/servicemapping"
)

var log = core.Log.WithName("dubbo")

func Setup(rt core_runtime.Runtime) error {
	cfg := rt.Config().ServiceNameMapping

	snpPusher := pusher.NewPusher(rt.ResourceManager(), rt.EventBus(), func() *time.Ticker {
		// todo: should configured by config in the future
		return time.NewTicker(time.Minute * 10)
	}, []core_model.ResourceType{core_mesh.MappingType})

	// register ServiceNameMappingService
	serviceMapping := servicemapping.NewSnpServer(
		rt.AppContext(),
		cfg.ServiceMapping,
		snpPusher,
		rt.ResourceManager(),
		rt.Transactions(),
		rt.Config().Multizone.Zone.Name,
	)
	mesh_proto.RegisterServiceNameMappingServiceServer(rt.DpServer().GrpcServer(), serviceMapping)

	return rt.Add(snpPusher, serviceMapping)
}
