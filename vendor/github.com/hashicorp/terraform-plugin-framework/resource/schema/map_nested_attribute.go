// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwschema/fwxschema"
	"github.com/hashicorp/terraform-plugin-framework/internal/fwtype"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisifies the desired interfaces.
var (
	_ NestedAttribute                              = MapNestedAttribute{}
	_ fwschema.AttributeWithValidateImplementation = MapNestedAttribute{}
	_ fwschema.AttributeWithMapDefaultValue        = MapNestedAttribute{}
	_ fwxschema.AttributeWithMapPlanModifiers      = MapNestedAttribute{}
	_ fwxschema.AttributeWithMapValidators         = MapNestedAttribute{}
)

// MapNestedAttribute represents an attribute that is a map of objects where
// the object attributes can be fully defined, including further nested
// attributes. When retrieving the value for this attribute, use types.Map
// as the value type unless the CustomType field is set. The NestedObject field
// must be set. Nested attributes are only compatible with protocol version 6.
//
// Use MapAttribute if the underlying elements are of a single type and do
// not require definition beyond type information.
//
// Terraform configurations configure this attribute using expressions that
// return a map of objects or directly via curly brace syntax.
//
//	# map of objects
//	example_attribute = {
//		key = {
//			nested_attribute = #...
//		},
//	]
//
// Terraform configurations reference this attribute using expressions that
// accept a map of objects or an element directly via square brace string
// syntax:
//
//	# known object at key
//	.example_attribute["key"]
//	# known object nested_attribute value at key
//	.example_attribute["key"].nested_attribute
type MapNestedAttribute struct {
	// NestedObject is the underlying object that contains nested attributes.
	// This field must be set.
	//
	// Nested attributes that contain a dynamic type (i.e. DynamicAttribute) are not supported.
	// If underlying dynamic values are required, replace this attribute definition with
	// DynamicAttribute instead.
	NestedObject NestedAttributeObject

	// CustomType enables the use of a custom attribute type in place of the
	// default types.MapType of types.ObjectType. When retrieving data, the
	// basetypes.MapValuable associated with this custom type must be used in
	// place of types.Map.
	CustomType basetypes.MapTypable

	// Required indicates whether the practitioner must enter a value for
	// this attribute or not. Required and Optional cannot both be true,
	// and Required and Computed cannot both be true.
	Required bool

	// Optional indicates whether the practitioner can choose to enter a value
	// for this attribute or not. Optional and Required cannot both be true.
	Optional bool

	// Computed indicates whether the provider may return its own value for
	// this Attribute or not. Required and Computed cannot both be true. If
	// Required and Optional are both false, Computed must be true, and the
	// attribute will be considered "read only" for the practitioner, with
	// only the provider able to set its value.
	Computed bool

	// Sensitive indicates whether the value of this attribute should be
	// considered sensitive data. Setting it to true will obscure the value
	// in CLI output. Sensitive does not impact how values are stored, and
	// practitioners are encouraged to store their state as if the entire
	// file is sensitive.
	Sensitive bool

	// Description is used in various tooling, like the language server, to
	// give practitioners more information about what this attribute is,
	// what it's for, and how it should be used. It should be written as
	// plain text, with no special formatting.
	Description string

	// MarkdownDescription is used in various tooling, like the
	// documentation generator, to give practitioners more information
	// about what this attribute is, what it's for, and how it should be
	// used. It should be formatted using Markdown.
	MarkdownDescription string

	// DeprecationMessage defines warning diagnostic details to display when
	// practitioner configurations use this Attribute. The warning diagnostic
	// summary is automatically set to "Attribute Deprecated" along with
	// configuration source file and line information.
	//
	// Set this field to a practitioner actionable message such as:
	//
	//  - "Configure other_attribute instead. This attribute will be removed
	//    in the next major version of the provider."
	//  - "Remove this attribute's configuration as it no longer is used and
	//    the attribute will be removed in the next major version of the
	//    provider."
	//
	// In Terraform 1.2.7 and later, this warning diagnostic is displayed any
	// time a practitioner attempts to configure a value for this attribute and
	// certain scenarios where this attribute is referenced.
	//
	// In Terraform 1.2.6 and earlier, this warning diagnostic is only
	// displayed when the Attribute is Required or Optional, and if the
	// practitioner configuration sets the value to a known or unknown value
	// (which may eventually be null). It has no effect when the Attribute is
	// Computed-only (read-only; not Required or Optional).
	//
	// Across any Terraform version, there are no warnings raised for
	// practitioner configuration values set directly to null, as there is no
	// way for the framework to differentiate between an unset and null
	// configuration due to how Terraform sends configuration information
	// across the protocol.
	//
	// Additional information about deprecation enhancements for read-only
	// attributes can be found in:
	//
	//  - https://github.com/hashicorp/terraform/issues/7569
	//
	DeprecationMessage string

	// Validators define value validation functionality for the attribute. All
	// elements of the slice of AttributeValidator are run, regardless of any
	// previous error diagnostics.
	//
	// Many common use case validators can be found in the
	// github.com/hashicorp/terraform-plugin-framework-validators Go module.
	//
	// If the Type field points to a custom type that implements the
	// xattr.TypeWithValidate interface, the validators defined in this field
	// are run in addition to the validation defined by the type.
	Validators []validator.Map

	// PlanModifiers defines a sequence of modifiers for this attribute at
	// plan time. Schema-based plan modifications occur before any
	// resource-level plan modifications.
	//
	// Schema-based plan modifications can adjust Terraform's plan by:
	//
	//  - Requiring resource recreation. Typically used for configuration
	//    updates which cannot be done in-place.
	//  - Setting the planned value. Typically used for enhancing the plan
	//    to replace unknown values. Computed must be true or Terraform will
	//    return an error. If the plan value is known due to a known
	//    configuration value, the plan value cannot be changed or Terraform
	//    will return an error.
	//
	// Any errors will prevent further execution of this sequence or modifiers.
	PlanModifiers []planmodifier.Map

	// Default defines a proposed new state (plan) value for the attribute
	// if the configuration value is null. Default prevents the framework
	// from automatically marking the value as unknown during planning when
	// other proposed new state changes are detected. If the attribute is
	// computed and the value could be altered by other changes then a default
	// should be avoided and a plan modifier should be used instead.
	Default defaults.Map

	// WriteOnly indicates that Terraform will not store this attribute value
	// in the plan or state artifacts.
	// If WriteOnly is true, either Optional or Required must also be true.
	// WriteOnly cannot be set with Computed.
	//
	// If WriteOnly is true for a nested attribute, all of its child attributes
	// must also set WriteOnly to true and no child attribute can be Computed.
	//
	// This functionality is only supported in Terraform 1.11 and later.
	// Practitioners that choose a value for this attribute with older
	// versions of Terraform will receive an error.
	WriteOnly bool
}

