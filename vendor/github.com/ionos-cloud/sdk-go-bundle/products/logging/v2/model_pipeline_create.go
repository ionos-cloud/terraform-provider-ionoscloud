/*
 * IONOS Logging Service REST API
 *
 * The Logging Service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data. This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the PipelineCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelineCreate{}

// PipelineCreate A pipeline consists of the building blocks of a centralized logging system including supported log agents and log sources and also public endpoints to push and access logs.
type PipelineCreate struct {
	Properties PipelineNoAddr `json:"properties"`
}

// NewPipelineCreate instantiates a new PipelineCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineCreate(properties PipelineNoAddr) *PipelineCreate {
	this := PipelineCreate{}

	this.Properties = properties

	return &this
}

// NewPipelineCreateWithDefaults instantiates a new PipelineCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineCreateWithDefaults() *PipelineCreate {
	this := PipelineCreate{}
	return &this
}

// GetProperties returns the Properties field value
func (o *PipelineCreate) GetProperties() PipelineNoAddr {
	if o == nil {
		var ret PipelineNoAddr
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *PipelineCreate) GetPropertiesOk() (*PipelineNoAddr, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *PipelineCreate) SetProperties(v PipelineNoAddr) {
	o.Properties = v
}

func (o PipelineCreate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PipelineCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullablePipelineCreate struct {
	value *PipelineCreate
	isSet bool
}

func (v NullablePipelineCreate) Get() *PipelineCreate {
	return v.value
}

func (v *NullablePipelineCreate) Set(val *PipelineCreate) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineCreate) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineCreate(val *PipelineCreate) *NullablePipelineCreate {
	return &NullablePipelineCreate{value: val, isSet: true}
}

func (v NullablePipelineCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
