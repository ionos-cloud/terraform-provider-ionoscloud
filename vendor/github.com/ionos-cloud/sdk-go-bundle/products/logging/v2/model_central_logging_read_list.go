/*
 * IONOS Logging Service REST API
 *
 * The Logging Service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data. This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the CentralLoggingReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CentralLoggingReadList{}

// CentralLoggingReadList struct for CentralLoggingReadList
type CentralLoggingReadList struct {
	// ID of the list of CentralLogging resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of CentralLogging resources.
	Href string `json:"href"`
	// The list of CentralLogging resources.
	Items []CentralLoggingRead `json:"items,omitempty"`
}

// NewCentralLoggingReadList instantiates a new CentralLoggingReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCentralLoggingReadList(id string, type_ string, href string) *CentralLoggingReadList {
	this := CentralLoggingReadList{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewCentralLoggingReadListWithDefaults instantiates a new CentralLoggingReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCentralLoggingReadListWithDefaults() *CentralLoggingReadList {
	this := CentralLoggingReadList{}
	return &this
}

// GetId returns the Id field value
func (o *CentralLoggingReadList) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *CentralLoggingReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *CentralLoggingReadList) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *CentralLoggingReadList) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *CentralLoggingReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *CentralLoggingReadList) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *CentralLoggingReadList) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *CentralLoggingReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *CentralLoggingReadList) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *CentralLoggingReadList) GetItems() []CentralLoggingRead {
	if o == nil || IsNil(o.Items) {
		var ret []CentralLoggingRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CentralLoggingReadList) GetItemsOk() ([]CentralLoggingRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *CentralLoggingReadList) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []CentralLoggingRead and assigns it to the Items field.
func (o *CentralLoggingReadList) SetItems(v []CentralLoggingRead) {
	o.Items = v
}

func (o CentralLoggingReadList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CentralLoggingReadList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableCentralLoggingReadList struct {
	value *CentralLoggingReadList
	isSet bool
}

func (v NullableCentralLoggingReadList) Get() *CentralLoggingReadList {
	return v.value
}

func (v *NullableCentralLoggingReadList) Set(val *CentralLoggingReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableCentralLoggingReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableCentralLoggingReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCentralLoggingReadList(val *CentralLoggingReadList) *NullableCentralLoggingReadList {
	return &NullableCentralLoggingReadList{value: val, isSet: true}
}

func (v NullableCentralLoggingReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCentralLoggingReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
