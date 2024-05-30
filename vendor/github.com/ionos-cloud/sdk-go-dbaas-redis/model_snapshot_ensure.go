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

// SnapshotEnsure struct for SnapshotEnsure
type SnapshotEnsure struct {
	// The ID (UUID) of the Snapshot.
	Id *string `json:"id"`
	// Metadata
	Metadata *map[string]interface{} `json:"metadata,omitempty"`
	// A point in time snapshot of a Redis replica set.
	Properties *map[string]interface{} `json:"properties"`
}

// NewSnapshotEnsure instantiates a new SnapshotEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSnapshotEnsure(id string, properties map[string]interface{}) *SnapshotEnsure {
	this := SnapshotEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewSnapshotEnsureWithDefaults instantiates a new SnapshotEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSnapshotEnsureWithDefaults() *SnapshotEnsure {
	this := SnapshotEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *SnapshotEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SnapshotEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *SnapshotEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *SnapshotEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *SnapshotEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SnapshotEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *SnapshotEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *SnapshotEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *SnapshotEnsure) GetProperties() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *SnapshotEnsure) GetPropertiesOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *SnapshotEnsure) SetProperties(v map[string]interface{}) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *SnapshotEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o SnapshotEnsure) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableSnapshotEnsure struct {
	value *SnapshotEnsure
	isSet bool
}

func (v NullableSnapshotEnsure) Get() *SnapshotEnsure {
	return v.value
}

func (v *NullableSnapshotEnsure) Set(val *SnapshotEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableSnapshotEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableSnapshotEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSnapshotEnsure(val *SnapshotEnsure) *NullableSnapshotEnsure {
	return &NullableSnapshotEnsure{value: val, isSet: true}
}

func (v NullableSnapshotEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSnapshotEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
