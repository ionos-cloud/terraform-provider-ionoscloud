package writeonly

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// GetStringValue gets a write only attribute, checking that it is of an expected type and returns it
func GetStringValue(d *schema.ResourceData, name string) (string, error) {
	valueWO, err := GetValue(d, name, cty.String)
	if err != nil {
		return "", err
	}
	value := ""
	if !valueWO.IsNull() {
		value = valueWO.AsString()
	}
	return value, err
}

// GetValue gets a write only attribute, checking that it is of an expected type and subsequently returns it
func GetValue(d *schema.ResourceData, name string, attributeType cty.Type) (cty.Value, error) {
	if d.GetRawConfig().IsNull() {
		return cty.Value{}, fmt.Errorf("retrieving write-only attribute `%s`: resource data is null", name)
	}
	value, diags := d.GetRawConfigAt(cty.GetAttrPath(name))
	if diags.HasError() {
		return cty.Value{}, fmt.Errorf("retrieving write-only attribute `%s`: %+v", name, diags)
	}

	if !value.Type().Equals(attributeType) {
		return cty.Value{}, fmt.Errorf("retrieving write-only attribute `%s`: value is not of type %v", name, attributeType)
	}
	return value, nil
}
