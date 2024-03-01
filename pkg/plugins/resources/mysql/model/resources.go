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

package model

import (
	"time"
)

type Resources struct {
	Name             string    `json:"name" gorm:"name"`
	Namespace        string    `json:"namespace" gorm:"namespace"`
	Mesh             string    `json:"mesh" gorm:"mesh"`
	Type             string    `json:"type" gorm:"type"`
	Version          int       `json:"version" gorm:"version"`
	Spec             string    `json:"spec" gorm:"spec"`
	Labels           string    `json:"labels" gorm:"labels"`
	ModificationTime time.Time `json:"modification_time" gorm:"modification_time"`
	CreateTime       time.Time `json:"create_time" gorm:"creation_time"`
	OwnerName        string    `json:"owner_name" gorm:"owner_name"`
	OwnerMesh        string    `json:"owner_mesh" gorm:"owner_mesh"`
	OwnerType        string    `json:"owner_type" gorm:"owner_type"`
}
