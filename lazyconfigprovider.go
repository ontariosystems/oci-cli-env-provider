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

package ocep

import (
	"crypto/rsa"

	"github.com/oracle/oci-go-sdk/v65/common"
)

// LazyConfigProvider returns a [common.ConfigurationProvider] that is initialized one time
// by calling the func argument. The initialization func is only called if the provider methods are
// called.
func LazyConfigProvider(providerFunc func() (common.ConfigurationProvider, error)) common.ConfigurationProvider {
	return &lazyProvider{providerFunc: providerFunc}
}

type lazyProvider struct {
	providerFunc func() (common.ConfigurationProvider, error)
	common.ConfigurationProvider
}

func (p *lazyProvider) initProvider() (err error) {
	if p.ConfigurationProvider == nil {
		if p.ConfigurationProvider, err = p.providerFunc(); err != nil {
			p.ConfigurationProvider = nil
		}
	}
	return
}

func (p *lazyProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	if err := p.initProvider(); err != nil {
		return nil, err
	}
	return p.ConfigurationProvider.PrivateRSAKey()
}

func (p *lazyProvider) KeyID() (string, error) {
	if err := p.initProvider(); err != nil {
		return "", err
	}
	return p.ConfigurationProvider.KeyID()
}

func (p *lazyProvider) TenancyOCID() (string, error) {
	if err := p.initProvider(); err != nil {
		return "", err
	}
	return p.ConfigurationProvider.TenancyOCID()
}

func (p *lazyProvider) UserOCID() (string, error) {
	if err := p.initProvider(); err != nil {
		return "", err
	}
	return p.ConfigurationProvider.UserOCID()
}

func (p *lazyProvider) KeyFingerprint() (string, error) {
	if err := p.initProvider(); err != nil {
		return "", err
	}
	return p.ConfigurationProvider.KeyFingerprint()
}

func (p *lazyProvider) Region() (string, error) {
	if err := p.initProvider(); err != nil {
		return "", err
	}
	return p.ConfigurationProvider.Region()
}

func (p *lazyProvider) AuthType() (common.AuthConfig, error) {
	if err := p.initProvider(); err != nil {
		return common.AuthConfig{AuthType: common.UnknownAuthenticationType}, err
	}
	return p.ConfigurationProvider.AuthType()
}
