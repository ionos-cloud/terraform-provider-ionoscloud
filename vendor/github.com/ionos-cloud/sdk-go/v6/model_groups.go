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

// Groups struct for Groups
type Groups struct {
	// The resource's unique identifier.
	Id *string `json:"id,omitempty"`
	// The type of the resource.
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path).
	Href *string `json:"href,omitempty"`
	// Array of items in the collection.
	Items *[]Group `json:"items,omitempty"`
}

// NewGroups instantiates a new Groups object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroups() *Groups {
	this := Groups{}

	return &this
}

// NewGroupsWithDefaults instantiates a new Groups object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupsWithDefaults() *Groups {
	this := Groups{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, nil is returned
func (o *Groups) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Groups) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *Groups) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *Groups) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, nil is returned
func (o *Groups) GetType() *Type {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Groups) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *Groups) SetType(v Type) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Groups) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, nil is returned
func (o *Groups) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Groups) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *Groups) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *Groups) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, nil is returned
func (o *Groups) GetItems() *[]Group {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Groups) GetItemsOk() (*[]Group, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *Groups) SetItems(v []Group) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *Groups) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o Groups) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	return json.Marshal(toSerialize)
}

type NullableGroups struct {
	value *Groups
	isSet bool
}

func (v NullableGroups) Get() *Groups {
	return v.value
}

func (v *NullableGroups) Set(val *Groups) {
	v.value = val
	v.isSet = true
}

func (v NullableGroups) IsSet() bool {
	return v.isSet
}

func (v *NullableGroups) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroups(val *Groups) *NullableGroups {
	return &NullableGroups{value: val, isSet: true}
}

func (v NullableGroups) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroups) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
