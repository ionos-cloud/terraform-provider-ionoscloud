/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// UserPost struct for UserPost
type UserPost struct {
	Properties *UserPropertiesPost `json:"properties"`
}



// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for UserPropertiesPost will be returned
func (o *UserPost) GetProperties() *UserPropertiesPost {
	if o == nil {
		return nil
	}


	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UserPost) GetPropertiesOk() (*UserPropertiesPost, bool) {
	if o == nil {
		return nil, false
	}


	return o.Properties, true
}

// SetProperties sets field value
func (o *UserPost) SetProperties(v UserPropertiesPost) {


	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *UserPost) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}


func (o UserPost) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}
	
	return json.Marshal(toSerialize)
}

type NullableUserPost struct {
	value *UserPost
	isSet bool
}

func (v NullableUserPost) Get() *UserPost {
	return v.value
}

func (v *NullableUserPost) Set(val *UserPost) {
	v.value = val
	v.isSet = true
}

func (v NullableUserPost) IsSet() bool {
	return v.isSet
}

func (v *NullableUserPost) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserPost(val *UserPost) *NullableUserPost {
	return &NullableUserPost{value: val, isSet: true}
}

func (v NullableUserPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserPost) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


