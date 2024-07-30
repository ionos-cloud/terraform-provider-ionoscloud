/*
 * IONOS Logging REST API
 *
 * The logging service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data.  This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CentralLoggingResponseProperties struct for CentralLoggingResponseProperties
type CentralLoggingResponseProperties struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// NewCentralLoggingResponseProperties instantiates a new CentralLoggingResponseProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCentralLoggingResponseProperties() *CentralLoggingResponseProperties {
	this := CentralLoggingResponseProperties{}

	var enabled bool = false
	this.Enabled = &enabled

	return &this
}

// NewCentralLoggingResponsePropertiesWithDefaults instantiates a new CentralLoggingResponseProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCentralLoggingResponsePropertiesWithDefaults() *CentralLoggingResponseProperties {
	this := CentralLoggingResponseProperties{}
	var enabled bool = false
	this.Enabled = &enabled
	return &this
}

// GetEnabled returns the Enabled field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *CentralLoggingResponseProperties) GetEnabled() *bool {
	if o == nil {
		return nil
	}

	return o.Enabled

}

// GetEnabledOk returns a tuple with the Enabled field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CentralLoggingResponseProperties) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}

	return o.Enabled, true
}

// SetEnabled sets field value
func (o *CentralLoggingResponseProperties) SetEnabled(v bool) {

	o.Enabled = &v

}

// HasEnabled returns a boolean if a field has been set.
func (o *CentralLoggingResponseProperties) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

func (o CentralLoggingResponseProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}

	return json.Marshal(toSerialize)
}

type NullableCentralLoggingResponseProperties struct {
	value *CentralLoggingResponseProperties
	isSet bool
}

func (v NullableCentralLoggingResponseProperties) Get() *CentralLoggingResponseProperties {
	return v.value
}

func (v *NullableCentralLoggingResponseProperties) Set(val *CentralLoggingResponseProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableCentralLoggingResponseProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableCentralLoggingResponseProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCentralLoggingResponseProperties(val *CentralLoggingResponseProperties) *NullableCentralLoggingResponseProperties {
	return &NullableCentralLoggingResponseProperties{value: val, isSet: true}
}

func (v NullableCentralLoggingResponseProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCentralLoggingResponseProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
