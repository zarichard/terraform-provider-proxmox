/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package resource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zarichard/terraform-provider-proxmox/proxmoxtf/test"
)

// TestRoleInstantiation tests whether the Role instance can be instantiated.
func TestRoleInstantiation(t *testing.T) {
	t.Parallel()
	s := Role()

	if s == nil {
		t.Fatalf("Cannot instantiate Role")
	}
}

// TestRoleSchema tests the Role schema.
func TestRoleSchema(t *testing.T) {
	t.Parallel()
	s := Role()

	test.AssertRequiredArguments(t, s, []string{
		mkResourceVirtualEnvironmentRolePrivileges,
		mkResourceVirtualEnvironmentRoleRoleID,
	})

	test.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkResourceVirtualEnvironmentRolePrivileges: schema.TypeSet,
		mkResourceVirtualEnvironmentRoleRoleID:     schema.TypeString,
	})
}
