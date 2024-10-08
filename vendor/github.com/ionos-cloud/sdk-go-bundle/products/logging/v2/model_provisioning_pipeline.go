/*
 * IONOS Logging REST API
 *
 * The logging service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data.  This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the ProvisioningPipeline type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ProvisioningPipeline{}

// ProvisioningPipeline pipeline response
type ProvisioningPipeline struct {
	// The unique ID of the resource.
	Id         *string               `json:"id,omitempty"`
	Metadata   *ProvisioningMetadata `json:"metadata,omitempty"`
	Properties *PipelineProperties   `json:"properties,omitempty"`
}

// NewProvisioningPipeline instantiates a new ProvisioningPipeline object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewProvisioningPipeline() *ProvisioningPipeline {
	this := ProvisioningPipeline{}

	return &this
}

// NewProvisioningPipelineWithDefaults instantiates a new ProvisioningPipeline object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewProvisioningPipelineWithDefaults() *ProvisioningPipeline {
	this := ProvisioningPipeline{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *ProvisioningPipeline) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProvisioningPipeline) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *ProvisioningPipeline) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *ProvisioningPipeline) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *ProvisioningPipeline) GetMetadata() ProvisioningMetadata {
	if o == nil || IsNil(o.Metadata) {
		var ret ProvisioningMetadata
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProvisioningPipeline) GetMetadataOk() (*ProvisioningMetadata, bool) {
	if o == nil || IsNil(o.Metadata) {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *ProvisioningPipeline) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given ProvisioningMetadata and assigns it to the Metadata field.
func (o *ProvisioningPipeline) SetMetadata(v ProvisioningMetadata) {
	o.Metadata = &v
}

// GetProperties returns the Properties field value if set, zero value otherwise.
func (o *ProvisioningPipeline) GetProperties() PipelineProperties {
	if o == nil || IsNil(o.Properties) {
		var ret PipelineProperties
		return ret
	}
	return *o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ProvisioningPipeline) GetPropertiesOk() (*PipelineProperties, bool) {
	if o == nil || IsNil(o.Properties) {
		return nil, false
	}
	return o.Properties, true
}

// HasProperties returns a boolean if a field has been set.
func (o *ProvisioningPipeline) HasProperties() bool {
	if o != nil && !IsNil(o.Properties) {
		return true
	}

	return false
}

// SetProperties gets a reference to the given PipelineProperties and assigns it to the Properties field.
func (o *ProvisioningPipeline) SetProperties(v PipelineProperties) {
	o.Properties = &v
}

func (o ProvisioningPipeline) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsNil(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableProvisioningPipeline struct {
	value *ProvisioningPipeline
	isSet bool
}

func (v NullableProvisioningPipeline) Get() *ProvisioningPipeline {
	return v.value
}

func (v *NullableProvisioningPipeline) Set(val *ProvisioningPipeline) {
	v.value = val
	v.isSet = true
}

func (v NullableProvisioningPipeline) IsSet() bool {
	return v.isSet
}

func (v *NullableProvisioningPipeline) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableProvisioningPipeline(val *ProvisioningPipeline) *NullableProvisioningPipeline {
	return &NullableProvisioningPipeline{value: val, isSet: true}
}

func (v NullableProvisioningPipeline) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableProvisioningPipeline) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
