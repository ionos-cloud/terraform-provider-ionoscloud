/*
 * CLOUD API
 *
 *  IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CreateSnapshot struct for CreateSnapshot
type CreateSnapshot struct {
	Properties *CreateSnapshotProperties `json:"properties,omitempty"`
}

// NewCreateSnapshot instantiates a new CreateSnapshot object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateSnapshot() *CreateSnapshot {
	this := CreateSnapshot{}

	return &this
}

// NewCreateSnapshotWithDefaults instantiates a new CreateSnapshot object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateSnapshotWithDefaults() *CreateSnapshot {
	this := CreateSnapshot{}
	return &this
}

// GetProperties returns the Properties field value
// If the value is explicit nil, nil is returned
func (o *CreateSnapshot) GetProperties() *CreateSnapshotProperties {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateSnapshot) GetPropertiesOk() (*CreateSnapshotProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *CreateSnapshot) SetProperties(v CreateSnapshotProperties) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *CreateSnapshot) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o CreateSnapshot) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableCreateSnapshot struct {
	value *CreateSnapshot
	isSet bool
}

func (v NullableCreateSnapshot) Get() *CreateSnapshot {
	return v.value
}

func (v *NullableCreateSnapshot) Set(val *CreateSnapshot) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateSnapshot) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateSnapshot) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateSnapshot(val *CreateSnapshot) *NullableCreateSnapshot {
	return &NullableCreateSnapshot{value: val, isSet: true}
}

func (v NullableCreateSnapshot) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateSnapshot) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}