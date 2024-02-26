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

package leader

import (
	leader_zookeeper "github.com/apache/dubbo-kubernetes/pkg/plugins/leader/zookeeper"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/config/core/resources/store"
	core_runtime "github.com/apache/dubbo-kubernetes/pkg/core/runtime"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	common_mysql "github.com/apache/dubbo-kubernetes/pkg/plugins/common/mysql"
	common_zookeeper "github.com/apache/dubbo-kubernetes/pkg/plugins/common/zookeeper"
	leader_memory "github.com/apache/dubbo-kubernetes/pkg/plugins/leader/memory"
	leader_mysql "github.com/apache/dubbo-kubernetes/pkg/plugins/leader/mysql"
)

func NewLeaderElector(b *core_runtime.Builder) (component.LeaderElector, error) {
	switch b.Config().Store.Type {
	case store.MemoryStore:
		return leader_memory.NewAlwaysLeaderElector(), nil
	case store.ZookeeperStore:
		cfg := *b.Config().Store.Zookeeper
		connect, err := common_zookeeper.ConnectToZK(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "cloud not connect to zookeeper")
		}
		return leader_zookeeper.NewZookeeperLeaderElector(connect), nil
	case store.MyStore:
		cfg := *b.Config().Store.Mysql
		db, err := common_mysql.ConnectToDb(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "cloud not connect to mysql")
		}
		return leader_mysql.NewMysqlLeaderElector(db), nil

	// In case of Kubernetes, Leader Elector is embedded in a Kubernetes ComponentManager
	default:
		return nil, errors.Errorf("no election leader for storage of type %s", b.Config().Store.Type)
	}
}
