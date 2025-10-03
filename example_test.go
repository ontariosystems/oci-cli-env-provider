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
	"context"
	"fmt"

	"github.com/ontariosystems/oci-cli-env-provider"
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/common/auth"
	"github.com/oracle/oci-go-sdk/v65/identity"
)

var (
	retryPolicy = common.DefaultRetryPolicy()
	metadata    = common.RequestMetadata{RetryPolicy: &retryPolicy}
)

// The most common use is to create a provider that can use the OCI CLI Environment variables as
// well as the providers from the sdk.
func Example() {
	provider := ocep.DefaultConfigProvider()
	client, _ := identity.NewIdentityClientWithConfigurationProvider(provider)

	tenancyID, _ := provider.TenancyOCID()
	req := identity.GetCompartmentRequest{
		CompartmentId:   common.String(tenancyID),
		RequestMetadata: metadata,
	}
	resp, _ := client.GetCompartment(context.TODO(), req)
	fmt.Printf("CompartmentId: %s\n", *resp.Id)
}

// This example requires the OCI CLI environment variables to be set.
// The output will be your Tenancy OCID.
func ExampleOciCliEnvironmentConfigurationProvider() {
	provider := ocep.OciCliEnvironmentConfigurationProvider()
	client, _ := identity.NewIdentityClientWithConfigurationProvider(provider)

	tenancyID, _ := provider.TenancyOCID()
	req := identity.GetCompartmentRequest{
		CompartmentId:   common.String(tenancyID),
		RequestMetadata: metadata,
	}
	resp, _ := client.GetCompartment(context.TODO(), req)
	fmt.Printf("CompartmentId: %s\n", *resp.Id)
}

// The [LazyConfigProvider] can be used to wrap providers that can only be created under
// certain circumstances (such as on a compute instance).
//
// Normally [auth.InstancePrincipalConfigurationProvider] will fail if not on a compute instance,
// but here it won't be called unless the other providers in the [ComposingConfigProvider] do not
// provide valid configuration
func ExampleLazyConfigProvider() {
	provider := ocep.ComposingConfigProvider(
		ocep.DefaultConfigProvider(),
		ocep.LazyConfigProvider(auth.InstancePrincipalConfigurationProvider),
	)
	client, _ := identity.NewIdentityClientWithConfigurationProvider(provider)

	tenancyID, _ := provider.TenancyOCID()
	req := identity.GetCompartmentRequest{
		CompartmentId:   common.String(tenancyID),
		RequestMetadata: metadata,
	}
	resp, _ := client.GetCompartment(context.TODO(), req)
	fmt.Printf("CompartmentId: %s\n", *resp.Id)
}
