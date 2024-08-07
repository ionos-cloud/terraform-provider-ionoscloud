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

// RequestProperties struct for RequestProperties
type RequestProperties struct {
	Method  *string            `json:"method,omitempty"`
	Headers *map[string]string `json:"headers,omitempty"`
	Body    *string            `json:"body,omitempty"`
	Url     *string            `json:"url,omitempty"`
}

// NewRequestProperties instantiates a new RequestProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRequestProperties() *RequestProperties {
	this := RequestProperties{}

	return &this
}

// NewRequestPropertiesWithDefaults instantiates a new RequestProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRequestPropertiesWithDefaults() *RequestProperties {
	this := RequestProperties{}
	return &this
}

// GetMethod returns the Method field value
// If the value is explicit nil, nil is returned
func (o *RequestProperties) GetMethod() *string {
	if o == nil {
		return nil
	}

	return o.Method

}

// GetMethodOk returns a tuple with the Method field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestProperties) GetMethodOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Method, true
}

// SetMethod sets field value
func (o *RequestProperties) SetMethod(v string) {

	o.Method = &v

}

// HasMethod returns a boolean if a field has been set.
func (o *RequestProperties) HasMethod() bool {
	if o != nil && o.Method != nil {
		return true
	}

	return false
}

// GetHeaders returns the Headers field value
// If the value is explicit nil, nil is returned
func (o *RequestProperties) GetHeaders() *map[string]string {
	if o == nil {
		return nil
	}

	return o.Headers

}

// GetHeadersOk returns a tuple with the Headers field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestProperties) GetHeadersOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Headers, true
}

// SetHeaders sets field value
func (o *RequestProperties) SetHeaders(v map[string]string) {

	o.Headers = &v

}

// HasHeaders returns a boolean if a field has been set.
func (o *RequestProperties) HasHeaders() bool {
	if o != nil && o.Headers != nil {
		return true
	}

	return false
}

// GetBody returns the Body field value
// If the value is explicit nil, nil is returned
func (o *RequestProperties) GetBody() *string {
	if o == nil {
		return nil
	}

	return o.Body

}

// GetBodyOk returns a tuple with the Body field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestProperties) GetBodyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Body, true
}

// SetBody sets field value
func (o *RequestProperties) SetBody(v string) {

	o.Body = &v

}

// HasBody returns a boolean if a field has been set.
func (o *RequestProperties) HasBody() bool {
	if o != nil && o.Body != nil {
		return true
	}

	return false
}

// GetUrl returns the Url field value
// If the value is explicit nil, nil is returned
func (o *RequestProperties) GetUrl() *string {
	if o == nil {
		return nil
	}

	return o.Url

}

// GetUrlOk returns a tuple with the Url field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RequestProperties) GetUrlOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Url, true
}

// SetUrl sets field value
func (o *RequestProperties) SetUrl(v string) {

	o.Url = &v

}

// HasUrl returns a boolean if a field has been set.
func (o *RequestProperties) HasUrl() bool {
	if o != nil && o.Url != nil {
		return true
	}

	return false
}

func (o RequestProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Method != nil {
		toSerialize["method"] = o.Method
	}

	if o.Headers != nil {
		toSerialize["headers"] = o.Headers
	}

	if o.Body != nil {
		toSerialize["body"] = o.Body
	}

	if o.Url != nil {
		toSerialize["url"] = o.Url
	}

	return json.Marshal(toSerialize)
}

type NullableRequestProperties struct {
	value *RequestProperties
	isSet bool
}

func (v NullableRequestProperties) Get() *RequestProperties {
	return v.value
}

func (v *NullableRequestProperties) Set(val *RequestProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableRequestProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableRequestProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRequestProperties(val *RequestProperties) *NullableRequestProperties {
	return &NullableRequestProperties{value: val, isSet: true}
}

func (v NullableRequestProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRequestProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
