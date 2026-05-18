package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Identity represents a reusable identity model for Terraform Plugin Framework resources
// that use a single "id" attribute in their Identity Schema.
type Identity struct {
	ID types.String `tfsdk:"id"`
}
