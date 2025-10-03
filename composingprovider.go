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
	"github.com/oracle/oci-go-sdk/v65/common"
)

// ComposingConfigProvider takes a list of providers and return a [common.ConfigurationProvider]
// that returns the configuration of the first non-error result
//
// One key difference from [common.ComposingConfigurationProvider] is that this loops through all
// the providers in the list when determining AuthType.
func ComposingConfigProvider(providers ...common.ConfigurationProvider) common.ConfigurationProvider {
	p, _ := common.ComposingConfigurationProvider(providers)
	return &composingProvider{providers: providers, ConfigurationProvider: p}
}

type composingProvider struct {
	providers []common.ConfigurationProvider
	common.ConfigurationProvider
}

// AuthType replaces the method from [common.ComposingConfigurationProvider] which only checks the first provider in the list
func (c composingProvider) AuthType() (common.AuthConfig, error) {
	for _, provider := range c.providers {
		if authConfig, err := provider.AuthType(); err == nil && authConfig.AuthType != common.UnknownAuthenticationType {
			return authConfig, nil
		}
	}
	return common.AuthConfig{AuthType: common.UnknownAuthenticationType}, ErrNoAuthType
}
