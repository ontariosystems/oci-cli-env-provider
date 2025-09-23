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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/oracle/oci-go-sdk/v65/common"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OCI CLI Env Provider Suite")
}

var _ = Describe("OciCliEnvironmentConfigurationProvider", func() {
	var (
		privateKeyPath string
		conf           = OciCliEnvironmentConfigurationProvider()
	)

	BeforeEach(func() {
		for _, env := range os.Environ() {
			key := strings.SplitN(env, "=", 2)[0]
			if strings.HasPrefix(key, "OCI_") {
				_ = os.Unsetenv(key)
			}
		}
	})

	AfterEach(func() {
		_ = os.Remove(privateKeyPath)
	})

	Context("valid configuration", func() {
		BeforeEach(func() {
			_ = os.Setenv("OCI_CLI_USER", testUser)
			_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
			_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
			_ = os.Setenv("OCI_CLI_REGION", testRegion)
		})

		When("environment variables are set for api key", func() {
			BeforeEach(func() {
				privateKeyPath = createTempFile(testPrivateKeyConf)
				_ = os.Setenv("OCI_CLI_KEY_FILE", privateKeyPath)
				_ = os.Setenv("OCI_CLI_AUTH", string(ApiKeyType))
			})

			It("has valid configuration", func() {
				valid, err := common.IsConfigurationProviderValid(conf)
				Expect(err).ToNot(HaveOccurred())
				Expect(valid).To(BeTrue())

				at, err := conf.AuthType()
				Expect(err).ToNot(HaveOccurred())
				Expect(at.AuthType).To(Equal(common.UserPrincipal))
			})

			Context("with key content environment variable", func() {
				BeforeEach(func() {
					_ = os.Unsetenv("OCI_CLI_KEY_FILE")
					_ = os.Setenv("OCI_CLI_KEY_CONTENT", string(testPrivateKeyConf))
				})

				It("has valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).ToNot(HaveOccurred())
					Expect(valid).To(BeTrue())
				})
			})

			Context("with encrypted key file", func() {
				BeforeEach(func() {
					privateKeyPath = createTempFile(testEncryptedPrivateKeyConf)
					_ = os.Setenv("OCI_CLI_KEY_FILE", privateKeyPath)
					_ = os.Setenv("OCI_CLI_PASSPHRASE", testPassphrase)
				})
				It("has valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).ToNot(HaveOccurred())
					Expect(valid).To(BeTrue())
				})
			})

			Context("with encrypted key content", func() {
				BeforeEach(func() {
					_ = os.Unsetenv("OCI_CLI_KEY_FILE")
					_ = os.Setenv("OCI_CLI_KEY_CONTENT", string(testEncryptedPrivateKeyConf))
					_ = os.Setenv("OCI_CLI_PASSPHRASE", testPassphrase)
				})
				It("has valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).ToNot(HaveOccurred())
					Expect(valid).To(BeTrue())
				})
			})
		})

		When("environment variables are set for sso token", func() {
			var tokenFilePath string
			BeforeEach(func() {
				privateKeyPath = createTempFile(testPrivateKeyConf)
				_ = os.Setenv("OCI_CLI_KEY_FILE", privateKeyPath)
				_ = os.Setenv("OCI_CLI_AUTH", string(SecurityTokenType))
				_ = os.Unsetenv("OCI_CLI_USER")
				tokenFilePath = createTempFile([]byte(testSecurityToken))
				_ = os.Setenv("OCI_CLI_SECURITY_TOKEN_FILE", tokenFilePath)
			})
			AfterEach(func() {
				_ = os.Remove(tokenFilePath)
			})

			It("has valid configuration", func() {
				key, err := conf.KeyID()
				Expect(err).ToNot(HaveOccurred())
				Expect(key).To(ContainSubstring(testSecurityToken))
				Expect(key).To(HavePrefix("ST$"))
			})
		})
	})

	Context("invalid configuration", func() {
		Context("invalid private key", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_USER", testUser)
				_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
				_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
				_ = os.Setenv("OCI_CLI_REGION", testRegion)
			})

			When("neither key content nor file set", func() {
				It("does not have valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).To(HaveOccurred())
					Expect(valid).To(BeFalse())
				})
			})

			When("key path is invalid", func() {
				BeforeEach(func() {
					_ = os.Setenv("OCI_CLI_KEY_FILE", "/does/not/exist")
				})
				It("does not have valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).To(HaveOccurred())
					Expect(valid).To(BeFalse())
				})
			})
		})

		Context("invalid tenancy", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_USER", testUser)
				_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
				_ = os.Setenv("OCI_CLI_REGION", testRegion)
			})
			It("does not have valid configuration", func() {
				valid, err := common.IsConfigurationProviderValid(conf)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("TENANCY"))
				Expect(valid).To(BeFalse())

				_, err = conf.KeyID()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("TENANCY"))
			})
		})

		Context("invalid fingerprint", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_USER", testUser)
				_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
				_ = os.Setenv("OCI_CLI_REGION", testRegion)
			})
			It("does not have valid configuration", func() {
				valid, err := common.IsConfigurationProviderValid(conf)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("FINGERPRINT"))
				Expect(valid).To(BeFalse())

				_, err = conf.KeyID()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("FINGERPRINT"))
			})
		})

		Context("invalid region", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_USER", testUser)
				_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
				_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
			})
			It("does not have valid configuration", func() {
				valid, err := common.IsConfigurationProviderValid(conf)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("REGION"))
				Expect(valid).To(BeFalse())
			})

			When("auth type is not sso", func() {
				BeforeEach(func() {
					_ = os.Setenv("OCI_CLI_AUTH", string(common.InstancePrincipal))
					_ = os.Unsetenv("OCI_CLI_REGION")
					_ = os.Unsetenv("OCI_CLI_USER")
				})

				It("does not have valid configuration", func() {
					_, err := conf.KeyID()
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ErrNoKeyId))
				})
			})
		})

		Context("no user and invalid auth type", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
				_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
				_ = os.Setenv("OCI_CLI_REGION", testRegion)
			})
			It("does not have valid configuration", func() {
				_, err := conf.KeyID()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("AUTH"))
			})
		})

		Context("invalid security token", func() {
			BeforeEach(func() {
				_ = os.Setenv("OCI_CLI_TENANCY", testTenancy)
				_ = os.Setenv("OCI_CLI_FINGERPRINT", testFingerprint)
				_ = os.Setenv("OCI_CLI_REGION", testRegion)
				_ = os.Setenv("OCI_CLI_AUTH", string(SecurityTokenType))
			})

			When("the token path is not set", func() {
				It("does not have valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).To(HaveOccurred())
					Expect(valid).To(BeFalse())

					_, err = conf.KeyID()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("SECURITY_TOKEN"))
				})
			})

			When("the token path is invalid", func() {
				BeforeEach(func() {
					_ = os.Setenv("OCI_CLI_SECURITY_TOKEN_FILE", "/does/not/exist")
				})
				It("does not have valid configuration", func() {
					valid, err := common.IsConfigurationProviderValid(conf)
					Expect(err).To(HaveOccurred())
					Expect(valid).To(BeFalse())

					_, err = conf.KeyID()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("no such file"))
				})
			})
		})
	})

})

func createTempFile(content []byte) string {
	f, _ := os.CreateTemp("", "ociclienvprovider")
	defer func() { _ = f.Close() }()
	_, _ = f.Write(content)
	return f.Name()
}

var (
	testUser          = "test-user"
	testFingerprint   = "test-fingerprint"
	testTenancy       = "test-tenancy"
	testRegion        = "test-region"
	testSecurityToken = "test-security-token"

	testPk, _          = rsa.GenerateKey(rand.Reader, 4096)
	testPkBlock        = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(testPk)}
	testPrivateKeyConf = pem.EncodeToMemory(testPkBlock)

	testPassphrase              = "super-secret-passphrase"
	testEncryptedPkBlock, _     = x509.EncryptPEMBlock(rand.Reader, testPkBlock.Type, testPkBlock.Bytes, []byte(testPassphrase), x509.PEMCipherAES256)
	testEncryptedPrivateKeyConf = pem.EncodeToMemory(testEncryptedPkBlock)
)
