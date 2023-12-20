/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// GroupPropertiesDatacenter The VMs for this VM Auto Scaling Description are created in this virtual data center.
type GroupPropertiesDatacenter struct {
	// The unique resource identifier.
	Id *string `json:"id"`
	// The resource type.
	Type *string `json:"type,omitempty"`
	// The absolute URL to the resource's representation.
	Href *string `json:"href,omitempty"`
}

// NewGroupPropertiesDatacenter instantiates a new GroupPropertiesDatacenter object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGroupPropertiesDatacenter(id string) *GroupPropertiesDatacenter {
	this := GroupPropertiesDatacenter{}

	this.Id = &id

	return &this
}

// NewGroupPropertiesDatacenterWithDefaults instantiates a new GroupPropertiesDatacenter object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGroupPropertiesDatacenterWithDefaults() *GroupPropertiesDatacenter {
	this := GroupPropertiesDatacenter{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *GroupPropertiesDatacenter) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GroupPropertiesDatacenter) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *GroupPropertiesDatacenter) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *GroupPropertiesDatacenter) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *GroupPropertiesDatacenter) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GroupPropertiesDatacenter) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *GroupPropertiesDatacenter) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *GroupPropertiesDatacenter) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *GroupPropertiesDatacenter) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *GroupPropertiesDatacenter) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *GroupPropertiesDatacenter) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *GroupPropertiesDatacenter) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

func (o GroupPropertiesDatacenter) MarshalJSON() ([]byte, error) {
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

	return json.Marshal(toSerialize)
}

type NullableGroupPropertiesDatacenter struct {
	value *GroupPropertiesDatacenter
	isSet bool
}

func (v NullableGroupPropertiesDatacenter) Get() *GroupPropertiesDatacenter {
	return v.value
}

func (v *NullableGroupPropertiesDatacenter) Set(val *GroupPropertiesDatacenter) {
	v.value = val
	v.isSet = true
}

func (v NullableGroupPropertiesDatacenter) IsSet() bool {
	return v.isSet
}

func (v *NullableGroupPropertiesDatacenter) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGroupPropertiesDatacenter(val *GroupPropertiesDatacenter) *NullableGroupPropertiesDatacenter {
	return &NullableGroupPropertiesDatacenter{value: val, isSet: true}
}

func (v NullableGroupPropertiesDatacenter) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGroupPropertiesDatacenter) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
