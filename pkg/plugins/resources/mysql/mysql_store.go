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
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"strconv"
	"time"
)

import (
	"github.com/pkg/errors"

	"gorm.io/gorm"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/config/plugins/resources/mysql"
	core_model "github.com/apache/dubbo-kubernetes/pkg/core/resources/model"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/store"
	common_mysql "github.com/apache/dubbo-kubernetes/pkg/plugins/common/mysql"
	"github.com/apache/dubbo-kubernetes/pkg/plugins/resources/mysql/model"
)

type mysqlResourceStore struct {
	db *gorm.DB
}

type ResourceNamesByMesh map[string][]string

func NewMysqlStore(config mysql.MysqlStoreConfig) (store.ResourceStore, error) {
	db, err := common_mysql.ConnectToDb(config)
	if err != nil {
		return nil, err
	}

	return &mysqlResourceStore{db: db}, nil
}

func (r *mysqlResourceStore) Create(_ context.Context, resource core_model.Resource, fs ...store.CreateOptionsFunc) error {
	opts := store.NewCreateOptions(fs...)

	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return errors.Wrap(err, "failed to convert spec to json")
	}

	version := 0
	labels, err := prepareLabels(opts.Labels)
	if err != nil {
		return err
	}

	resources := model.Resources{
		Name:             opts.Name,
		Mesh:             opts.Mesh,
		Type:             string(resource.Descriptor().Name),
		Version:          version,
		Spec:             string(bytes),
		Labels:           labels,
		ModificationTime: opts.CreationTime.UTC(),
		CreateTime:       opts.CreationTime.UTC(),
		OwnerMesh:        opts.Owner.GetMeta().GetName(),
		OwnerName:        opts.Owner.GetMeta().GetName(),
		OwnerType:        string(opts.Owner.Descriptor().Name),
	}
	err = r.db.Create(&resources).Error
	if err != nil {
		return errors.Wrapf(err, "failed to create resource")
	}

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

func (r *mysqlResourceStore) Update(_ context.Context, resource core_model.Resource, fs ...store.UpdateOptionsFunc) error {
	bytes, err := core_model.ToJSON(resource.GetSpec())
	if err != nil {
		return err
	}

	opts := store.NewUpdateOptions(fs...)

	version, err := strconv.Atoi(resource.GetMeta().GetVersion())
	newVersion := version + 1
	if err != nil {
		return errors.Wrap(err, "failed to convert meta version to int")
	}

	labels, err := prepareLabels(opts.Labels)
	if err != nil {
		return err
	}

	err = r.db.Raw("UPDATE resources SET spec = ?, version = ?, modification_time = ?, labels = ? WHERE name = ? AND mesh = ? AND type = ? AND version = ?;",
		string(bytes),
		newVersion,
		opts.ModificationTime.UTC().String(),
		labels,
		resource.GetMeta().GetName(),
		resource.GetMeta().GetMesh(),
		resource.Descriptor().Name,
		version,
	).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update resources")
	}
	resource.SetMeta(&resourceMetaObject{
		Name:             resource.GetMeta().GetName(),
		Mesh:             resource.GetMeta().GetMesh(),
		Version:          strconv.Itoa(newVersion),
		ModificationTime: opts.ModificationTime,
		Labels:           maps.Clone(opts.Labels),
	})
	return nil
}

func (r *mysqlResourceStore) Delete(_ context.Context, resource core_model.Resource, fs ...store.DeleteOptionsFunc) error {
	opts := store.NewDeleteOptions(fs...)

	err := r.db.Raw("DELETE FROM resources WHERE name = ? AND type = ? AND mesh = ?;",
		opts.Name,
		resource.Descriptor().Name,
		opts.Mesh,
	).Error
	if err != nil {
		return errors.Wrapf(err, "failed to delete resource")
	}

	return nil
}

func (r *mysqlResourceStore) Get(_ context.Context, resource core_model.Resource, fs ...store.GetOptionsFunc) error {
	opts := store.NewGetOptions(fs...)

	resources := model.Resources{}
	err := r.db.Raw("SELECT spec, version, creation_time, modification_time, labels FROM resources WHERE name = ? AND mesh = ? AND type = ?;",
		opts.Name,
		opts.Mesh,
		resource.Descriptor().Name,
	).Scan(&resources).Error
	if err != nil {
		return errors.Wrap(err, "failed to get resources")
	}

	if err := core_model.FromJSON([]byte(resources.Spec), resource.GetSpec()); err != nil {
		return errors.Wrap(err, "failed to convert json to spec")
	}

	meta := &resourceMetaObject{
		Name:             opts.Name,
		Mesh:             opts.Mesh,
		Version:          strconv.Itoa(resources.Version),
		CreationTime:     resources.CreateTime.Local(),
		ModificationTime: resources.ModificationTime.Local(),
		Labels:           map[string]string{},
	}
	if err := json.Unmarshal([]byte(resources.Labels), &meta.Labels); err != nil {
		return errors.Wrap(err, "failed to convert json to labels")
	}
	resource.SetMeta(meta)

	if opts.Version != "" && resource.GetMeta().GetVersion() != opts.Version {
		return store.ErrorResourceConflict(resource.Descriptor().Name, opts.Name, opts.Mesh)
	}
	return nil
}

