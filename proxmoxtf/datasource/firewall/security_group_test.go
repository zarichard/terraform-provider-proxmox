/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package firewall

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/zarichard/terraform-provider-proxmox/proxmoxtf/structure"
)

// TestSecurityGroupSchemaInstantiation tests whether the SecurityGroupSchema instance can be instantiated.
func TestSecurityGroupSchemaInstantiation(t *testing.T) {
	t.Parallel()
	require.NotNil(t, SecurityGroupSchema(), "Cannot instantiate SecurityGroupSchema")
}

// TestSecurityGroupSchema tests the SecurityGroupSchema.
func TestSecurityGroupSchema(t *testing.T) {
	t.Parallel()
	s := SecurityGroupSchema()

	structure.AssertRequiredArguments(t, s, []string{
		mkSecurityGroupName,
	})

	structure.AssertComputedAttributes(t, s, []string{
		mkSecurityGroupComment,
		mkRules,
	})

	structure.AssertValueTypes(t, s, map[string]schema.ValueType{
		mkSecurityGroupName:    schema.TypeString,
		mkSecurityGroupComment: schema.TypeString,
		mkRules:                schema.TypeList,
	})
}
