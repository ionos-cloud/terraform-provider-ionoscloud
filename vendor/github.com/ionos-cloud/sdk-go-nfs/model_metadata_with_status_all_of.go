/*
 * IONOS Cloud - Network File Storage API
 *
 * The RESTful API for managing Network File Storage.
 *
 * API version: 0.1.3
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// MetadataWithStatusAllOf struct for MetadataWithStatusAllOf
type MetadataWithStatusAllOf struct {
	// The status of the resource can be one of the following:  * `AVAILABLE` - The resource exists and is healthy. * `PROVISIONING` - The resource is being created or updated. * `DESTROYING` - A delete command was issued, and the resource is being deleted. * `FAILED` - The resource failed, with details provided in `statusMessage`.
	Status *string `json:"status"`
	// The message of the failure if the status is `FAILED`.
	StatusMessage *string `json:"statusMessage,omitempty"`
}

// NewMetadataWithStatusAllOf instantiates a new MetadataWithStatusAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataWithStatusAllOf(status string) *MetadataWithStatusAllOf {
	this := MetadataWithStatusAllOf{}

	this.Status = &status

	return &this
}

// NewMetadataWithStatusAllOfWithDefaults instantiates a new MetadataWithStatusAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithStatusAllOfWithDefaults() *MetadataWithStatusAllOf {
	this := MetadataWithStatusAllOf{}
	return &this
}

// GetStatus returns the Status field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithStatusAllOf) GetStatus() *string {
	if o == nil {
		return nil
	}

	return o.Status

}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithStatusAllOf) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Status, true
}

// SetStatus sets field value
func (o *MetadataWithStatusAllOf) SetStatus(v string) {

	o.Status = &v

}

// HasStatus returns a boolean if a field has been set.
func (o *MetadataWithStatusAllOf) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// GetStatusMessage returns the StatusMessage field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MetadataWithStatusAllOf) GetStatusMessage() *string {
	if o == nil {
		return nil
	}

	return o.StatusMessage

}

// GetStatusMessageOk returns a tuple with the StatusMessage field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MetadataWithStatusAllOf) GetStatusMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.StatusMessage, true
}

// SetStatusMessage sets field value
func (o *MetadataWithStatusAllOf) SetStatusMessage(v string) {

	o.StatusMessage = &v

}

// HasStatusMessage returns a boolean if a field has been set.
func (o *MetadataWithStatusAllOf) HasStatusMessage() bool {
	if o != nil && o.StatusMessage != nil {
		return true
	}

	return false
}

func (o MetadataWithStatusAllOf) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Status != nil {
		toSerialize["status"] = o.Status
	}

	if o.StatusMessage != nil {
		toSerialize["statusMessage"] = o.StatusMessage
	}

	return json.Marshal(toSerialize)
}

type NullableMetadataWithStatusAllOf struct {
	value *MetadataWithStatusAllOf
	isSet bool
}

func (v NullableMetadataWithStatusAllOf) Get() *MetadataWithStatusAllOf {
	return v.value
}

func (v *NullableMetadataWithStatusAllOf) Set(val *MetadataWithStatusAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataWithStatusAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataWithStatusAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataWithStatusAllOf(val *MetadataWithStatusAllOf) *NullableMetadataWithStatusAllOf {
	return &NullableMetadataWithStatusAllOf{value: val, isSet: true}
}

func (v NullableMetadataWithStatusAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataWithStatusAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