// ApplyTerraform5AttributePathStep returns the Attributes field value if step
// is ElementKeyString, otherwise returns an error.
func (a MapNestedAttribute) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	_, ok := step.(tftypes.ElementKeyString)

	if !ok {
		return nil, fmt.Errorf("cannot apply step %T to MapNestedAttribute", step)
	}

	return a.NestedObject, nil
}

// Equal returns true if the given Attribute is a MapNestedAttribute
// and all fields are equal.
func (a MapNestedAttribute) Equal(o fwschema.Attribute) bool {
	other, ok := o.(MapNestedAttribute)

	if !ok {
		return false
	}

	return fwschema.NestedAttributesEqual(a, other)
}

// GetDeprecationMessage returns the DeprecationMessage field value.
func (a MapNestedAttribute) GetDeprecationMessage() string {
	return a.DeprecationMessage
}

// GetDescription returns the Description field value.
func (a MapNestedAttribute) GetDescription() string {
	return a.Description
}

// GetMarkdownDescription returns the MarkdownDescription field value.
func (a MapNestedAttribute) GetMarkdownDescription() string {
	return a.MarkdownDescription
}

// GetNestedObject returns the NestedObject field value.
func (a MapNestedAttribute) GetNestedObject() fwschema.NestedAttributeObject {
	return a.NestedObject
}

