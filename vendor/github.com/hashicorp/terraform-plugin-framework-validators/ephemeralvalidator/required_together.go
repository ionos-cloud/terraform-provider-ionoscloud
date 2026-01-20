// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// RequiredTogether checks that a set of path.Expression either has all known
// or all null values.
func RequiredTogether(expressions ...path.Expression) ephemeral.ConfigValidator {
	return &configvalidator.RequiredTogetherValidator{
		PathExpressions: expressions,
	}
}
