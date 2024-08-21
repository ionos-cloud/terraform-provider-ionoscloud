/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 0.1.7
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the DistributionsAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DistributionsAllOf{}

// DistributionsAllOf struct for DistributionsAllOf
type DistributionsAllOf struct {
	// ID of the list of Distribution resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of Distribution resources.
	Href string `json:"href"`
	// The list of Distribution resources.
	Items []Distribution `json:"items,omitempty"`
}

// NewDistributionsAllOf instantiates a new DistributionsAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDistributionsAllOf(id string, type_ string, href string) *DistributionsAllOf {
	this := DistributionsAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewDistributionsAllOfWithDefaults instantiates a new DistributionsAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDistributionsAllOfWithDefaults() *DistributionsAllOf {
	this := DistributionsAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *DistributionsAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *DistributionsAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *DistributionsAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *DistributionsAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *DistributionsAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *DistributionsAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *DistributionsAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *DistributionsAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *DistributionsAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *DistributionsAllOf) GetItems() []Distribution {
	if o == nil || IsNil(o.Items) {
		var ret []Distribution
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DistributionsAllOf) GetItemsOk() ([]Distribution, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *DistributionsAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []Distribution and assigns it to the Items field.
func (o *DistributionsAllOf) SetItems(v []Distribution) {
	o.Items = v
}

func (o DistributionsAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableDistributionsAllOf struct {
	value *DistributionsAllOf
	isSet bool
}

func (v NullableDistributionsAllOf) Get() *DistributionsAllOf {
	return v.value
}

func (v *NullableDistributionsAllOf) Set(val *DistributionsAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableDistributionsAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableDistributionsAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDistributionsAllOf(val *DistributionsAllOf) *NullableDistributionsAllOf {
	return &NullableDistributionsAllOf{value: val, isSet: true}
}

func (v NullableDistributionsAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDistributionsAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
