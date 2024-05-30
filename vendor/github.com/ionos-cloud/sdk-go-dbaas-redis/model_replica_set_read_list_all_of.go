/*
 * Redis DB API
 *
 * Redis Database API
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ReplicaSetReadListAllOf struct for ReplicaSetReadListAllOf
type ReplicaSetReadListAllOf struct {
	// ID of the ReplicaSet.
	Id *string `json:"id"`
	// The type of the resource.
	Type *string `json:"type"`
	// The URL of the ReplicaSet.
	Href *string `json:"href"`
	// The list of ReplicaSet resources.
	Items *[]ReplicaSetRead `json:"items,omitempty"`
}

// NewReplicaSetReadListAllOf instantiates a new ReplicaSetReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewReplicaSetReadListAllOf(id string, type_ string, href string) *ReplicaSetReadListAllOf {
	this := ReplicaSetReadListAllOf{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href

	return &this
}

// NewReplicaSetReadListAllOfWithDefaults instantiates a new ReplicaSetReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewReplicaSetReadListAllOfWithDefaults() *ReplicaSetReadListAllOf {
	this := ReplicaSetReadListAllOf{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ReplicaSetReadListAllOf) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaSetReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *ReplicaSetReadListAllOf) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ReplicaSetReadListAllOf) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ReplicaSetReadListAllOf) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaSetReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *ReplicaSetReadListAllOf) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ReplicaSetReadListAllOf) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ReplicaSetReadListAllOf) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaSetReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *ReplicaSetReadListAllOf) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *ReplicaSetReadListAllOf) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []ReplicaSetRead will be returned
func (o *ReplicaSetReadListAllOf) GetItems() *[]ReplicaSetRead {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ReplicaSetReadListAllOf) GetItemsOk() (*[]ReplicaSetRead, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *ReplicaSetReadListAllOf) SetItems(v []ReplicaSetRead) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *ReplicaSetReadListAllOf) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o ReplicaSetReadListAllOf) MarshalJSON() ([]byte, error) {
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

type NullableReplicaSetReadListAllOf struct {
	value *ReplicaSetReadListAllOf
	isSet bool
}

func (v NullableReplicaSetReadListAllOf) Get() *ReplicaSetReadListAllOf {
	return v.value
}

func (v *NullableReplicaSetReadListAllOf) Set(val *ReplicaSetReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableReplicaSetReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableReplicaSetReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableReplicaSetReadListAllOf(val *ReplicaSetReadListAllOf) *NullableReplicaSetReadListAllOf {
	return &NullableReplicaSetReadListAllOf{value: val, isSet: true}
}

func (v NullableReplicaSetReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableReplicaSetReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
