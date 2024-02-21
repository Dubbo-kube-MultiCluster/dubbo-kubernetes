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
	core_runtime "github.com/apache/dubbo-kubernetes/pkg/core/runtime"
	"time"
)

var log = core.Log.WithName("clusterID")

// clusterIDReader tries to read cluster ID and sets it in the runtime. Cluster ID does not change during CP lifecycle
// therefore once cluster ID is read and set, the component exits.
// In single-zone setup, followers are waiting until leader creates a cluster ID
// In multi-zone setup, the global followers and all zones waits until global leader creates a cluster ID
type clusterIDReader struct {
	rt core_runtime.Runtime
}

func (c *clusterIDReader) Start(stop <-chan struct{}) error {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			clusterID, err := c.read(context.Background())
			if err != nil {
				log.Error(err, "could not read cluster ID") // just log, do not exit to retry operation
			}
			if clusterID != "" {
				log.Info("setting cluster ID", "clusterID", clusterID)
				c.rt.SetClusterId(clusterID)
				return nil
			}
		case <-stop:
			return nil
		}
	}
}

func (c *clusterIDReader) NeedLeaderElection() bool {
	return false
}

func (c *clusterIDReader) read(ctx context.Context) (string, error) {
	resource := config_model.NewConfigResource()
	err := c.rt.ConfigManager().Get(ctx, resource, store.GetByKey(config_manager.ClusterIdConfigKey, core_model.NoMesh))
	if err != nil {
		if store.IsResourceNotFound(err) {
			return "", nil
		}
		return "", err
	}
	return resource.Spec.Config, nil
}
