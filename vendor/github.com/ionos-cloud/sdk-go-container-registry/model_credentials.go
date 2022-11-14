/*
 * Container Registry service
 *
 * Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls.
 *
 * API version: 1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Credentials struct for Credentials
type Credentials struct {
	Password *string `json:"password"`
	Username *string `json:"username"`
}

// NewCredentials instantiates a new Credentials object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCredentials(password string, username string) *Credentials {
	this := Credentials{}

	this.Password = &password
	this.Username = &username

	return &this
}

// NewCredentialsWithDefaults instantiates a new Credentials object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCredentialsWithDefaults() *Credentials {
	this := Credentials{}
	return &this
}

// GetPassword returns the Password field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Credentials) GetPassword() *string {
	if o == nil {
		return nil
	}

	return o.Password

}

// GetPasswordOk returns a tuple with the Password field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Credentials) GetPasswordOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Password, true
}

// SetPassword sets field value
func (o *Credentials) SetPassword(v string) {

	o.Password = &v

}

// HasPassword returns a boolean if a field has been set.
func (o *Credentials) HasPassword() bool {
	if o != nil && o.Password != nil {
		return true
	}

	return false
}

// GetUsername returns the Username field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Credentials) GetUsername() *string {
	if o == nil {
		return nil
	}

	return o.Username

}

// GetUsernameOk returns a tuple with the Username field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Credentials) GetUsernameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Username, true
}

// SetUsername sets field value
func (o *Credentials) SetUsername(v string) {

	o.Username = &v

}

// HasUsername returns a boolean if a field has been set.
func (o *Credentials) HasUsername() bool {
	if o != nil && o.Username != nil {
		return true
	}

	return false
}

func (o Credentials) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Password != nil {
		toSerialize["password"] = o.Password
	}

	if o.Username != nil {
		toSerialize["username"] = o.Username
	}

	return json.Marshal(toSerialize)
}

type NullableCredentials struct {
	value *Credentials
	isSet bool
}

func (v NullableCredentials) Get() *Credentials {
	return v.value
}

func (v *NullableCredentials) Set(val *Credentials) {
	v.value = val
	v.isSet = true
}

func (v NullableCredentials) IsSet() bool {
	return v.isSet
}

func (v *NullableCredentials) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCredentials(val *Credentials) *NullableCredentials {
	return &NullableCredentials{value: val, isSet: true}
}

func (v NullableCredentials) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCredentials) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
