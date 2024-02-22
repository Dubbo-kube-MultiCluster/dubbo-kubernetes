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
	"github.com/apache/dubbo-kubernetes/pkg/core"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	util_channels "github.com/apache/dubbo-kubernetes/pkg/util/channels"
	"github.com/dubbogo/go-zookeeper/zk"
	"time"
)

const dubboLockName = "/dubbo/cp-lock"
const backoffTime = 5 * time.Second

var log = core.Log.WithName("zookeeper-leader")

type zookeeperLeaderElector struct {
	alwaysLeader bool
	lockClient   *zk.Lock
	callbacks    []component.LeaderCallbacks
}

func NewZookeeperLeaderElector(connect *zk.Conn) component.LeaderElector {
	lockClient := zk.NewLock(connect, dubboLockName, zk.WorldACL(zk.PermAll))
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
	retries := 0
	for {
		log.Info("waiting for lock")
		err := n.lockClient.Lock()
		if err != nil {
			if retries >= 3 {
				log.Error(err, "error waiting for lock", "retries", retries)
			} else {
				log.V(1).Info("error waiting for lock", "err", err, "retries", retries)
			}
			retries += 1
		} else {
			retries = 0
		}
		n.leaderAcquired()
		n.lockClient.Unlock()
		n.leaderLost()

		if util_channels.IsClosed(stop) {
			break
		}
		time.Sleep(backoffTime)
	}

	if n.alwaysLeader {
		for _, callback := range n.callbacks {
			callback.OnStartedLeading()
		}
	}
}
func (n *zookeeperLeaderElector) leaderAcquired() {

}
func (n *zookeeperLeaderElector) leaderLost() {

}