// GetNestingMode always returns NestingModeMap.
func (a MapNestedAttribute) GetNestingMode() fwschema.NestingMode {
	return fwschema.NestingModeMap
}

// GetType returns MapType of ObjectType or CustomType.
func (a MapNestedAttribute) GetType() attr.Type {
	if a.CustomType != nil {
		return a.CustomType
	}

	return types.MapType{
		ElemType: a.NestedObject.Type(),
	}
}

// IsComputed returns the Computed field value.
func (a MapNestedAttribute) IsComputed() bool {
	return a.Computed
}

// IsOptional returns the Optional field value.
func (a MapNestedAttribute) IsOptional() bool {
	return a.Optional
}

// IsRequired returns the Required field value.
func (a MapNestedAttribute) IsRequired() bool {
	return a.Required
}

// IsSensitive returns the Sensitive field value.
func (a MapNestedAttribute) IsSensitive() bool {
	return a.Sensitive
}

// IsWriteOnly returns the WriteOnly field value.
func (a MapNestedAttribute) IsWriteOnly() bool {
	return a.WriteOnly
}

// MapDefaultValue returns the Default field value.
func (a MapNestedAttribute) MapDefaultValue() defaults.Map {
	return a.Default
}

// MapPlanModifiers returns the PlanModifiers field value.
func (a MapNestedAttribute) MapPlanModifiers() []planmodifier.Map {
	return a.PlanModifiers
}

// MapValidators returns the Validators field value.
func (a MapNestedAttribute) MapValidators() []validator.Map {
	return a.Validators
}

// ValidateImplementation contains logic for validating the
// provider-defined implementation of the attribute to prevent unexpected
// errors or panics. This logic runs during the GetProviderSchema RPC and
// should never include false positives.
func (a MapNestedAttribute) ValidateImplementation(ctx context.Context, req fwschema.ValidateImplementationRequest, resp *fwschema.ValidateImplementationResponse) {
	if a.CustomType == nil && fwtype.ContainsCollectionWithDynamic(a.GetType()) {
		resp.Diagnostics.Append(fwtype.AttributeCollectionWithDynamicTypeDiag(req.Path))
	}

	if a.IsWriteOnly() && !fwschema.ContainsAllWriteOnlyChildAttributes(a) {
		resp.Diagnostics.Append(fwschema.InvalidWriteOnlyNestedAttributeDiag(req.Path))
	}

	if a.IsComputed() && fwschema.ContainsAnyWriteOnlyChildAttributes(a) {
		resp.Diagnostics.Append(fwschema.InvalidComputedNestedAttributeWithWriteOnlyDiag(req.Path))
	}

	if a.MapDefaultValue() != nil {
		if !a.IsComputed() {
			resp.Diagnostics.Append(nonComputedAttributeWithDefaultDiag(req.Path))
		}

		// Validate Default implementation. This is safe unless the framework
		// ever allows more dynamic Default implementations at which the
		// implementation would be required to be validated at runtime.
		// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/930
		defaultReq := defaults.MapRequest{
			Path: req.Path,
		}
		defaultResp := &defaults.MapResponse{}

		a.MapDefaultValue().DefaultMap(ctx, defaultReq, defaultResp)

		resp.Diagnostics.Append(defaultResp.Diagnostics...)

		if defaultResp.Diagnostics.HasError() {
			return
		}

		if a.CustomType == nil && a.NestedObject.CustomType == nil && !a.NestedObject.Type().Equal(defaultResp.PlanValue.ElementType(ctx)) {
			resp.Diagnostics.Append(fwschema.AttributeDefaultElementTypeMismatchDiag(req.Path, a.NestedObject.Type(), defaultResp.PlanValue.ElementType(ctx)))
		}
	}
}
