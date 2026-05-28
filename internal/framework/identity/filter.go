package identity

import (
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// FiltersKey is the standard config key used by list resources for filters.
const FiltersKey = "filters"

// Filter represents a single filter criterion for list resources.
type Filter struct {
	FieldName  types.String `tfsdk:"field_name"`
	FieldValue types.String `tfsdk:"field_value"`
}

// FilterAttribute returns the standard "filters" list attribute for list resource config schemas.
func FilterAttribute() listschema.ListNestedAttribute {
	return listschema.ListNestedAttribute{
		Optional:    true,
		Description: "Filters to apply when listing resources. All filters must match (AND logic).",
		NestedObject: listschema.NestedAttributeObject{
			Attributes: map[string]listschema.Attribute{
				"field_name":  listschema.StringAttribute{Required: true, Description: "The name of the field to filter on."},
				"field_value": listschema.StringAttribute{Required: true, Description: "The value to match against."},
			},
		},
	}
}

// MatchesFilters returns true if all filters match the given fields map.
// An empty or nil filters slice always returns true.
func MatchesFilters(fields map[string]string, filters []Filter) bool {
	for _, f := range filters {
		val, ok := fields[f.FieldName.ValueString()]
		if !ok || val != f.FieldValue.ValueString() {
			return false
		}
	}
	return true
}
