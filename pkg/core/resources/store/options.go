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

package store

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

import (
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
)

const (
	Group = "group"
	Key   = "key"

	Protocol  = "protocol"
	Address   = "address"
	DataId    = "data-id"
	Cluster   = "cluster"
	Username  = "username"
	Password  = "password"
	Namespace = "namespace"
	AppID     = "dubbo"
	Timeout   = "timeout"
)

type CreateOptions struct {
	Name         string
	Mesh         string
	CreationTime time.Time
	Owner        core_model.Resource
	Labels       map[string]string
}

type CreateOptionsFunc func(*CreateOptions)

func NewCreateOptions(fs ...CreateOptionsFunc) *CreateOptions {
	opts := &CreateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func CreateByGroupAndKey(group string, key string) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Labels = map[string]string{
			Group: group,
			Key:   key,
		}
	}
}

func CreateBy(key core_model.ResourceKey) CreateOptionsFunc {
	return CreateByKey(key.Name, key.Mesh)
}

func CreateByKey(name, mesh string) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func CreatedAt(creationTime time.Time) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.CreationTime = creationTime
	}
}

func CreateWithOwner(owner core_model.Resource) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Owner = owner
	}
}

func CreateWithLabels(labels map[string]string) CreateOptionsFunc {
	return func(opts *CreateOptions) {
		opts.Labels = labels
	}
}

type UpdateOptions struct {
	ModificationTime time.Time
	Labels           map[string]string
}

func ModifiedAt(modificationTime time.Time) UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.ModificationTime = modificationTime
	}
}

func UpdateWithLabels(labels map[string]string) UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.Labels = labels
	}
}

type UpdateOptionsFunc func(*UpdateOptions)

func NewUpdateOptions(fs ...UpdateOptionsFunc) *UpdateOptions {
	opts := &UpdateOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func UpdateGroupAndKey(group string, key string) UpdateOptionsFunc {
	return func(opts *UpdateOptions) {
		opts.Labels = map[string]string{
			Group: group,
			Key:   key,
		}
	}
}

type DeleteOptions struct {
	Name   string
	Mesh   string
	Labels map[string]string
}

type DeleteOptionsFunc func(*DeleteOptions)

func NewDeleteOptions(fs ...DeleteOptionsFunc) *DeleteOptions {
	opts := &DeleteOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func DeleteByGroupAndKey(group string, key string) DeleteOptionsFunc {
	return func(opts *DeleteOptions) {
		opts.Labels = map[string]string{
			Group: group,
			Key:   key,
		}
	}
}

func DeleteBy(key core_model.ResourceKey) DeleteOptionsFunc {
	return DeleteByKey(key.Name, key.Mesh)
}

func DeleteByKey(name, mesh string) DeleteOptionsFunc {
	return func(opts *DeleteOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

type DeleteAllOptions struct {
	Mesh string
}

type DeleteAllOptionsFunc func(*DeleteAllOptions)

func DeleteAllByMesh(mesh string) DeleteAllOptionsFunc {
	return func(opts *DeleteAllOptions) {
		opts.Mesh = mesh
	}
}

func NewDeleteAllOptions(fs ...DeleteAllOptionsFunc) *DeleteAllOptions {
	opts := &DeleteAllOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

type GetOptions struct {
	Name       string
	Mesh       string
	Version    string
	Consistent bool
	Labels     map[string]string
}

type GetOptionsFunc func(*GetOptions)

func NewGetOptions(fs ...GetOptionsFunc) *GetOptions {
	opts := &GetOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

func (g *GetOptions) HashCode() string {
	return fmt.Sprintf("%s:%s", g.Name, g.Mesh)
}

func GetByAddress(address string) GetOptionsFunc {
	return func(options *GetOptions) {
		if i := strings.Index(address, "://"); i > 0 {
			options.Labels = map[string]string{
				Protocol: address[0:i],
			}
			options.Labels = map[string]string{
				Address: address,
			}
		}
	}
}

func GetByTimeout(timeout time.Duration) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			Timeout: strconv.Itoa(int(timeout.Milliseconds())),
		}
	}
}

func GetByUsername(username string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			Username: username,
		}
	}
}

func GetByPassword(password string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			Password: password,
		}
	}
}

func GetByDataID(id string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{}
	}
}

func GetByCluster(cluster string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			Cluster: cluster,
		}
	}
}

func GetByNamespace(namespace string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			Namespace: namespace,
		}
	}
}

func GetByAppID(id string) GetOptionsFunc {
	return func(options *GetOptions) {
		options.Labels = map[string]string{
			AppID: id,
		}
	}
}

func GetBy(key core_model.ResourceKey) GetOptionsFunc {
	return GetByKey(key.Name, key.Mesh)
}

func GetByKey(name, mesh string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Name = name
		opts.Mesh = mesh
	}
}

func GetByVersion(version string) GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Version = version
	}
}

// GetConsistent forces consistency if storage provides eventual consistency like read replica for Postgres.
func GetConsistent() GetOptionsFunc {
	return func(opts *GetOptions) {
		opts.Consistent = true
	}
}

type (
	ListFilterFunc func(rs core_model.Resource) bool
)

type ListOptions struct {
	Mesh         string
	PageSize     int
	PageOffset   string
	FilterFunc   ListFilterFunc
	NameContains string
	Ordered      bool
	ResourceKeys map[core_model.ResourceKey]struct{}
}

type ListOptionsFunc func(*ListOptions)

func NewListOptions(fs ...ListOptionsFunc) *ListOptions {
	opts := &ListOptions{}
	for _, f := range fs {
		f(opts)
	}
	return opts
}

// Filter returns true if the item passes the filtering criteria
func (l *ListOptions) Filter(rs core_model.Resource) bool {
	if l.FilterFunc == nil {
		return true
	}

	return l.FilterFunc(rs)
}

func ListByNameContains(name string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.NameContains = name
	}
}

func ListByMesh(mesh string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.Mesh = mesh
	}
}

func ListByPage(size int, offset string) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.PageSize = size
		opts.PageOffset = offset
	}
}

func ListByFilterFunc(filterFunc ListFilterFunc) ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.FilterFunc = filterFunc
	}
}

func ListOrdered() ListOptionsFunc {
	return func(opts *ListOptions) {
		opts.Ordered = true
	}
}

func ListByResourceKeys(rk []core_model.ResourceKey) ListOptionsFunc {
	return func(opts *ListOptions) {
		resourcesKeys := map[core_model.ResourceKey]struct{}{}
		for _, val := range rk {
			resourcesKeys[val] = struct{}{}
		}
		opts.ResourceKeys = resourcesKeys
	}
}

func (l *ListOptions) IsCacheable() bool {
	return l.FilterFunc == nil
}

func (l *ListOptions) HashCode() string {
	return fmt.Sprintf("%s:%t:%s:%d:%s", l.Mesh, l.Ordered, l.NameContains, l.PageSize, l.PageOffset)
}
