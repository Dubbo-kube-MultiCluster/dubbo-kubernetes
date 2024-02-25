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
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"

	"github.com/dubbogo/go-zookeeper/zk"

	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"

	"github.com/pkg/errors"
)

import (
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	"github.com/apache/dubbo-kubernetes/pkg/events"
)

const (
	RootPath      = "/dubbo/config"
	pathSeparator = "/"
	DefaultGroup  = "dubbo"
)

type zookeeperStore struct {
	url    *common.URL
	client *gxzookeeper.ZookeeperClient
}

func NewStore() store.ResourceStore {
	return &zookeeperStore{}
}

func (c *zookeeperStore) getPath(key string, group string) string {
	if len(key) == 0 {
		return c.buildPath(group)
	}
	return c.buildPath(group) + pathSeparator + key
}

func (c *zookeeperStore) buildPath(group string) string {
	if len(group) == 0 {
		group = DefaultGroup
	}
	return RootPath + pathSeparator + group
}

func (c *zookeeperStore) SetEventWriter(writer events.Emitter) {
}

func (c *zookeeperStore) Create(_ context.Context, resource core_model.Resource, fs ...store.CreateOptionsFunc) error {
	opts := store.NewCreateOptions(fs...)
	group := opts.Labels[store.Group]
	key := opts.Labels[store.Key]
	path := c.getPath(key, group)

	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	err = c.client.CreateWithValue(path, bytes)
	if err != nil {
		// try update value if node already exists
		if errors.Is(err, zk.ErrNodeExists) {
			_, stat, _ := c.client.GetContent(path)
			_, setErr := c.client.SetContent(path, bytes, stat.Version)
			if setErr != nil {
				return errors.WithStack(setErr)
			}
			return nil
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *zookeeperStore) Update(_ context.Context, resource core_model.Resource, fs ...store.UpdateOptionsFunc) error {
	opts := store.NewUpdateOptions(fs...)
	group := opts.Labels[store.Group]
	key := opts.Labels[store.Key]
	path := c.getPath(key, group)

	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}
	bytes = []byte(base64.StdEncoding.EncodeToString(bytes))
	err = c.client.CreateWithValue(path, bytes)
	if err != nil {
		// try update value if node already exists
		if errors.Is(err, zk.ErrNodeExists) {
			_, stat, _ := c.client.GetContent(path)
			_, setErr := c.client.SetContent(path, bytes, stat.Version)
			if setErr != nil {
				return errors.WithStack(setErr)
			}
			return nil
		}
		return errors.WithStack(err)
	}

	return nil
}

func (c *zookeeperStore) Delete(ctx context.Context, resource core_model.Resource, fs ...store.DeleteOptionsFunc) error {
	opts := store.NewDeleteOptions(fs...)
	group := opts.Labels[store.Group]
	key := opts.Labels[store.Key]
	path := c.getPath(key, group)

	err := c.client.Delete(path)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *zookeeperStore) Get(_ context.Context, resource core_model.Resource, fs ...store.GetOptionsFunc) error {
	opts := store.NewGetOptions(fs...)
	key := opts.Labels[store.Key]
	content, err := c.GetProperties(key, opts.Labels)
	if err != nil {
		return errors.Wrap(err, "No configuration found in config center.")
	}
	if err := core_model.FromJSON([]byte(content), resource.GetSpec()); err != nil {
		return errors.Wrap(err, "failed to convert json to spec")
	}
	meta := &resourceMetaObject{
		Name: opts.Name,
		Mesh: opts.Mesh,
	}
	resource.SetMeta(meta)

	if opts.Version != "" && resource.GetMeta().GetVersion() != opts.Version {
		return store.ErrorResourceConflict(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}
	return nil
}

func (c *zookeeperStore) List(_ context.Context, resource core_model.ResourceList, fs ...store.ListOptionsFunc) error {
	return nil
}

func (c *zookeeperStore) GetProperties(key string, labels map[string]string) (string, error) {
	/**
	 * when group is not null, we are getting startup configs from Config Center, for example:
	 * group=dubbo, key=dubbo.properties
	 */
	group := labels[store.Group]
	if group != "" {
		key = group + "/" + key
	} else {
		// TODO
	}
	content, _, err := c.client.GetContent(RootPath + "/" + key)
	if err != nil {
		return "", errors.WithStack(err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(content))
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(decoded), nil
}

func (c *zookeeperStore) GetURL() *common.URL {
	return c.url
}
