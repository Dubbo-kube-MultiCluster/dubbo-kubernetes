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

package events

import (
	"encoding/json"
	"github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/zookeeper"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	"github.com/apache/dubbo-kubernetes/pkg/events"
	common_zookeeper "github.com/apache/dubbo-kubernetes/pkg/plugins/common/zookeeper"
	"github.com/pkg/errors"
)

var log = core.Log.WithName("zookeeper-event-listener")

type listener struct {
	cfg zookeeper.ZookeeperStoreConfig
	out events.Emitter
}

func NewListener(cfg zookeeper.ZookeeperStoreConfig, out events.Emitter) component.Component {
	return &listener{
		out: out,
		cfg: cfg,
	}
}

func (k *listener) Start(stop <-chan struct{}) error {
	var err error
	var listener common_zookeeper.Listener
	listener, err = common_zookeeper.NewListener(k.cfg, core.Log.WithName("zk-event-listener"))
	if err != nil {
		return err
	}
	defer func() {
		if err := listener.Close(); err != nil {
			log.Error(err, "error closing zk listener")
		}
	}()

	for {
		select {
		case err := <-listener.Error():
			log.Error(err, "failed to listen on events")
			return err
		case n := <-listener.Notify():
			if n == nil {
				continue
			}
			obj := &struct {
				Action string `json:"action"`
				Data   struct {
					Name     string `json:"name"`
					Mesh     string `json:"mesh"`
					Type     string `json:"type"`
					TenantID string `json:"tenant_id"` // It is always empty with current implementation
				}
			}{}
			if err := json.Unmarshal([]byte(n.Payload), obj); err != nil {
				log.Error(err, "unable to unmarshal event from PostgreSQL")
				continue
			}
			var op events.Op
			switch obj.Action {
			case "INSERT":
				op = events.Create
			case "UPDATE":
				op = events.Update
			case "DELETE":
				op = events.Delete
			default:
				log.Error(errors.New("unknown Action"), "failed to parse action", "action", op)
				continue
			}
			k.out.Send(events.ResourceChangedEvent{
				Operation: op,
				Type:      model.ResourceType(obj.Data.Type),
				Key:       model.ResourceKey{Mesh: obj.Data.Mesh, Name: obj.Data.Name},
				TenantID:  obj.Data.TenantID,
			})
		case <-stop:
			log.Info("stop")
			return nil
		}
	}
}

func (k *listener) NeedLeaderElection() bool {
	return false
}
