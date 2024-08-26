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

// checks if the IPSecTunnelReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecTunnelReadListAllOf{}

// IPSecTunnelReadListAllOf struct for IPSecTunnelReadListAllOf
type IPSecTunnelReadListAllOf struct {
	// ID of the list of IPSecTunnel resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of IPSecTunnel resources.
	Href string `json:"href"`
	// The list of IPSecTunnel resources.
	Items []IPSecTunnelRead `json:"items,omitempty"`
}

// NewIPSecTunnelReadListAllOf instantiates a new IPSecTunnelReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecTunnelReadListAllOf(id string, type_ string, href string) *IPSecTunnelReadListAllOf {
	this := IPSecTunnelReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewIPSecTunnelReadListAllOfWithDefaults instantiates a new IPSecTunnelReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecTunnelReadListAllOfWithDefaults() *IPSecTunnelReadListAllOf {
	this := IPSecTunnelReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *IPSecTunnelReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnelReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *IPSecTunnelReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *IPSecTunnelReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnelReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *IPSecTunnelReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *IPSecTunnelReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnelReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *IPSecTunnelReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *IPSecTunnelReadListAllOf) GetItems() []IPSecTunnelRead {
	if o == nil || IsNil(o.Items) {
		var ret []IPSecTunnelRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecTunnelReadListAllOf) GetItemsOk() ([]IPSecTunnelRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *IPSecTunnelReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []IPSecTunnelRead and assigns it to the Items field.
func (o *IPSecTunnelReadListAllOf) SetItems(v []IPSecTunnelRead) {
	o.Items = v
}

func (o IPSecTunnelReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableIPSecTunnelReadListAllOf struct {
	value *IPSecTunnelReadListAllOf
	isSet bool
}

func (v NullableIPSecTunnelReadListAllOf) Get() *IPSecTunnelReadListAllOf {
	return v.value
}

func (v *NullableIPSecTunnelReadListAllOf) Set(val *IPSecTunnelReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecTunnelReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecTunnelReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecTunnelReadListAllOf(val *IPSecTunnelReadListAllOf) *NullableIPSecTunnelReadListAllOf {
	return &NullableIPSecTunnelReadListAllOf{value: val, isSet: true}
}

func (v NullableIPSecTunnelReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecTunnelReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
