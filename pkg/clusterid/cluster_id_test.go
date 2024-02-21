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

package clusterid_test

import (
	"context"
	"github.com/apache/dubbo-kubernetes/pkg/clusterid"
	"github.com/apache/dubbo-kubernetes/pkg/test/runtime"

	dubbo_cp "github.com/apache/dubbo-kubernetes/pkg/config/app/dubbo-cp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cluster ID", func() {
	var stop chan struct{}

	AfterEach(func() {
		if stop != nil {
			close(stop)
		}
	})

	It("should create and set cluster ID", func() {
		// given runtime with cluster ID components
		cfg := dubbo_cp.DefaultConfig()
		builder, err := runtime.BuilderFor(context.Background(), cfg)
		Expect(err).ToNot(HaveOccurred())
		runtime, err := builder.Build()
		Expect(err).ToNot(HaveOccurred())

		err = clusterid.Setup(runtime)
		Expect(err).ToNot(HaveOccurred())

		// when runtime is started
		stop = make(chan struct{})
		go func() {
			defer GinkgoRecover()
			err := runtime.Start(stop)
			Expect(err).ToNot(HaveOccurred())
		}()

		// then cluster ID is created and set in Runtime object
		Eventually(runtime.GetClusterId, "5s", "100ms").ShouldNot(BeEmpty())
	})
})
