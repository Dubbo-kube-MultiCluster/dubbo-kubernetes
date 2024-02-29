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

package envoy

import (
	"github.com/apache/dubbo-kubernetes/pkg/core"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/model/rest"
	"io"
	"sync"
)

var runLog = core.Log.WithName("dubbo-proxy").WithName("run").WithName("envoy")

type BootstrapParams struct {
	Dataplane           rest.Resource
	DNSPort             uint32
	EmptyDNSPort        uint32
	EnvoyVersion        EnvoyVersion
	DynamicMetadata     map[string]string
	Workdir             string
	MetricsSocketPath   string
	AccessLogSocketPath string
	MetricsCertPath     string
	MetricsKeyPath      string
}

type EnvoyVersion struct {
	Build            string
	Version          string
	KumaDpCompatible bool
}

type Opts struct {
	BootstrapConfig []byte
	AdminPort       uint32
	Dataplane       rest.Resource
	Stdout          io.Writer
	Stderr          io.Writer
	OnFinish        func()
}

type Envoy struct {
	opts Opts

	wg sync.WaitGroup
}

func New(opts Opts) (*Envoy, error) {

	if opts.OnFinish == nil {
		opts.OnFinish = func() {}
	}
	return &Envoy{opts: opts}, nil
}
