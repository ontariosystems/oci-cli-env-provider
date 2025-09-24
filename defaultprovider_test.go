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

package oci_cli_env_provider

import (
	"bytes"
	"html/template"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/oracle/oci-go-sdk/v65/common"
)

var _ = Describe("DefaultConfigProvider", func() {
	var (
		privateKeyPath string
		conf           = DefaultConfigProvider()
	)

	AfterEach(func() {
		_ = os.Remove(privateKeyPath)
	})

	Context("valid configuration", func() {
		BeforeEach(func() {
			_ = os.Setenv(EnvUser, testUser)
			_ = os.Setenv(EnvTenancy, testTenancy)
			_ = os.Setenv(EnvFingerprint, testFingerprint)
			_ = os.Setenv(EnvRegion, testRegion)
			privateKeyPath = createTempFile(testPrivateKeyConf)
			_ = os.Setenv(EnvKeyFile, privateKeyPath)
			_ = os.Setenv(EnvAuth, string(ApiKeyType))
		})

		It("has valid configuration", func() {
			valid, err := common.IsConfigurationProviderValid(conf)
			Expect(err).ToNot(HaveOccurred())
			Expect(valid).To(BeTrue())

			at, err := conf.AuthType()
			Expect(err).ToNot(HaveOccurred())
			Expect(at.AuthType).To(Equal(common.UserPrincipal))
		})

		Context("with cli config file", func() {
			var configFile string

			BeforeEach(func() {
				for _, env := range []string{EnvUser, EnvTenancy, EnvFingerprint, EnvRegion, EnvKeyFile, EnvAuth} {
					_ = os.Unsetenv(env)
				}

				testCliConfigFileTmplData.TestKeyFile = privateKeyPath

				b := &bytes.Buffer{}
				_ = template.Must(template.New("").Parse(testCliConfigFileTmpl)).Execute(b, testCliConfigFileTmplData)
				configFile = createTempFile(b.Bytes())
				_ = os.Setenv(EnvConfigFile, configFile)

				conf = DefaultConfigProvider()
			})

			AfterEach(func() {
				_ = os.Remove(configFile)
			})

			It("has valid configuration", func() {
				valid, err := common.IsConfigurationProviderValid(conf)
				Expect(err).ToNot(HaveOccurred())
				Expect(valid).To(BeTrue())

				u, err := conf.UserOCID()
				Expect(err).ToNot(HaveOccurred())
				Expect(u).To(Equal(testUser))
			})

			Context("with default cli config file", func() {
				var testHomePath string

				BeforeEach(func() {
					_ = os.Unsetenv(EnvConfigFile)
					testHomePath = createTempDir()
					_ = os.Mkdir(path.Join(testHomePath, ".oci"), 0700)
					cliConf, _ := os.ReadFile(configFile)
					_ = os.WriteFile(path.Join(testHomePath, ".oci", "config"), cliConf, 0600)
					_ = os.Setenv("HOME", testHomePath)

					conf = DefaultConfigProvider()
				})

				AfterEach(func() {
					_ = os.RemoveAll(testHomePath)
				})

				It("has valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).ToNot(HaveOccurred())
					Expect(valid).To(BeTrue())

					u, err := conf.UserOCID()
					Expect(err).ToNot(HaveOccurred())
					Expect(u).To(Equal(testUser))
				})
			})

			When("cli config profile environment variable is set", func() {
				var tokenFilePath string

				BeforeEach(func() {
					_ = os.Remove(configFile)

					tokenFilePath = createTempFile([]byte(testSecurityToken))
					testCliConfigFileTmplData.AltSecurityTokenFile = tokenFilePath
					testCliConfigFileTmplData.AltKeyFile = privateKeyPath
					_ = os.Setenv(EnvProfile, "alt")
					_ = os.Setenv(EnvAuth, string(SecurityTokenType))

					b := &bytes.Buffer{}
					_ = template.Must(template.New("").Parse(testCliConfigFileTmpl)).Execute(b, testCliConfigFileTmplData)
					configFile = createTempFile(b.Bytes())
					_ = os.Setenv(EnvConfigFile, configFile)

					conf = DefaultConfigProvider()
				})

				AfterEach(func() {
					_ = os.Remove(tokenFilePath)
				})

				It("has valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).ToNot(HaveOccurred())
					Expect(valid).To(BeTrue())

					Expect(conf.KeyID()).To(ContainSubstring(testSecurityToken))
					at, err := conf.AuthType()
					Expect(err).ToNot(HaveOccurred())
					Expect(at.AuthType).To(Equal(SecurityTokenType))
				})
			})
		})
	})
})

func createTempDir() string {
	tmp, _ := os.MkdirTemp("", "ociclienvprovider")
	return tmp
}

var (
	testCliConfigFileTmpl = `
[OCI_CLI_SETTINGS]
default_profile = test

[DEFAULT]

[test]
fingerprint = {{ .TestFingerprint }}
user = {{ .TestUser }}
key_file = {{ .TestKeyFile }}
tenancy = {{ .TestTenancy }}
region = {{ .TestRegion }}

[alt]
fingerprint = {{ .AltFingerprint }}
key_file = {{ .AltKeyFile }}
tenancy = {{ .AltTenancy }}
region = {{ .AltRegion }}
security_token_file = {{ .AltSecurityTokenFile }}
`
	testCliConfigFileTmplData = struct {
		TestFingerprint, TestUser, TestKeyFile, TestTenancy, TestRegion         string
		AltFingerprint, AltKeyFile, AltTenancy, AltRegion, AltSecurityTokenFile string
	}{
		TestFingerprint: testFingerprint,
		TestUser:        testUser,
		TestTenancy:     testTenancy,
		TestRegion:      testRegion,
		AltFingerprint:  "alt-" + testFingerprint,
		AltTenancy:      "alt-" + testTenancy,
		AltRegion:       "alt-" + testRegion,
	}
)
