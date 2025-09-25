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
Package ocep creates a [common.ConfigurationProvider] which reads the
configuration specified by [oci-cli environment variables].

[oci-cli environment variables]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/clienvironmentvariables.htm
*/
package ocep

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

func (p *ociCliEnvProvider) Passphrase() string {
	return os.Getenv(EnvPassphrase)
}

func (p *ociCliEnvProvider) PrivateRSAKey() (*rsa.PrivateKey, error) {
	passphrase := p.Passphrase()

	if value, ok := os.LookupEnv(EnvKeyContent); ok {
		return common.PrivateKeyFromBytesWithPassword([]byte(value), []byte(passphrase))
	}

	if value, ok := os.LookupEnv(EnvKeyFile); ok {
		content, err := os.ReadFile(internal.ExpandPath(value))
		if err != nil {
			return nil, err
		}

		return common.PrivateKeyFromBytesWithPassword(content, []byte(passphrase))
	}

	return nil, errors.Join(&EnvError{EnvKeyContent}, &EnvError{EnvKeyFile})
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
		tokenPath, ok := os.LookupEnv(EnvSecurityTokenFile)
		if !ok {
			return "", &EnvError{EnvSecurityTokenFile}
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
	value, ok := os.LookupEnv(EnvTenancy)
	if !ok {
		return "", &EnvError{EnvTenancy}
	}
	return value, nil
}

func (p *ociCliEnvProvider) UserOCID() (string, error) {
	value, ok := os.LookupEnv(EnvUser)
	if !ok {
		return "", &EnvError{EnvUser}
	}
	return value, nil
}

func (p *ociCliEnvProvider) KeyFingerprint() (string, error) {
	value, ok := os.LookupEnv(EnvFingerprint)
	if !ok {
		return "", &EnvError{EnvFingerprint}
	}
	return value, nil
}

func (p *ociCliEnvProvider) Region() (string, error) {
	value, ok := os.LookupEnv(EnvRegion)
	if !ok {
		return "", &EnvError{EnvRegion}
	}
	return value, nil
}

func (p *ociCliEnvProvider) AuthType() (common.AuthConfig, error) {
	value, ok := os.LookupEnv(EnvAuth)
	if !ok {
		return common.AuthConfig{AuthType: common.UnknownAuthenticationType}, &EnvError{EnvAuth}
	}

	switch at := common.AuthenticationType(value); at {
	case ApiKeyType:
		return common.AuthConfig{AuthType: common.UserPrincipal}, nil
	default:
		return common.AuthConfig{AuthType: at}, nil
	}
}
