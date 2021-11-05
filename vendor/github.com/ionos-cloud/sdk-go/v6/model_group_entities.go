/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// GroupEntities struct for GroupEntities
type GroupEntities struct {
	Users     *GroupMembers   `json:"users,omitempty"`
	Resources *ResourceGroups `json:"resources,omitempty"`
}

// GetUsers returns the Users field value
// If the value is explicit nil, the zero value for GroupMembers will be returned
func (o *GroupEntities) GetUsers() *GroupMembers {
	if o == nil {
		return nil
	}

	return o.Users

}

// GetUsersOk returns a tuple with the Users field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GroupEntities) GetUsersOk() (*GroupMembers, bool) {
	if o == nil {
		return nil, false
	}

	return o.Users, true
}

// SetUsers sets field value
func (o *GroupEntities) SetUsers(v GroupMembers) {

	o.Users = &v

}

// HasUsers returns a boolean if a field has been set.
func (o *GroupEntities) HasUsers() bool {
	if o != nil && o.Users != nil {
		return true
	}

	return false
}

// GetResources returns the Resources field value
// If the value is explicit nil, the zero value for ResourceGroups will be returned
func (o *GroupEntities) GetResources() *ResourceGroups {
	if o == nil {
		return nil
	}

	return o.Resources

}

// GetResourcesOk returns a tuple with the Resources field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GroupEntities) GetResourcesOk() (*ResourceGroups, bool) {
	if o == nil {
		return nil, false
	}

	return o.Resources, true
}

// SetResources sets field value
func (o *GroupEntities) SetResources(v ResourceGroups) {

	o.Resources = &v

}

// HasResources returns a boolean if a field has been set.
func (o *GroupEntities) HasResources() bool {
	if o != nil && o.Resources != nil {
		return true
	}

	return false
}

func (o GroupEntities) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Users != nil {
		toSerialize["users"] = o.Users
	}

	if o.Resources != nil {
		toSerialize["resources"] = o.Resources
	}
	return json.Marshal(toSerialize)
}

type NullableGroupEntities struct {
	value *GroupEntities
	isSet bool
}

func (v NullableGroupEntities) Get() *GroupEntities {
	return v.value
}

func (v *NullableGroupEntities) Set(val *GroupEntities) {
	v.value = val
	v.isSet = true
}

func (v NullableGroupEntities) IsSet() bool {
	return v.isSet
}

func (v *NullableGroupEntities) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroupEntities(val *GroupEntities) *NullableGroupEntities {
	return &NullableGroupEntities{value: val, isSet: true}
}

func (v NullableGroupEntities) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroupEntities) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
