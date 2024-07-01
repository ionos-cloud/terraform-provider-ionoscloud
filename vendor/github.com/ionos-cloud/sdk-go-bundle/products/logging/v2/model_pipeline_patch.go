/*
 * IONOS Logging REST API
 *
 * Logging as a Service (LaaS) is a service that provides a centralized logging system where users are able to push and aggregate their system or application logs. This service also provides a visualization platform where users are able to observe, search and filter the logs and also create dashboards and alerts for their data points. This service can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an API. The API allows you to create logging pipelines or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the PipelinePatch type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelinePatch{}

// PipelinePatch Request payload with any data that is possible to patch a logging pipeline
type PipelinePatch struct {
	Properties PipelinePatchProperties `json:"properties"`
}

// NewPipelinePatch instantiates a new PipelinePatch object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelinePatch(properties PipelinePatchProperties) *PipelinePatch {
	this := PipelinePatch{}

	this.Properties = properties

	return &this
}

// NewPipelinePatchWithDefaults instantiates a new PipelinePatch object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelinePatchWithDefaults() *PipelinePatch {
	this := PipelinePatch{}
	return &this
}

// GetProperties returns the Properties field value
func (o *PipelinePatch) GetProperties() PipelinePatchProperties {
	if o == nil {
		var ret PipelinePatchProperties
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *PipelinePatch) GetPropertiesOk() (*PipelinePatchProperties, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *PipelinePatch) SetProperties(v PipelinePatchProperties) {
	o.Properties = v
}

func (o PipelinePatch) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PipelinePatch) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullablePipelinePatch struct {
	value *PipelinePatch
	isSet bool
}

func (v NullablePipelinePatch) Get() *PipelinePatch {
	return v.value
}

func (v *NullablePipelinePatch) Set(val *PipelinePatch) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelinePatch) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelinePatch) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelinePatch(val *PipelinePatch) *NullablePipelinePatch {
	return &NullablePipelinePatch{value: val, isSet: true}
}

func (v NullablePipelinePatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelinePatch) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}