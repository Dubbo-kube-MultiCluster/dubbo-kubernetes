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

package metadata

import (
	"context"
)

import (
	core_manager "github.com/apache/dubbo-kubernetes/pkg/core/resources/manager"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	core_store "github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
)

type metadataManager struct {
	core_manager.ResourceManager
}

func NewMetadataManager() core_manager.ResourceManager {
	return nil
}

func (m *metadataManager) Create(ctx context.Context, r model.Resource, opts ...core_store.CreateOptionsFunc) error {
	return nil
}

func (m *metadataManager) Update(ctx context.Context, r model.Resource, opts ...core_store.UpdateOptionsFunc) error {
	return nil
}

func (m *metadataManager) Delete(ctx context.Context, r model.Resource, opts ...core_store.DeleteOptionsFunc) error {
	return nil
}

func (m *metadataManager) Get(ctx context.Context, r model.Resource, opts ...core_store.GetOptionsFunc) error {
	return nil
}

func (m *metadataManager) List(ctx context.Context, r model.ResourceList, opts ...core_store.ListOptionsFunc) error {
	return nil
}
