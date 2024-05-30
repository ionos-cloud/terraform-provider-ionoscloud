/*
 * Redis DB API
 *
 * Redis Database API
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// SnapshotCreate struct for SnapshotCreate
type SnapshotCreate struct {
	// Metadata
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	// A point in time snapshot of a Redis replica set.
	Properties *map[string]interface{} `json:"properties"`
}

// NewSnapshotCreate instantiates a new SnapshotCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSnapshotCreate(properties map[string]interface{}) *SnapshotCreate {
	this := SnapshotCreate{}

	this.Properties = &properties

	return &this
}

// NewSnapshotCreateWithDefaults instantiates a new SnapshotCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSnapshotCreateWithDefaults() *SnapshotCreate {
	this := SnapshotCreate{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *SnapshotCreate) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SnapshotCreate) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *SnapshotCreate) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *SnapshotCreate) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *SnapshotCreate) GetProperties() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SnapshotCreate) GetPropertiesOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *SnapshotCreate) SetProperties(v map[string]interface{}) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *SnapshotCreate) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o SnapshotCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableSnapshotCreate struct {
	value *SnapshotCreate
	isSet bool
}

func (v NullableSnapshotCreate) Get() *SnapshotCreate {
	return v.value
}

func (v *NullableSnapshotCreate) Set(val *SnapshotCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableSnapshotCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableSnapshotCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSnapshotCreate(val *SnapshotCreate) *NullableSnapshotCreate {
	return &NullableSnapshotCreate{value: val, isSet: true}
}

func (v NullableSnapshotCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSnapshotCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
