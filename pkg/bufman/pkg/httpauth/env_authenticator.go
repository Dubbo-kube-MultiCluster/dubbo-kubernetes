// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpauth

import (
	"net/http"
)

import (
	"github.com/apache/dubbo-kubernetes/pkg/bufman/pkg/app"
)

type envAuthenticator struct {
	usernameKey string
	passwordKey string
}

func newEnvAuthenticator(
	usernameKey string,
	passwordKey string,
) *envAuthenticator {
	return &envAuthenticator{
		usernameKey: usernameKey,
		passwordKey: passwordKey,
	}
}

func (a *envAuthenticator) SetAuth(envContainer app.EnvContainer, request *http.Request) (bool, error) {
	return setBasicAuth(
		request,
		envContainer.Env(a.usernameKey),
		envContainer.Env(a.passwordKey),
		a.usernameKey,
		a.passwordKey,
	)
}
