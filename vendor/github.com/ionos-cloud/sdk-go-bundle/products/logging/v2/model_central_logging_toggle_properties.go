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

// checks if the CentralLoggingToggleProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CentralLoggingToggleProperties{}

// CentralLoggingToggleProperties struct for CentralLoggingToggleProperties
type CentralLoggingToggleProperties struct {
	Enabled bool `json:"enabled"`
}

// NewCentralLoggingToggleProperties instantiates a new CentralLoggingToggleProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCentralLoggingToggleProperties(enabled bool) *CentralLoggingToggleProperties {
	this := CentralLoggingToggleProperties{}

	this.Enabled = enabled

	return &this
}

// NewCentralLoggingTogglePropertiesWithDefaults instantiates a new CentralLoggingToggleProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCentralLoggingTogglePropertiesWithDefaults() *CentralLoggingToggleProperties {
	this := CentralLoggingToggleProperties{}
	var enabled bool = false
	this.Enabled = enabled
	return &this
}

// GetEnabled returns the Enabled field value
func (o *CentralLoggingToggleProperties) GetEnabled() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value
// and a boolean to check if the value has been set.
func (o *CentralLoggingToggleProperties) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Enabled, true
}

// SetEnabled sets field value
func (o *CentralLoggingToggleProperties) SetEnabled(v bool) {
	o.Enabled = v
}

func (o CentralLoggingToggleProperties) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CentralLoggingToggleProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	return toSerialize, nil
}

type NullableCentralLoggingToggleProperties struct {
	value *CentralLoggingToggleProperties
	isSet bool
}

func (v NullableCentralLoggingToggleProperties) Get() *CentralLoggingToggleProperties {
	return v.value
}

func (v *NullableCentralLoggingToggleProperties) Set(val *CentralLoggingToggleProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableCentralLoggingToggleProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableCentralLoggingToggleProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCentralLoggingToggleProperties(val *CentralLoggingToggleProperties) *NullableCentralLoggingToggleProperties {
	return &NullableCentralLoggingToggleProperties{value: val, isSet: true}
}

func (v NullableCentralLoggingToggleProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCentralLoggingToggleProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
