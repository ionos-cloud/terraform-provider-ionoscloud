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
	"time"
)

// RestoreMetadataAllOf struct for RestoreMetadataAllOf
type RestoreMetadataAllOf struct {
	// The time the snapshot was dumped from the replica set.
	RestoreTime *IonosTime `json:"restoreTime,omitempty"`
	// The ID of the snapshot that was restored.
	RestoredSnapshotId *string `json:"restoredSnapshotId,omitempty"`
}

// NewRestoreMetadataAllOf instantiates a new RestoreMetadataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRestoreMetadataAllOf() *RestoreMetadataAllOf {
	this := RestoreMetadataAllOf{}

	return &this
}

// NewRestoreMetadataAllOfWithDefaults instantiates a new RestoreMetadataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRestoreMetadataAllOfWithDefaults() *RestoreMetadataAllOf {
	this := RestoreMetadataAllOf{}
	return &this
}

// GetRestoreTime returns the RestoreTime field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *RestoreMetadataAllOf) GetRestoreTime() *time.Time {
	if o == nil {
		return nil
	}

	if o.RestoreTime == nil {
		return nil
	}
	return &o.RestoreTime.Time

}

// GetRestoreTimeOk returns a tuple with the RestoreTime field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreMetadataAllOf) GetRestoreTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.RestoreTime == nil {
		return nil, false
	}
	return &o.RestoreTime.Time, true

}

// SetRestoreTime sets field value
func (o *RestoreMetadataAllOf) SetRestoreTime(v time.Time) {

	o.RestoreTime = &IonosTime{v}

}

// HasRestoreTime returns a boolean if a field has been set.
func (o *RestoreMetadataAllOf) HasRestoreTime() bool {
	if o != nil && o.RestoreTime != nil {
		return true
	}

	return false
}

// GetRestoredSnapshotId returns the RestoredSnapshotId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RestoreMetadataAllOf) GetRestoredSnapshotId() *string {
	if o == nil {
		return nil
	}

	return o.RestoredSnapshotId

}

// GetRestoredSnapshotIdOk returns a tuple with the RestoredSnapshotId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreMetadataAllOf) GetRestoredSnapshotIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.RestoredSnapshotId, true
}

// SetRestoredSnapshotId sets field value
func (o *RestoreMetadataAllOf) SetRestoredSnapshotId(v string) {

	o.RestoredSnapshotId = &v

}

// HasRestoredSnapshotId returns a boolean if a field has been set.
func (o *RestoreMetadataAllOf) HasRestoredSnapshotId() bool {
	if o != nil && o.RestoredSnapshotId != nil {
		return true
	}

	return false
}

func (o RestoreMetadataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.RestoreTime != nil {
		toSerialize["restoreTime"] = o.RestoreTime
	}

	if o.RestoredSnapshotId != nil {
		toSerialize["restoredSnapshotId"] = o.RestoredSnapshotId
	}

	return json.Marshal(toSerialize)
}

type NullableRestoreMetadataAllOf struct {
	value *RestoreMetadataAllOf
	isSet bool
}

func (v NullableRestoreMetadataAllOf) Get() *RestoreMetadataAllOf {
	return v.value
}

func (v *NullableRestoreMetadataAllOf) Set(val *RestoreMetadataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableRestoreMetadataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableRestoreMetadataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRestoreMetadataAllOf(val *RestoreMetadataAllOf) *NullableRestoreMetadataAllOf {
	return &NullableRestoreMetadataAllOf{value: val, isSet: true}
}

func (v NullableRestoreMetadataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRestoreMetadataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
