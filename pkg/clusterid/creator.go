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

package clusterid

import (
	"context"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	config_manager "github.com/apache/dubbo-kubernetes/pkg/core/config/manager"
	config_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/system"
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/pkg/errors"
)

type clusterIDCreator struct {
	configManager config_manager.ConfigManager
}

func (c *clusterIDCreator) Start(_ <-chan struct{}) error {
	return c.create()
}

func (c *clusterIDCreator) NeedLeaderElection() bool {
	return true
}

func (c *clusterIDCreator) create() error {
	ctx := context.Background()
	resource := config_model.NewConfigResource()
	err := c.configManager.Get(ctx, resource, store.GetByKey(config_manager.ClusterIdConfigKey, core_model.NoMesh))
	if err != nil {
		if !store.IsResourceNotFound(err) {
			return err
		}
		resource.Spec.Config = core.NewUUID()
		log.Info("creating cluster ID", "clusterID", resource.Spec.Config)
		if err := c.configManager.Create(ctx, resource, store.CreateByKey(config_manager.ClusterIdConfigKey, core_model.NoMesh)); err != nil {
			return errors.Wrap(err, "could not create config")
		}
	}

	return nil
}
