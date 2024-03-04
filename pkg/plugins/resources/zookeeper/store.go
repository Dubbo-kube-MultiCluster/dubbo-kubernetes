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
	"context"
	"encoding/base64"
	"maps"
	"strconv"
)

import (
	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"

	"github.com/pkg/errors"
)

import (
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/events"
)

type zookeeperStore struct {
	client *gxzookeeper.ZookeeperClient
}

func NewStore(zkClient *gxzookeeper.ZookeeperClient) store.ResourceStore {
	return &zookeeperStore{
		client: zkClient,
	}
}

func (c *zookeeperStore) SetEventWriter(writer events.Emitter) {
}

func (c *zookeeperStore) Create(_ context.Context, resource core_model.Resource, fs ...store.CreateOptionsFunc) error {
	opts := store.NewCreateOptions(fs...)

	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	valueBytes := []byte(base64.StdEncoding.EncodeToString(bytes))

	path := opts.Labels[store.PathLabel]

	err = c.client.CreateWithValue(path, valueBytes)
	// 如果创建失败的话直接返回空
	if err != nil {
		return err
	}

	// 创建成功, 通过GetContent获取到version
	_, stat, _ := c.client.GetContent(path)
	version := int(stat.Version)

	resource.SetMeta(&resourceMetaObject{
		Name:             opts.Name,
		Mesh:             opts.Mesh,
		Version:          strconv.Itoa(version),
		CreationTime:     opts.CreationTime,
		ModificationTime: opts.CreationTime,
		Labels:           maps.Clone(opts.Labels),
	})
	return nil
}

func (c *zookeeperStore) Update(_ context.Context, resource core_model.Resource, fs ...store.UpdateOptionsFunc) error {
	opts := store.NewUpdateOptions(fs...)

	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	valueBytes := []byte(base64.StdEncoding.EncodeToString(bytes))

	path := opts.Labels[store.PathLabel]

	_, stat, _ := c.client.GetContent(path)
	zkStat, setErr := c.client.SetContent(path, valueBytes, stat.Version)
	if setErr != nil {
		return errors.WithStack(setErr)
	}

	resource.SetMeta(&resourceMetaObject{
		Name:             resource.GetMeta().GetName(),
		Mesh:             resource.GetMeta().GetMesh(),
		Version:          strconv.Itoa(int(zkStat.Version)),
		ModificationTime: opts.ModificationTime,
		Labels:           maps.Clone(opts.Labels),
	})
	return nil
}

func (c *zookeeperStore) Delete(ctx context.Context, resource core_model.Resource, fs ...store.DeleteOptionsFunc) error {
	opts := store.NewDeleteOptions(fs...)
	path := opts.Labels[store.PathLabel]
	err := c.client.Delete(path)
	if err != nil {
		return err
	}
	return nil
}

func (c *zookeeperStore) Get(_ context.Context, resource core_model.Resource, fs ...store.GetOptionsFunc) error {
	opts := store.NewGetOptions(fs...)
	path := opts.Labels[store.PathLabel]
	content, stat, err := c.client.GetContent(path)
	if err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return err
	}
	if err := core_model.FromJSON(decoded, resource.GetSpec()); err != nil {
		return errors.Wrap(err, "failed to convert json to spec")
	}

	meta := &resourceMetaObject{
		Name:    opts.Name,
		Mesh:    opts.Mesh,
		Version: strconv.Itoa(int(stat.Version)),
	}
	resource.SetMeta(meta)
	return nil
}

// List 在zk store的场景下List和Get实现一样, List取决于path, path由manager层传
// TODO
func (c *zookeeperStore) List(_ context.Context, resource core_model.ResourceList, fs ...store.ListOptionsFunc) error {
	return nil
}
