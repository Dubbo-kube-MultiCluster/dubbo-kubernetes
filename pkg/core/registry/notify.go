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

package registry

import (
	"context"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/manager"
)

type NotifyListener struct {
	manager.ResourceManager
}

func NewNotifyListener(manager manager.ResourceManager) *NotifyListener {
	return &NotifyListener{
		manager,
	}
}

func (l *NotifyListener) Notify(event *dubboRegistry.ServiceEvent) {
	switch event.Action {
	case remoting.EventTypeAdd, remoting.EventTypeUpdate:
		if err := l.createOrUpdateDataplane(context.Background()); err != nil {
			return
		}
	case remoting.EventTypeDel:
		if err := l.deleteDataplane(context.Background()); err != nil {
			return
		}
	}
}

func (l *NotifyListener) NotifyAll(events []*dubboRegistry.ServiceEvent, f func()) {
	for _, event := range events {
		l.Notify(event)
	}
}

func (l *NotifyListener) deleteDataplane(ctx context.Context) error {
	return nil
}

func (l *NotifyListener) createOrUpdateDataplane(ctx context.Context) error {
	return nil
}
