/*
Copyright 2025 Finvi, Ontario Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package ocep_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/oracle/oci-go-sdk/v65/common"

	. "github.com/onsi/gomega"
	ocep "github.com/ontariosystems/oci-cli-env-provider"
)

var _ = Describe("ComposingConfigProvider", func() {
	It("returns first valid AuthType", func() {
		conf := ocep.ComposingConfigProvider(&noOpProvider{}, &noOpProvider{}, &testProvider{authType: common.UserPrincipal}, &testProvider{authType: common.InstancePrincipal})
		at, err := conf.AuthType()
		Expect(err).ToNot(HaveOccurred())
		Expect(at.AuthType).To(Equal(common.UserPrincipal))
	})

	It("returns error if no valid AuthType", func() {
		conf := ocep.ComposingConfigProvider(&noOpProvider{}, &noOpProvider{})
		at, err := conf.AuthType()
		Expect(err).To(MatchError(ocep.ErrNoAuthType))
		Expect(at.AuthType).To(Equal(common.UnknownAuthenticationType))
	})
})

type testProvider struct {
	authType common.AuthenticationType
	noOpProvider
}

func (p *testProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{AuthType: p.authType}, nil
}
