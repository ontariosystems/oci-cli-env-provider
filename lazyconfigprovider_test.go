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
	"crypto/rsa"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	ocep "github.com/ontariosystems/oci-cli-env-provider"
	"github.com/oracle/oci-go-sdk/v65/common"
)

var _ = Describe("LazyConfigProvider", func() {
	var (
		initCounter      int
		initError        error
		conf             common.ConfigurationProvider
		noOpTestProvider = func() (common.ConfigurationProvider, error) {
			initCounter += 1
			return &noOpProvider{}, initError
		}
	)

	BeforeEach(func() {
		initError = nil
		initCounter = 0
		conf = ocep.LazyConfigProvider(noOpTestProvider)
	})

	It("only initializes once", func() {
		_, _ = conf.PrivateRSAKey()
		_, _ = conf.KeyID()
		_, _ = conf.TenancyOCID()
		_, _ = conf.UserOCID()
		_, _ = conf.KeyFingerprint()
		_, _ = conf.Region()
		_, _ = conf.AuthType()

		Expect(initCounter).To(Equal(1))
	})

	It("does not initialize if not called", func() {
		providers := []common.ConfigurationProvider{&noOpProvider{}, conf}
		compProvider, _ := common.ComposingConfigurationProvider(providers)

		_, _ = compProvider.KeyID()
		Expect(initCounter).To(Equal(0))
	})

	When("it throws an error", func() {
		BeforeEach(func() {
			initError = errors.New("some error")
		})

		It("returns an error", func() {
			_, err := conf.PrivateRSAKey()
			Expect(err).To(MatchError(initError))

			_, err = conf.KeyID()
			Expect(err).To(MatchError(initError))

			_, err = conf.TenancyOCID()
			Expect(err).To(MatchError(initError))

			_, err = conf.UserOCID()
			Expect(err).To(MatchError(initError))

			_, err = conf.KeyFingerprint()
			Expect(err).To(MatchError(initError))

			_, err = conf.Region()
			Expect(err).To(MatchError(initError))

			_, err = conf.AuthType()
			Expect(err).To(MatchError(initError))
		})
	})
})

type noOpProvider struct{}

func (n noOpProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	return nil, nil
}

func (n noOpProvider) KeyID() (string, error) {
	return "", nil
}

func (n noOpProvider) TenancyOCID() (string, error) {
	return "", nil
}

func (n noOpProvider) UserOCID() (string, error) {
	return "", nil
}

func (n noOpProvider) KeyFingerprint() (string, error) {
	return "", nil
}

func (n noOpProvider) Region() (string, error) {
	return "", nil
}

func (n noOpProvider) AuthType() (common.AuthConfig, error) {
	return common.AuthConfig{AuthType: common.UnknownAuthenticationType}, nil
}
