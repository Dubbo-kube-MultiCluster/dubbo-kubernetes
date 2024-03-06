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

package universal

import (
	"dubbo.apache.org/dubbo-go/v3/registry"
	gxset "github.com/dubbogo/gost/container/set"
	"sync"
)

type ServiceMappingChangedListenerImpl struct {
	oldServiceNames *gxset.HashSet
	listener        registry.NotifyListener
	interfaceKey    string

	mux           sync.Mutex
	delSDRegistry registry.ServiceDiscovery
}

func NewMappingListener(oldServiceNames *gxset.HashSet, listener registry.NotifyListener) *ServiceMappingChangedListenerImpl {
	return &ServiceMappingChangedListenerImpl{
		listener:        listener,
		oldServiceNames: oldServiceNames,
	}
}

func (lstn *ServiceMappingChangedListenerImpl) updateListener(interfaceKey string, apps *gxset.HashSet) error {
	return nil
}

// Stop on ServiceMappingChangedEvent the service mapping change event
func (lstn *ServiceMappingChangedListenerImpl) Stop() {}
