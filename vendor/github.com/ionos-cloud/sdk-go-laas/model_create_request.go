/*
 * IONOS Logging REST API
 *
 * Logging as a Service (LaaS) is a service that provides a centralized logging system where users are able to push and aggregate their system or application logs. This service also provides a visualization platform where users are able to observe, search and filter the logs and also create dashboards and alerts for their data points. This service can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an API. The API allows you to create logging pipelines or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CreateRequest Request payload with all data needed to create a new logging pipeline
type CreateRequest struct {
	Properties *CreateRequestProperties `json:"properties"`
}

// NewCreateRequest instantiates a new CreateRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateRequest(properties CreateRequestProperties) *CreateRequest {
	this := CreateRequest{}

	this.Properties = &properties

	return &this
}

// NewCreateRequestWithDefaults instantiates a new CreateRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateRequestWithDefaults() *CreateRequest {
	this := CreateRequest{}
	return &this
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for CreateRequestProperties will be returned
func (o *CreateRequest) GetProperties() *CreateRequestProperties {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateRequest) GetPropertiesOk() (*CreateRequestProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *CreateRequest) SetProperties(v CreateRequestProperties) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *CreateRequest) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o CreateRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableCreateRequest struct {
	value *CreateRequest
	isSet bool
}

func (v NullableCreateRequest) Get() *CreateRequest {
	return v.value
}

func (v *NullableCreateRequest) Set(val *CreateRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateRequest(val *CreateRequest) *NullableCreateRequest {
	return &NullableCreateRequest{value: val, isSet: true}
}

func (v NullableCreateRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
