/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zarichard/terraform-provider-proxmox/proxmoxtf/test"
)

// TestProviderInstantiation() tests whether the ProxmoxVirtualEnvironment instance can be instantiated.
func TestProviderInstantiation(t *testing.T) {
	t.Parallel()
	s := ProxmoxVirtualEnvironment()

	if s == nil {
		t.Fatalf("Cannot instantiate ProxmoxVirtualEnvironment")
	}
}

// TestProviderSchema() tests the ProxmoxVirtualEnvironment schema.
func TestProviderSchema(t *testing.T) {
	t.Parallel()
	s := &schema.Resource{
		Schema: ProxmoxVirtualEnvironment().Schema,
	}

	test.AssertOptionalArguments(t, s, []string{
		mkProviderVirtualEnvironment,
		mkProviderUsername,
		mkProviderPassword,
		mkProviderEndpoint,
		mkProviderInsecure,
		mkProviderOTP,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkProviderVirtualEnvironment: schema.TypeList,
		mkProviderUsername:           schema.TypeString,
		mkProviderPassword:           schema.TypeString,
		mkProviderEndpoint:           schema.TypeString,
		mkProviderInsecure:           schema.TypeBool,
		mkProviderOTP:                schema.TypeString,
	})

	veSchema := test.AssertNestedSchemaExistence(t, s, mkProviderVirtualEnvironment)

	test.AssertOptionalArguments(t, veSchema, []string{
		mkProviderEndpoint,
		mkProviderInsecure,
		mkProviderOTP,
		mkProviderPassword,
		mkProviderUsername,
	})

	test.AssertValueTypes(t, veSchema, map[string]schema.ValueType{
		mkProviderEndpoint: schema.TypeString,
		mkProviderInsecure: schema.TypeBool,
		mkProviderOTP:      schema.TypeString,
		mkProviderPassword: schema.TypeString,
		mkProviderUsername: schema.TypeString,
	})
}
