// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ephemeralvalidator

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/internal/configvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Conflicting checks that a set of path.Expression, are not configured
// simultaneously.
func Conflicting(expressions ...path.Expression) ephemeral.ConfigValidator {
	return &configvalidator.ConflictingValidator{
		PathExpressions: expressions,
	}
}
