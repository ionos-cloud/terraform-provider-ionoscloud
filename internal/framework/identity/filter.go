package identity

import (
	"slices"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
func FilterAttribute(allowedFields ...string) listschema.ListNestedAttribute {
	fieldNameAttr := listschema.StringAttribute{Required: true, Description: "The name of the field to filter on."}
	if len(allowedFields) > 0 {
		fieldNameAttr.Validators = []validator.String{stringvalidator.OneOf(allowedFields...)}
	}

	return listschema.ListNestedAttribute{
		Optional:    true,
		Description: "Filters to apply when listing resources. All filters must match (AND logic).",
		NestedObject: listschema.NestedAttributeObject{
			Attributes: map[string]listschema.Attribute{
				"field_name":  fieldNameAttr,
				"field_value": listschema.StringAttribute{Required: true, Description: "The value to match against."},
			},
		},
	}
}

// FilterValue returns the field_value of the first filter whose field_name matches fieldName,
// or "" if no such filter exists. Useful for pushing a filter down to an API call.
func FilterValue(filters []Filter, fieldName string) string {
	idx := slices.IndexFunc(filters, func(f Filter) bool {
		return f.FieldName.ValueString() == fieldName
	})
	if idx < 0 {
		return ""
	}
	return filters[idx].FieldValue.ValueString()
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
