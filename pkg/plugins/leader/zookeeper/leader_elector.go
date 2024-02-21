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

package zookeeper

import (
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	"github.com/apex/log"
	"github.com/dubbogo/go-zookeeper/zk"
)

type zookeeperLeaderElector struct {
	alwaysLeader bool
	lockClient   *zk.Lock
	callbacks    []component.LeaderCallbacks
}

func NewZookeeperLeaderElector(lockClient *zk.Lock) component.LeaderElector {
	return &zookeeperLeaderElector{
		lockClient: lockClient,
	}
}

func (n *zookeeperLeaderElector) AddCallbacks(callbacks component.LeaderCallbacks) {
	n.callbacks = append(n.callbacks, callbacks)
}

func (n *zookeeperLeaderElector) IsLeader() bool {
	return n.alwaysLeader
}

func (n *zookeeperLeaderElector) Start(stop <-chan struct{}) {
	for {
		log.Info("waiting for lock")
		n.lockClient.Lock()
	}
	if n.alwaysLeader {
		for _, callback := range n.callbacks {
			callback.OnStartedLeading()
		}
	}
}
