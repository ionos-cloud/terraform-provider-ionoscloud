/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the DistributionCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DistributionCreate{}

// DistributionCreate struct for DistributionCreate
type DistributionCreate struct {
	// Metadata
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Properties DistributionProperties `json:"properties"`
}

// NewDistributionCreate instantiates a new DistributionCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDistributionCreate(properties DistributionProperties) *DistributionCreate {
	this := DistributionCreate{}

	this.Properties = properties

	return &this
}

// NewDistributionCreateWithDefaults instantiates a new DistributionCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDistributionCreateWithDefaults() *DistributionCreate {
	this := DistributionCreate{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *DistributionCreate) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DistributionCreate) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *DistributionCreate) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *DistributionCreate) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *DistributionCreate) GetProperties() DistributionProperties {
	if o == nil {
		var ret DistributionProperties
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *DistributionCreate) GetPropertiesOk() (*DistributionProperties, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *DistributionCreate) SetProperties(v DistributionProperties) {
	o.Properties = v
}

func (o DistributionCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableDistributionCreate struct {
	value *DistributionCreate
	isSet bool
}

func (v NullableDistributionCreate) Get() *DistributionCreate {
	return v.value
}

func (v *NullableDistributionCreate) Set(val *DistributionCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableDistributionCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableDistributionCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDistributionCreate(val *DistributionCreate) *NullableDistributionCreate {
	return &NullableDistributionCreate{value: val, isSet: true}
}

func (v NullableDistributionCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDistributionCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}