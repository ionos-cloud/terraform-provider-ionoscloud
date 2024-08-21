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

// checks if the WireguardPeerReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardPeerReadListAllOf{}

// WireguardPeerReadListAllOf struct for WireguardPeerReadListAllOf
type WireguardPeerReadListAllOf struct {
	// ID of the list of WireguardPeer resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of WireguardPeer resources.
	Href string `json:"href"`
	// The list of WireguardPeer resources.
	Items []WireguardPeerRead `json:"items,omitempty"`
}

// NewWireguardPeerReadListAllOf instantiates a new WireguardPeerReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardPeerReadListAllOf(id string, type_ string, href string) *WireguardPeerReadListAllOf {
	this := WireguardPeerReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewWireguardPeerReadListAllOfWithDefaults instantiates a new WireguardPeerReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardPeerReadListAllOfWithDefaults() *WireguardPeerReadListAllOf {
	this := WireguardPeerReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *WireguardPeerReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WireguardPeerReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *WireguardPeerReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WireguardPeerReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *WireguardPeerReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *WireguardPeerReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *WireguardPeerReadListAllOf) GetItems() []WireguardPeerRead {
	if o == nil || IsNil(o.Items) {
		var ret []WireguardPeerRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardPeerReadListAllOf) GetItemsOk() ([]WireguardPeerRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *WireguardPeerReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []WireguardPeerRead and assigns it to the Items field.
func (o *WireguardPeerReadListAllOf) SetItems(v []WireguardPeerRead) {
	o.Items = v
}

func (o WireguardPeerReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableWireguardPeerReadListAllOf struct {
	value *WireguardPeerReadListAllOf
	isSet bool
}

func (v NullableWireguardPeerReadListAllOf) Get() *WireguardPeerReadListAllOf {
	return v.value
}

func (v *NullableWireguardPeerReadListAllOf) Set(val *WireguardPeerReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardPeerReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardPeerReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardPeerReadListAllOf(val *WireguardPeerReadListAllOf) *NullableWireguardPeerReadListAllOf {
	return &NullableWireguardPeerReadListAllOf{value: val, isSet: true}
}

func (v NullableWireguardPeerReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardPeerReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
