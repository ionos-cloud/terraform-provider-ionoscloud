package identity

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Model represents a reusable identity model for Terraform Plugin Framework resources
// that use a single "id" attribute in their Identity Schema.
type Model struct {
	ID types.String `tfsdk:"id"`
}

// MappedItem is the value returned by a mapper function. StreamList uses it to
// populate the list result — the mapper itself never mutates *list.ListResult directly.
type MappedItem struct {
	DisplayName string
	Identity    *Model
	Resource    any
}
