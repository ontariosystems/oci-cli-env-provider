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

/*
The standard [oci-cli environment variables] we use

[oci-cli environment variables]: https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/clienvironmentvariables.htm
*/
const (
	EnvAuth              = "OCI_CLI_AUTH"
	EnvConfigFile        = "OCI_CLI_CONFIG_FILE"
	EnvFingerprint       = "OCI_CLI_FINGERPRINT"
	EnvKeyContent        = "OCI_CLI_KEY_CONTENT"
	EnvKeyFile           = "OCI_CLI_KEY_FILE"
	EnvPassphrase        = "OCI_CLI_PASSPHRASE"
	EnvProfile           = "OCI_CLI_PROFILE"
	EnvRegion            = "OCI_CLI_REGION"
	EnvSecurityTokenFile = "OCI_CLI_SECURITY_TOKEN_FILE"
	EnvTenancy           = "OCI_CLI_TENANCY"
	EnvUser              = "OCI_CLI_USER"
)
