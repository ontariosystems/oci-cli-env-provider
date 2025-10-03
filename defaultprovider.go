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
	"os"

	"github.com/ontariosystems/oci-cli-env-provider/internal"
	"github.com/oracle/oci-go-sdk/v65/common"
	"gopkg.in/ini.v1"
)

// DefaultConfigProvider returns a [common.ConfigurationProvider] containing providers for oci cli
// environment variables, as well as those returned by [common.DefaultConfigProvider]
func DefaultConfigProvider() common.ConfigurationProvider {
	var providers []common.ConfigurationProvider

	envProvider := OciCliEnvironmentConfigurationProvider()
	providers = append(providers, envProvider)

	configFilePath := os.Getenv(EnvConfigFile)
	if configFilePath == "" {
		configFilePath = internal.ExpandPath("~/.oci/config")
	}

	profileName := os.Getenv(EnvProfile)
	if profileName == "" {
		if cliConf, err := ini.Load(configFilePath); err == nil {
			profileName = cliConf.Section("OCI_CLI_SETTINGS").Key("default_profile").String()
		}
	}

	if profileName != "" {
		p, _ := common.ConfigurationProviderFromFileWithProfile(configFilePath, profileName, envProvider.(*ociCliEnvProvider).Passphrase())
		providers = append(providers, p)
	}

	providers = append(providers, common.DefaultConfigProvider())
	return ComposingConfigProvider(providers...)
}
