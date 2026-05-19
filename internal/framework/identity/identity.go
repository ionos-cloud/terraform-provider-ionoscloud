package identity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents a reusable identity model for Terraform Plugin Framework resources
// that use a single "id" attribute in their Identity Schema.
type Model struct {
	ID types.String `tfsdk:"id"`
}
