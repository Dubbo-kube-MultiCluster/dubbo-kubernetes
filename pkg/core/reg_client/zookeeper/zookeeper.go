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
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"

	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/core/consts"
	"github.com/apache/dubbo-kubernetes/pkg/core/extensions"
	"github.com/apache/dubbo-kubernetes/pkg/core/reg_client"
	"github.com/apache/dubbo-kubernetes/pkg/core/reg_client/factory"
)

func init() {
	mf := &zookeeperRegClientFactory{}
	extensions.SetRegClientFactory("zookeeper", func() factory.RegClientFactory {
		return mf
	})
}

type zookeeperRegClient struct {
	client *gxzookeeper.ZookeeperClient
}

func (m *zookeeperRegClient) GetKeys() {

}

type zookeeperRegClientFactory struct{}

func (mf *zookeeperRegClientFactory) CreateRegClient(url *common.URL) reg_client.RegClient {
	client, err := gxzookeeper.NewZookeeperClient(
		"zookeeperRegClient",
		strings.Split(url.Location, ","),
		false,
		gxzookeeper.WithZkTimeOut(url.GetParamDuration(consts.TimeoutKey, "15s")),
	)
	if err != nil {
		panic(err)
	}

	return &zookeeperRegClient{client: client}
}
