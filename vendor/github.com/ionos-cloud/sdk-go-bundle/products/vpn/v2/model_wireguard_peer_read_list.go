/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vpn

import (
	"encoding/json"
)

// checks if the WireguardPeerReadList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardPeerReadList{}

// WireguardPeerReadList struct for WireguardPeerReadList
type WireguardPeerReadList struct {
	// ID of the list of WireguardPeer resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of WireguardPeer resources.
	Href string `json:"href"`
	// The list of WireguardPeer resources.
	Items []WireguardPeerRead `json:"items,omitempty"`
	// The offset specified in the request (if none was specified, the default offset is 0).
	Offset int32 `json:"offset"`
	// The limit specified in the request (if none was specified, use the endpoint's default pagination limit).
	Limit int32 `json:"limit"`
	Links Links `json:"_links"`
}

// NewWireguardPeerReadList instantiates a new WireguardPeerReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardPeerReadList(id string, type_ string, href string, offset int32, limit int32, links Links) *WireguardPeerReadList {
	this := WireguardPeerReadList{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Offset = offset
	this.Limit = limit
	this.Links = links

	return &this
}

// NewWireguardPeerReadListWithDefaults instantiates a new WireguardPeerReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardPeerReadListWithDefaults() *WireguardPeerReadList {
	this := WireguardPeerReadList{}
	return &this
}

// GetId returns the Id field value
func (o *WireguardPeerReadList) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WireguardPeerReadList) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *WireguardPeerReadList) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WireguardPeerReadList) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *WireguardPeerReadList) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *WireguardPeerReadList) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *WireguardPeerReadList) GetItems() []WireguardPeerRead {
	if o == nil || IsNil(o.Items) {
		var ret []WireguardPeerRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetItemsOk() ([]WireguardPeerRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *WireguardPeerReadList) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []WireguardPeerRead and assigns it to the Items field.
func (o *WireguardPeerReadList) SetItems(v []WireguardPeerRead) {
	o.Items = v
}

// GetOffset returns the Offset field value
func (o *WireguardPeerReadList) GetOffset() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Offset
}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Offset, true
}

// SetOffset sets field value
func (o *WireguardPeerReadList) SetOffset(v int32) {
	o.Offset = v
}

// GetLimit returns the Limit field value
func (o *WireguardPeerReadList) GetLimit() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Limit
}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetLimitOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Limit, true
}

// SetLimit sets field value
func (o *WireguardPeerReadList) SetLimit(v int32) {
	o.Limit = v
}

// GetLinks returns the Links field value
func (o *WireguardPeerReadList) GetLinks() Links {
	if o == nil {
		var ret Links
		return ret
	}

	return o.Links
}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadList) GetLinksOk() (*Links, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Links, true
}

// SetLinks sets field value
func (o *WireguardPeerReadList) SetLinks(v Links) {
	o.Links = v
}

func (o WireguardPeerReadList) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WireguardPeerReadList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsZero(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsZero(o.Href) {
		toSerialize["href"] = o.Href
	}
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	if !IsZero(o.Offset) {
		toSerialize["offset"] = o.Offset
	}
	if !IsZero(o.Limit) {
		toSerialize["limit"] = o.Limit
	}
	if !IsZero(o.Links) {
		toSerialize["_links"] = o.Links
	}
	return toSerialize, nil
}

type NullableWireguardPeerReadList struct {
	value *WireguardPeerReadList
	isSet bool
}

func (v NullableWireguardPeerReadList) Get() *WireguardPeerReadList {
	return v.value
}

func (v *NullableWireguardPeerReadList) Set(val *WireguardPeerReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardPeerReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardPeerReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardPeerReadList(val *WireguardPeerReadList) *NullableWireguardPeerReadList {
	return &NullableWireguardPeerReadList{value: val, isSet: true}
}

func (v NullableWireguardPeerReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardPeerReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
