/*
 * IONOS Cloud - Monitoring REST API
 *
 * The monitoring service offers a centralized platform to collect and store metrics.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// checks if the KeyCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &KeyCreate{}

// KeyCreate Generates a new key for a pipeline invalidating the old one. The key is used for authentication when sending metrics.
type KeyCreate struct {
	Description *string `json:"description,omitempty"`
}

// NewKeyCreate instantiates a new KeyCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewKeyCreate() *KeyCreate {
	this := KeyCreate{}

	return &this
}

// NewKeyCreateWithDefaults instantiates a new KeyCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewKeyCreateWithDefaults() *KeyCreate {
	this := KeyCreate{}
	return &this
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *KeyCreate) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *KeyCreate) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *KeyCreate) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *KeyCreate) SetDescription(v string) {
	o.Description = &v
}

func (o KeyCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	return toSerialize, nil
}

type NullableKeyCreate struct {
	value *KeyCreate
	isSet bool
}

func (v NullableKeyCreate) Get() *KeyCreate {
	return v.value
}

func (v *NullableKeyCreate) Set(val *KeyCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableKeyCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableKeyCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKeyCreate(val *KeyCreate) *NullableKeyCreate {
	return &NullableKeyCreate{value: val, isSet: true}
}

func (v NullableKeyCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKeyCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
