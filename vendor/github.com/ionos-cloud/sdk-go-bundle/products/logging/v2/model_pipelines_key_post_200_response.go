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

// checks if the PipelinesKeyPost200Response type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelinesKeyPost200Response{}

// PipelinesKeyPost200Response struct for PipelinesKeyPost200Response
type PipelinesKeyPost200Response struct {
	Key *string `json:"key,omitempty"`
}

// NewPipelinesKeyPost200Response instantiates a new PipelinesKeyPost200Response object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelinesKeyPost200Response() *PipelinesKeyPost200Response {
	this := PipelinesKeyPost200Response{}

	return &this
}

// NewPipelinesKeyPost200ResponseWithDefaults instantiates a new PipelinesKeyPost200Response object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelinesKeyPost200ResponseWithDefaults() *PipelinesKeyPost200Response {
	this := PipelinesKeyPost200Response{}
	return &this
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *PipelinesKeyPost200Response) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelinesKeyPost200Response) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *PipelinesKeyPost200Response) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *PipelinesKeyPost200Response) SetKey(v string) {
	o.Key = &v
}

func (o PipelinesKeyPost200Response) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Key) {
		toSerialize["key"] = o.Key
	}
	return toSerialize, nil
}

type NullablePipelinesKeyPost200Response struct {
	value *PipelinesKeyPost200Response
	isSet bool
}

func (v NullablePipelinesKeyPost200Response) Get() *PipelinesKeyPost200Response {
	return v.value
}

func (v *NullablePipelinesKeyPost200Response) Set(val *PipelinesKeyPost200Response) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelinesKeyPost200Response) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelinesKeyPost200Response) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelinesKeyPost200Response(val *PipelinesKeyPost200Response) *NullablePipelinesKeyPost200Response {
	return &NullablePipelinesKeyPost200Response{value: val, isSet: true}
}

func (v NullablePipelinesKeyPost200Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelinesKeyPost200Response) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
