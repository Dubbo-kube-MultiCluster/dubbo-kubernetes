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
	"errors"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/mysql"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	core_plugins "github.com/apache/dubbo-kubernetes/pkg/core/plugins"
	core_store "github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
	"github.com/apache/dubbo-kubernetes/pkg/events"
	mysql_events "github.com/apache/dubbo-kubernetes/pkg/plugins/resources/mysql/events"
)

type plugin struct{}

func init() {
	core_plugins.Register(core_plugins.MySQL, &plugin{})
}

func (p *plugin) NewResourceStore(pc core_plugins.PluginContext, config core_plugins.PluginConfig) (core_store.ResourceStore, core_store.Transactions, error) {
	cfg, ok := config.(*mysql.MysqlStoreConfig)
	if !ok {
		return nil, nil, errors.New("invalid type of the config. Passed config should be a PostgresStoreConfig")
	}
	store, err := NewMysqlStore(*cfg)
	if err != nil {
		return nil, nil, err
	}
	return store, nil, err
}

func (p *plugin) Migrate(pc core_plugins.PluginContext, config core_plugins.PluginConfig) (core_plugins.DbVersion, error) {
	cfg, ok := config.(*mysql.MysqlStoreConfig)
	if !ok {
		return 0, errors.New("invalid type of the config. Passed config should be a MysqlStoreConfig")
	}
	return MigrateDb(*cfg)
}

func (p *plugin) EventListener(pc core_plugins.PluginContext, out events.Emitter) error {
	mysqlListener := mysql_events.NewListener(*pc.Config().Store.Mysql, out)
	return pc.ComponentManager().Add(component.NewResilientComponent(core.Log.WithName("mysql-event-listener-component"), mysqlListener))
}
