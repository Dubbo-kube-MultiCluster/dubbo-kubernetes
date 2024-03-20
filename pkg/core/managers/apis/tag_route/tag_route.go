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

package tag_route

import (
	"context"
)

import (
	core_manager "github.com/apache/dubbo-kubernetes/pkg/core/resources/manager"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	core_store "github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
)

type tagRouteManager struct {
	core_manager.ResourceManager
}

func NewTagRouteManager() core_manager.ResourceManager {
	return nil
}

func (m *tagRouteManager) Create(ctx context.Context, r model.Resource, opts ...core_store.CreateOptionsFunc) error {
	return nil
}

func (m *tagRouteManager) Update(ctx context.Context, r model.Resource, opts ...core_store.UpdateOptionsFunc) error {
	return nil
}

func (m *tagRouteManager) Delete(ctx context.Context, r model.Resource, opts ...core_store.DeleteOptionsFunc) error {
	return nil
}

func (m *tagRouteManager) Get(ctx context.Context, r model.Resource, opts ...core_store.GetOptionsFunc) error {
	return nil
}

func (m *tagRouteManager) List(ctx context.Context, r model.ResourceList, opts ...core_store.ListOptionsFunc) error {
	return nil
}
