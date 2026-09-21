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
//
// Identity and Resource are whatever tfsdk.ResourceIdentity.Set and tfsdk.Resource.Set
// accept: a struct with attr.Value fields for a framework-native resource, or a
// tftypes.Value for an SDKv2 one (see MappedItemFromResourceData).
type MappedItem struct {
	DisplayName string
	Identity    any
	Resource    any
}
