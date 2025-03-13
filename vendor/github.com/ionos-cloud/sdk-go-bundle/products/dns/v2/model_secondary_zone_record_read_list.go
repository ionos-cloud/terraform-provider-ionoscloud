/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.17.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the SecondaryZoneRecordReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SecondaryZoneRecordReadList{}

// SecondaryZoneRecordReadList struct for SecondaryZoneRecordReadList
type SecondaryZoneRecordReadList struct {
	// The resource's unique identifier.
	Id       string                              `json:"id"`
	Type     string                              `json:"type"`
	Href     string                              `json:"href"`
	Metadata SecondaryZoneRecordReadListMetadata `json:"metadata"`
	Items    []SecondaryZoneRecordRead           `json:"items"`
	// Pagination offset.
	Offset float32 `json:"offset"`
	// Pagination limit.
	Limit float32 `json:"limit"`
	Links Links   `json:"_links"`
}

// NewSecondaryZoneRecordReadList instantiates a new SecondaryZoneRecordReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSecondaryZoneRecordReadList(id string, type_ string, href string, metadata SecondaryZoneRecordReadListMetadata, items []SecondaryZoneRecordRead, offset float32, limit float32, links Links) *SecondaryZoneRecordReadList {
	this := SecondaryZoneRecordReadList{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Items = items
	this.Offset = offset
	this.Limit = limit
	this.Links = links

	return &this
}

// NewSecondaryZoneRecordReadListWithDefaults instantiates a new SecondaryZoneRecordReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSecondaryZoneRecordReadListWithDefaults() *SecondaryZoneRecordReadList {
	this := SecondaryZoneRecordReadList{}
	return &this
}

// GetId returns the Id field value
func (o *SecondaryZoneRecordReadList) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *SecondaryZoneRecordReadList) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *SecondaryZoneRecordReadList) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *SecondaryZoneRecordReadList) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *SecondaryZoneRecordReadList) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *SecondaryZoneRecordReadList) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *SecondaryZoneRecordReadList) GetMetadata() SecondaryZoneRecordReadListMetadata {
	if o == nil {
		var ret SecondaryZoneRecordReadListMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetMetadataOk() (*SecondaryZoneRecordReadListMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *SecondaryZoneRecordReadList) SetMetadata(v SecondaryZoneRecordReadListMetadata) {
	o.Metadata = v
}

// GetItems returns the Items field value
func (o *SecondaryZoneRecordReadList) GetItems() []SecondaryZoneRecordRead {
	if o == nil {
		var ret []SecondaryZoneRecordRead
		return ret
	}

	return o.Items
}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetItemsOk() ([]SecondaryZoneRecordRead, bool) {
	if o == nil {
		return nil, false
	}
	return o.Items, true
}

// SetItems sets field value
func (o *SecondaryZoneRecordReadList) SetItems(v []SecondaryZoneRecordRead) {
	o.Items = v
}

// GetOffset returns the Offset field value
func (o *SecondaryZoneRecordReadList) GetOffset() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Offset
}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetOffsetOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Offset, true
}

// SetOffset sets field value
func (o *SecondaryZoneRecordReadList) SetOffset(v float32) {
	o.Offset = v
}

// GetLimit returns the Limit field value
func (o *SecondaryZoneRecordReadList) GetLimit() float32 {
	if o == nil {
		var ret float32
		return ret
	}

	return o.Limit
}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetLimitOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Limit, true
}

// SetLimit sets field value
func (o *SecondaryZoneRecordReadList) SetLimit(v float32) {
	o.Limit = v
}

// GetLinks returns the Links field value
func (o *SecondaryZoneRecordReadList) GetLinks() Links {
	if o == nil {
		var ret Links
		return ret
	}

	return o.Links
}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
func (o *SecondaryZoneRecordReadList) GetLinksOk() (*Links, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Links, true
}

// SetLinks sets field value
func (o *SecondaryZoneRecordReadList) SetLinks(v Links) {
	o.Links = v
}

func (o SecondaryZoneRecordReadList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SecondaryZoneRecordReadList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	toSerialize["metadata"] = o.Metadata
	toSerialize["items"] = o.Items
	toSerialize["offset"] = o.Offset
	toSerialize["limit"] = o.Limit
	toSerialize["_links"] = o.Links
	return toSerialize, nil
}

type NullableSecondaryZoneRecordReadList struct {
	value *SecondaryZoneRecordReadList
	isSet bool
}

func (v NullableSecondaryZoneRecordReadList) Get() *SecondaryZoneRecordReadList {
	return v.value
}

func (v *NullableSecondaryZoneRecordReadList) Set(val *SecondaryZoneRecordReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableSecondaryZoneRecordReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableSecondaryZoneRecordReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSecondaryZoneRecordReadList(val *SecondaryZoneRecordReadList) *NullableSecondaryZoneRecordReadList {
	return &NullableSecondaryZoneRecordReadList{value: val, isSet: true}
}

func (v NullableSecondaryZoneRecordReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSecondaryZoneRecordReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
