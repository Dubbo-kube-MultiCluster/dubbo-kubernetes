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

package server

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/config/snp"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	"github.com/apache/dubbo-kubernetes/pkg/core/runtime/component"
)

var log = core.Log.WithName("dubbo").WithName("server")

const (
	grpcMaxConcurrentStreams = 1000000
	grpcKeepAliveTime        = 15 * time.Second
)

var _ component.Component = &Server{}

type Server struct {
	config     snp.SNPServerConfig
	grpcServer *grpc.Server
	instanceId string
}

func NewServer(
	config snp.SNPServerConfig,
) *Server {
	grpcOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepAliveTime,
			Timeout: grpcKeepAliveTime,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepAliveTime,
			PermitWithoutStream: true,
		}),
	}

	grpcOptions = append(grpcOptions)
	grpcServer := grpc.NewServer(grpcOptions...)

	return &Server{
		config:     config,
		grpcServer: grpcServer,
	}
}

func (s *Server) Start(stop <-chan struct{}) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.config.Port))
	if err != nil {
		return err
	}
	log := log.WithValues(
		"instanceId",
		s.instanceId,
	)

	errChan := make(chan error)
	go func() {
		defer close(errChan)
		if err := s.grpcServer.Serve(lis); err != nil {
			if err != http.ErrServerClosed {
				log.Error(err, "terminated with an error")
				errChan <- err
				return
			}
		}
		log.Info("terminated normally")
	}()
	log.Info("starting", "interface", "0.0.0.0", "port", s.config.Port, "tls", true)

	select {
	case <-stop:
		log.Info("stopping gracefully")
		s.grpcServer.GracefulStop()
		log.Info("stopped")
		return nil
	case err := <-errChan:
		return err
	}
}

func (s *Server) NeedLeaderElection() bool {
	return false
}

func (s *Server) GrpcServer() *grpc.Server {
	return s.grpcServer
}
