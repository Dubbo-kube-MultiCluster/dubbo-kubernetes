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
	"github.com/apache/dubbo-kubernetes/pkg/config/app/dubboctl"
	"github.com/apache/dubbo-kubernetes/pkg/core"
	"github.com/apache/dubbo-kubernetes/pkg/core/resources/model/rest"
	"github.com/apache/dubbo-kubernetes/pkg/util/files"
	"github.com/pkg/errors"
	"io"
	"os/exec"
	"regexp"
	"strings"
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
	ProxyArgs       dubboctl.ProxyArgs
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

func lookupEnvoyPath(configuredPath string) (string, error) {
	return files.LookupBinaryPath(
		files.LookupInPath(configuredPath),
		files.LookupInCurrentDirectory("envoy"),
		files.LookupNextToCurrentExecutable("envoy"),
	)
}

func GetEnvoyVersion(binaryPath string) (*EnvoyVersion, error) {
	resolvedPath, err := lookupEnvoyPath(binaryPath)
	if err != nil {
		return nil, err
	}
	arg := "--version"
	command := exec.Command(resolvedPath, arg)
	output, err := command.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute %s with arguments %q", resolvedPath, arg)
	}
	build := strings.ReplaceAll(string(output), "\r\n", "\n")
	build = strings.Trim(build, "\n")
	build = regexp.MustCompile(`version:(.*)`).FindString(build)
	build = strings.Trim(build, "version:")
	build = strings.Trim(build, " ")

	parts := strings.Split(build, "/")
	if len(parts) != 5 { // revision/build_version_number/revision_status/build_type/ssl_version
		return nil, errors.Errorf("wrong Envoy build format: %s", build)
	}
	return &EnvoyVersion{
		Build:   build,
		Version: parts[1],
	}, nil
}
