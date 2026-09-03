// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package translate

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-mux/internal/tfprotov5tov6"
)

func Schema(in *tfprotov5.Schema) *tfprotov6.Schema {
	return tfprotov5tov6.Schema(in)
}

func ResourceIdentitySchema(in *tfprotov5.ResourceIdentitySchema) *tfprotov6.ResourceIdentitySchema {
	return tfprotov5tov6.ResourceIdentitySchema(in)
}
