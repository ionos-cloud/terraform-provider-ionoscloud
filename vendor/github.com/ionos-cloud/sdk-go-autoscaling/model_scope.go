/*
 * Container Registry service (CloudAPI)
 *
 * Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their manage Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls.
 *
 * API version: 1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Scope struct for Scope
type Scope struct {
	Actions *[]string `json:"actions"`
	Name    *string   `json:"name"`
	Type    *string   `json:"type"`
}

// NewScope instantiates a new Scope object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewScope(actions []string, name string, type_ string) *Scope {
	this := Scope{}

	this.Actions = &actions
	this.Name = &name
	this.Type = &type_

	return &this
}

// NewScopeWithDefaults instantiates a new Scope object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewScopeWithDefaults() *Scope {
	this := Scope{}
	return &this
}

// GetActions returns the Actions field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *Scope) GetActions() *[]string {
	if o == nil {
		return nil
	}

	return o.Actions

}

// GetActionsOk returns a tuple with the Actions field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Scope) GetActionsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Actions, true
}

// SetActions sets field value
func (o *Scope) SetActions(v []string) {

	o.Actions = &v

}

// HasActions returns a boolean if a field has been set.
func (o *Scope) HasActions() bool {
	if o != nil && o.Actions != nil {
		return true
	}

	return false
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Scope) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Scope) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *Scope) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *Scope) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Scope) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Scope) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *Scope) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Scope) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

func (o Scope) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["actions"] = o.Actions

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	return json.Marshal(toSerialize)
}

type NullableScope struct {
	value *Scope
	isSet bool
}

func (v NullableScope) Get() *Scope {
	return v.value
}

func (v *NullableScope) Set(val *Scope) {
	v.value = val
	v.isSet = true
}

func (v NullableScope) IsSet() bool {
	return v.isSet
}

func (v *NullableScope) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableScope(val *Scope) *NullableScope {
	return &NullableScope{value: val, isSet: true}
}

func (v NullableScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableScope) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