func (r *mysqlResourceStore) List(_ context.Context, resources core_model.ResourceList, args ...store.ListOptionsFunc) error {
	opts := store.NewListOptions(args...)

	statement := "SELECT name, mesh, spec, version, creation_time, modification_time, labels FROM resources WHERE type= ?;"
	var statementArgs []interface{}
	statementArgs = append(statementArgs, resources.GetItemType())
	argsIndex := 1
	statement += " AND ("
	res := resourceNamesByMesh(opts.ResourceKeys)
	iter := 0
	for mesh, names := range res {
		if iter > 0 {
			statement += " OR "
		}
		argsIndex++
		statement += fmt.Sprintf("(mesh=$%d AND", argsIndex)
		statementArgs = append(statementArgs, mesh)
		for idx, name := range names {
			argsIndex++
			if idx == 0 {
				statement += fmt.Sprintf(" name IN ($%d", argsIndex)
			} else {
				statement += fmt.Sprintf(",$%d", argsIndex)
			}
			statementArgs = append(statementArgs, name)
		}
		statement += "))"
		iter++
	}
	statement += ")"

	if opts.Mesh != "" {
		argsIndex++
		statement += fmt.Sprintf(" AND mesh=$%d", argsIndex)
		statementArgs = append(statementArgs, opts.Mesh)
	}
	if opts.NameContains != "" {
		argsIndex++
		statement += fmt.Sprintf(" AND name LIKE $%d", argsIndex)
		statementArgs = append(statementArgs, "%"+opts.NameContains+"%")
	}
	statement += " ORDER BY name, mesh"

	rows, err := r.db.Raw(statement, statementArgs...).Rows()
	if err != nil {
		return errors.Wrapf(err, "failed to execute query: %s", statement)
	}
	defer rows.Close()

	total := 0
	for rows.Next() {
		item, err := rowToItem(resources, rows)
		if err != nil {
			return err
		}
		if err := resources.AddItem(item); err != nil {
			return err
		}
		total++
	}

	resources.GetPagination().SetTotal(uint32(total))
	return nil
}

func rowToItem(resources core_model.ResourceList, rows *sql.Rows) (core_model.Resource, error) {
	var name, mesh, spec string
	var version int
	var creationTime, modificationTime time.Time
	var labels string
	if err := rows.Scan(&name, &mesh, &spec, &version, &creationTime, &modificationTime, &labels); err != nil {
		return nil, errors.Wrap(err, "failed to retrieve elements from query")
	}

	item := resources.NewItem()
	if err := core_model.FromJSON([]byte(spec), item.GetSpec()); err != nil {
		return nil, errors.Wrap(err, "failed to convert json to spec")
	}

	meta := &resourceMetaObject{
		Name:             name,
		Mesh:             mesh,
		Version:          strconv.Itoa(version),
		CreationTime:     creationTime.Local(),
		ModificationTime: modificationTime.Local(),
		Labels:           map[string]string{},
	}
	if err := json.Unmarshal([]byte(labels), &meta.Labels); err != nil {
		return nil, errors.Wrap(err, "failed to convert json to labels")
	}
	item.SetMeta(meta)

	return item, nil
}

func prepareLabels(labels map[string]string) (string, error) {
	lblBytes, err := json.Marshal(labels)
	if err != nil {
		return "", errors.Wrap(err, "failed to convert labels to json")
	}
	return string(lblBytes), nil
}

func resourceNamesByMesh(resourceKeys map[core_model.ResourceKey]struct{}) ResourceNamesByMesh {
	resourceNamesByMesh := ResourceNamesByMesh{}
	for key := range resourceKeys {
		if val, exists := resourceNamesByMesh[key.Mesh]; exists {
			resourceNamesByMesh[key.Mesh] = append(val, key.Name)
		} else {
			resourceNamesByMesh[key.Mesh] = []string{key.Name}
		}
	}
	return resourceNamesByMesh
}
