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

// RestoreSnapshotProperties struct for RestoreSnapshotProperties
type RestoreSnapshotProperties struct {
	// The id of the snapshot
	SnapshotId *string `json:"snapshotId,omitempty"`
}

// NewRestoreSnapshotProperties instantiates a new RestoreSnapshotProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRestoreSnapshotProperties() *RestoreSnapshotProperties {
	this := RestoreSnapshotProperties{}

	return &this
}

// NewRestoreSnapshotPropertiesWithDefaults instantiates a new RestoreSnapshotProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRestoreSnapshotPropertiesWithDefaults() *RestoreSnapshotProperties {
	this := RestoreSnapshotProperties{}
	return &this
}

// GetSnapshotId returns the SnapshotId field value
// If the value is explicit nil, nil is returned
func (o *RestoreSnapshotProperties) GetSnapshotId() *string {
	if o == nil {
		return nil
	}

	return o.SnapshotId

}

// GetSnapshotIdOk returns a tuple with the SnapshotId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreSnapshotProperties) GetSnapshotIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.SnapshotId, true
}

// SetSnapshotId sets field value
func (o *RestoreSnapshotProperties) SetSnapshotId(v string) {

	o.SnapshotId = &v

}

// HasSnapshotId returns a boolean if a field has been set.
func (o *RestoreSnapshotProperties) HasSnapshotId() bool {
	if o != nil && o.SnapshotId != nil {
		return true
	}

	return false
}

func (o RestoreSnapshotProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.SnapshotId != nil {
		toSerialize["snapshotId"] = o.SnapshotId
	}

	return json.Marshal(toSerialize)
}

type NullableRestoreSnapshotProperties struct {
	value *RestoreSnapshotProperties
	isSet bool
}

func (v NullableRestoreSnapshotProperties) Get() *RestoreSnapshotProperties {
	return v.value
}

func (v *NullableRestoreSnapshotProperties) Set(val *RestoreSnapshotProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableRestoreSnapshotProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableRestoreSnapshotProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRestoreSnapshotProperties(val *RestoreSnapshotProperties) *NullableRestoreSnapshotProperties {
	return &NullableRestoreSnapshotProperties{value: val, isSet: true}
}

func (v NullableRestoreSnapshotProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRestoreSnapshotProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
