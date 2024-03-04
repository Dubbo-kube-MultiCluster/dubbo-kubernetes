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

package client

import (
	"io"
)

import (
	"github.com/go-logr/logr"

	"github.com/pkg/errors"
)

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	core_mesh "github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/mesh"
)

type Callbacks struct {
	OnResourcesReceived func(request *mesh_proto.MappingSyncRequest) (core_mesh.MappingResourceList, error)
}

// MappingSyncClient Handle MappingSyncRequest from client
type MappingSyncClient interface {
	HandleReceive() error
}

type mappingSyncClient struct {
	log        logr.Logger
	syncStream MappingSyncStream
	callbacks  *Callbacks
}

func NewMappingSyncClient(log logr.Logger, syncStream MappingSyncStream, cb *Callbacks) MappingSyncClient {
	return &mappingSyncClient{
		log:        log,
		syncStream: syncStream,
		callbacks:  cb,
	}
}

func (s *mappingSyncClient) HandleReceive() error {
	for {
		received, err := s.syncStream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return errors.Wrap(err, "failed to receive a MappingSyncRequest")
		}

		var mappingList core_mesh.MappingResourceList
		if s.callbacks == nil {
			// if no callbacks, send NACK
			s.log.Info("no callback set, sending NACK")
			if err := s.syncStream.NACK(received.InterfaceNames, "no callbacks in control plane, could not get Mapping"); err != nil {
				if err == io.EOF {
					return nil
				}
				return errors.Wrap(err, "failed to NACK a MappingSyncRequest")
			}
			continue
		}

		// callbacks to get Mapping Resource List
		mappingList, err = s.callbacks.OnResourcesReceived(received)
		if err != nil {
			s.log.Info("error during search Mapping Resource, sending NACK", "err", err)
			if err := s.syncStream.NACK(received.InterfaceNames, err.Error()); err != nil {
				if err == io.EOF {
					return nil
				}
				return errors.Wrap(err, "failed to NACK a MappingSyncRequest")
			}
		} else {
			s.log.V(1).Info("sending ACK")
			if err := s.syncStream.ACK(mappingList); err != nil {
				if err == io.EOF {
					return nil
				}
				return errors.Wrap(err, "failed to ACK a MappingSyncRequest")
			}
		}

	}
}
