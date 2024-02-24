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

package mysql

import (
	"context"
)

import (
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/events"
)

type mySqlStore struct{}

func NewStore() store.ResourceStore {
	return &mySqlStore{}
}

func (c *mySqlStore) SetEventWriter(writer events.Emitter) {
}

func (c *mySqlStore) Create(_ context.Context, resource core_model.Resource, fs ...store.CreateOptionsFunc) error {
	return nil
}

func (c *mySqlStore) Update(_ context.Context, resource core_model.Resource, fs ...store.UpdateOptionsFunc) error {
	return nil
}

func (c *mySqlStore) Delete(ctx context.Context, resource core_model.Resource, fs ...store.DeleteOptionsFunc) error {
	return nil
}

func (c *mySqlStore) Get(_ context.Context, resource core_model.Resource, fs ...store.GetOptionsFunc) error {
	return nil
}

func (c *mySqlStore) List(_ context.Context, resource core_model.ResourceList, fs ...store.ListOptionsFunc) error {
	return nil
}
