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

/*
Package oci_cli_env_provider creates a [common.ConfigurationProvider] which reads the
configuration specified by [oci-cli environment variables].

[oci-cli environment variables]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/clienvironmentvariables.htm
*/
package oci_cli_env_provider

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"

	"github.com/ontariosystems/oci-cli-env-provider/internal"
	"github.com/oracle/oci-go-sdk/v65/common"
)

const (
	ApiKeyType        common.AuthenticationType = "api_key"
	SecurityTokenType common.AuthenticationType = "security_token"
)

// OciCliEnvironmentConfigurationProvider returns a [common.ConfigurationProvider] that
// gets values from [oci-cli environment variables]
//
// [oci-cli environment variables]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/clienvironmentvariables.htm
func OciCliEnvironmentConfigurationProvider() common.ConfigurationProvider {
	return &ociCliEnvProvider{}
}

type ociCliEnvProvider struct{}

func (p *ociCliEnvProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	passphrase := os.Getenv("OCI_CLI_PASSPHRASE")

	if value, ok := os.LookupEnv("OCI_CLI_KEY_CONTENT"); ok {
		return common.PrivateKeyFromBytesWithPassword([]byte(value), []byte(passphrase))
	}

	if value, ok := os.LookupEnv("OCI_CLI_KEY_FILE"); ok {
		content, err := os.ReadFile(internal.ExpandPath(value))
		if err != nil {
			return nil, err
		}

		return common.PrivateKeyFromBytesWithPassword(content, []byte(passphrase))
	}

	return nil, errors.Join(&EnvError{"OCI_CLI_KEY_CONTENT"}, &EnvError{"OCI_CLI_KEY_FILE"})
}

func (p *ociCliEnvProvider) KeyID() (keyID string, err error) {
	tenancy, err := p.TenancyOCID()
	if err != nil {
		return
	}

	fingerprint, err := p.KeyFingerprint()
	if err != nil {
		return
	}

	user, err := p.UserOCID()
	if err == nil {
		return fmt.Sprintf("%s/%s/%s", tenancy, user, fingerprint), nil
	}

	at, err := p.AuthType()
	if err != nil {
		return
	}
	switch at.AuthType {
	case SecurityTokenType:
		tokenPath, ok := os.LookupEnv("OCI_CLI_SECURITY_TOKEN_FILE")
		if !ok {
			return "", &EnvError{"OCI_CLI_SECURITY_TOKEN_FILE"}
		}

		var token []byte
		if token, err = os.ReadFile(internal.ExpandPath(tokenPath)); err != nil {
			return
		}
		keyID = fmt.Sprintf("ST$%s", token)
		return
	}

	err = ErrNoKeyId
	return
}

func (p *ociCliEnvProvider) TenancyOCID() (string, error) {
	value, ok := os.LookupEnv("OCI_CLI_TENANCY")
	if !ok {
		return "", &EnvError{"OCI_CLI_TENANCY"}
	}
	return value, nil
}

func (p *ociCliEnvProvider) UserOCID() (string, error) {
	value, ok := os.LookupEnv("OCI_CLI_USER")
	if !ok {
		return "", &EnvError{"OCI_CLI_USER"}
	}
	return value, nil
}

func (p *ociCliEnvProvider) KeyFingerprint() (string, error) {
	value, ok := os.LookupEnv("OCI_CLI_FINGERPRINT")
	if !ok {
		return "", &EnvError{"OCI_CLI_FINGERPRINT"}
	}
	return value, nil
}

func (p *ociCliEnvProvider) Region() (string, error) {
	value, ok := os.LookupEnv("OCI_CLI_REGION")
	if !ok {
		return "", &EnvError{"OCI_CLI_REGION"}
	}
	return value, nil
}

func (p *ociCliEnvProvider) AuthType() (common.AuthConfig, error) {
	value, ok := os.LookupEnv("OCI_CLI_AUTH")
	if !ok {
		return common.AuthConfig{AuthType: common.UnknownAuthenticationType}, &EnvError{"OCI_CLI_AUTH"}
	}

	switch at := common.AuthenticationType(value); at {
	case ApiKeyType:
		return common.AuthConfig{AuthType: common.UserPrincipal}, nil
	default:
		return common.AuthConfig{AuthType: at}, nil
	}
}
