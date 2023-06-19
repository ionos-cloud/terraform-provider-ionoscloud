/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ResourceGroups Resources assigned to this group.
type ResourceGroups struct {
	// URL to the object representation (absolute path).
	Href *string `json:"href,omitempty"`
	// The resource's unique identifier.
	Id *string `json:"id,omitempty"`
	// Array of items in the collection.
	Items *[]Resource `json:"items,omitempty"`
	// The type of the resource.
	Type *Type `json:"type,omitempty"`
}

// NewResourceGroups instantiates a new ResourceGroups object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResourceGroups() *ResourceGroups {
	this := ResourceGroups{}

	return &this
}

// NewResourceGroupsWithDefaults instantiates a new ResourceGroups object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResourceGroupsWithDefaults() *ResourceGroups {
	this := ResourceGroups{}
	return &this
}

// GetHref returns the Href field value
// If the value is explicit nil, nil is returned
func (o *ResourceGroups) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceGroups) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *ResourceGroups) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *ResourceGroups) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetId returns the Id field value
// If the value is explicit nil, nil is returned
func (o *ResourceGroups) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceGroups) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *ResourceGroups) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ResourceGroups) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, nil is returned
func (o *ResourceGroups) GetItems() *[]Resource {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceGroups) GetItemsOk() (*[]Resource, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *ResourceGroups) SetItems(v []Resource) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *ResourceGroups) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, nil is returned
func (o *ResourceGroups) GetType() *Type {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceGroups) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *ResourceGroups) SetType(v Type) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ResourceGroups) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

func (o ResourceGroups) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	return json.Marshal(toSerialize)
}

type NullableResourceGroups struct {
	value *ResourceGroups
	isSet bool
}

func (v NullableResourceGroups) Get() *ResourceGroups {
	return v.value
}

func (v *NullableResourceGroups) Set(val *ResourceGroups) {
	v.value = val
	v.isSet = true
}

func (v NullableResourceGroups) IsSet() bool {
	return v.isSet
}

func (v *NullableResourceGroups) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResourceGroups(val *ResourceGroups) *NullableResourceGroups {
	return &NullableResourceGroups{value: val, isSet: true}
}

func (v NullableResourceGroups) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResourceGroups) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
