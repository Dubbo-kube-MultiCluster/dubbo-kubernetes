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

import Mock from 'mockjs'

Mock.mock('/mock/application/search', 'get', () => {
  let total = Mock.mock('@integer(8, 1000)')
  let list = []
  for (let i = 0; i < total; i++) {
    list.push({
      appName: 'app_' + Mock.mock('@string(2,10)'),
      instanceNum: Mock.mock('@integer(80, 200)'),
      deployCluster: 'cluster_' + Mock.mock('@string(5)'),
      'registerClusters|1-3': ['cluster_' + Mock.mock('@string(5)')]
    })
  }
  return {
    code: 200,
    message: 'success',
    data: Mock.mock({
      total: total,
      curPage: 1,
      pageSize: 10,
      data: list
    })
  }
})
