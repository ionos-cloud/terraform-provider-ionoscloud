/*
 * IONOS Cloud - Object Storage Management API
 *
 * Object Storage Management API is a RESTful API that manages the object storage service configuration for IONOS Cloud.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// AccessKey Per user access key.
type AccessKey struct {
	// Description of the Access key.
	Description *string `json:"description"`
	// Access key metadata is a string of 92 characters.
	AccessKey *string `json:"accessKey"`
	// The secret key of the Access key.
	SecretKey *string `json:"secretKey"`
	// The canonical user ID which is valid for user-owned buckets.
	CanonicalUserId *string `json:"canonicalUserId,omitempty"`
	// The contract user ID which is valid for contract-owned buckets.
	ContractUserId *string `json:"contractUserId,omitempty"`
}

// NewAccessKey instantiates a new AccessKey object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccessKey(description string, accessKey string, secretKey string) *AccessKey {
	this := AccessKey{}

	this.Description = &description
	this.AccessKey = &accessKey
	this.SecretKey = &secretKey

	return &this
}

// NewAccessKeyWithDefaults instantiates a new AccessKey object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccessKeyWithDefaults() *AccessKey {
	this := AccessKey{}
	return &this
}

// GetDescription returns the Description field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKey) GetDescription() *string {
	if o == nil {
		return nil
	}

	return o.Description

}

// GetDescriptionOk returns a tuple with the Description field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKey) GetDescriptionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Description, true
}

// SetDescription sets field value
func (o *AccessKey) SetDescription(v string) {

	o.Description = &v

}

// HasDescription returns a boolean if a field has been set.
func (o *AccessKey) HasDescription() bool {
	if o != nil && o.Description != nil {
		return true
	}

	return false
}

// GetAccessKey returns the AccessKey field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKey) GetAccessKey() *string {
	if o == nil {
		return nil
	}

	return o.AccessKey

}

// GetAccessKeyOk returns a tuple with the AccessKey field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKey) GetAccessKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.AccessKey, true
}

// SetAccessKey sets field value
func (o *AccessKey) SetAccessKey(v string) {

	o.AccessKey = &v

}

// HasAccessKey returns a boolean if a field has been set.
func (o *AccessKey) HasAccessKey() bool {
	if o != nil && o.AccessKey != nil {
		return true
	}

	return false
}

// GetSecretKey returns the SecretKey field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKey) GetSecretKey() *string {
	if o == nil {
		return nil
	}

	return o.SecretKey

}

// GetSecretKeyOk returns a tuple with the SecretKey field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKey) GetSecretKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.SecretKey, true
}

// SetSecretKey sets field value
func (o *AccessKey) SetSecretKey(v string) {

	o.SecretKey = &v

}

// HasSecretKey returns a boolean if a field has been set.
func (o *AccessKey) HasSecretKey() bool {
	if o != nil && o.SecretKey != nil {
		return true
	}

	return false
}

// GetCanonicalUserId returns the CanonicalUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKey) GetCanonicalUserId() *string {
	if o == nil {
		return nil
	}

	return o.CanonicalUserId

}

// GetCanonicalUserIdOk returns a tuple with the CanonicalUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKey) GetCanonicalUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CanonicalUserId, true
}

// SetCanonicalUserId sets field value
func (o *AccessKey) SetCanonicalUserId(v string) {

	o.CanonicalUserId = &v

}

// HasCanonicalUserId returns a boolean if a field has been set.
func (o *AccessKey) HasCanonicalUserId() bool {
	if o != nil && o.CanonicalUserId != nil {
		return true
	}

	return false
}

// GetContractUserId returns the ContractUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKey) GetContractUserId() *string {
	if o == nil {
		return nil
	}

	return o.ContractUserId

}

// GetContractUserIdOk returns a tuple with the ContractUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKey) GetContractUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ContractUserId, true
}

// SetContractUserId sets field value
func (o *AccessKey) SetContractUserId(v string) {

	o.ContractUserId = &v

}

// HasContractUserId returns a boolean if a field has been set.
func (o *AccessKey) HasContractUserId() bool {
	if o != nil && o.ContractUserId != nil {
		return true
	}

	return false
}

func (o AccessKey) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Description != nil {
		toSerialize["description"] = o.Description
	}

	if o.AccessKey != nil {
		toSerialize["accessKey"] = o.AccessKey
	}

	if o.SecretKey != nil {
		toSerialize["secretKey"] = o.SecretKey
	}

	if o.CanonicalUserId != nil {
		toSerialize["canonicalUserId"] = o.CanonicalUserId
	}

	if o.ContractUserId != nil {
		toSerialize["contractUserId"] = o.ContractUserId
	}

	return json.Marshal(toSerialize)
}

type NullableAccessKey struct {
	value *AccessKey
	isSet bool
}

func (v NullableAccessKey) Get() *AccessKey {
	return v.value
}

func (v *NullableAccessKey) Set(val *AccessKey) {
	v.value = val
	v.isSet = true
}

func (v NullableAccessKey) IsSet() bool {
	return v.isSet
}

func (v *NullableAccessKey) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccessKey(val *AccessKey) *NullableAccessKey {
	return &NullableAccessKey{value: val, isSet: true}
}

func (v NullableAccessKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccessKey) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}