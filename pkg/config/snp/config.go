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

package snp

import (
	"time"

	config_types "github.com/apache/dubbo-kubernetes/pkg/config/types"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

func DefaultServiceNameMappingConfig() SNPConfig {
	return SNPConfig{
		Server: SNPServerConfig{
			Port:            57384,
			TlsMinVersion:   "TLSv1_2",
			TlsCipherSuites: []string{},
		},
	}
}

type SNPConfig struct {
	ServiceMapping ServiceMapping `json:"serviceMapping"`
	// Service Name Mapping server configuration
	Server SNPServerConfig `json:"server"`
}

func (s *SNPConfig) Validate() error {
	if err := s.ServiceMapping.Validate(); err != nil {
		return errors.Wrap(err, ".ServiceMapping validation failed")
	}

	if err := s.Server.Validate(); err != nil {
		return errors.Wrap(err, ".Server validation failed")
	}
	return nil
}

type SNPServerConfig struct {
	// Port on which Service Name Mapping server will listen
	Port uint16 `json:"port" envconfig:"dubbo_snp_server_port"`
	// TlsMinVersion defines the minimum TLS version to be used
	TlsMinVersion string `json:"tlsMinVersion" envconfig:"dubbo_snp_server_tls_min_version"`
	// TlsMaxVersion defines the maximum TLS version to be used
	TlsMaxVersion string `json:"tlsMaxVersion" envconfig:"dubbo_snp_server_tls_max_version"`
	// TlsCipherSuites defines the list of ciphers to use
	TlsCipherSuites []string `json:"tlsCipherSuites" envconfig:"dubbo_snp_server_tls_cipher_suites"`
}

func (s *SNPServerConfig) Validate() error {
	var errs error
	if s.Port == 0 {
		errs = multierr.Append(errs, errors.New(".Port cannot be zero"))
	}
	if _, err := config_types.TLSVersion(s.TlsMinVersion); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsMinVersion "+err.Error()))
	}
	if _, err := config_types.TLSVersion(s.TlsMaxVersion); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsMaxVersion "+err.Error()))
	}
	if _, err := config_types.TLSCiphers(s.TlsCipherSuites); err != nil {
		errs = multierr.Append(errs, errors.New(".TlsCipherSuites "+err.Error()))
	}
	return errs
}

type ServiceMapping struct {
	Debounce Debounce `json:"debounce"`
}

func (s *ServiceMapping) Validate() error {
	var errs error
	if err := s.Debounce.Validate(); err != nil {
		errs = multierr.Append(errs, errors.Wrap(err, ".Debounce validation failed"))
	}

	return errs
}

type Debounce struct {
	After  time.Duration `yaml:"after"`
	Max    time.Duration `yaml:"max"`
	Enable bool          `yaml:"enable"`
}

func (s *Debounce) Validate() error {
	var errs error

	afterThreshold := time.Second * 10
	if s.After > afterThreshold {
		errs = multierr.Append(errs, errors.New(".After can not greater than "+afterThreshold.String()))
	}

	maxThreshold := time.Second * 10
	if s.Max > maxThreshold {
		errs = multierr.Append(errs, errors.New(".Max can not greater than "+maxThreshold.String()))
	}

	return errs
}
