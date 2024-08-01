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

// checks if the WireguardPeerRead type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardPeerRead{}

// WireguardPeerRead struct for WireguardPeerRead
type WireguardPeerRead struct {
	// The ID (UUID) of the WireguardPeer.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the WireguardPeer.
	Href       string                `json:"href"`
	Metadata   WireguardPeerMetadata `json:"metadata"`
	Properties WireguardPeer         `json:"properties"`
}

// NewWireguardPeerRead instantiates a new WireguardPeerRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardPeerRead(id string, type_ string, href string, metadata WireguardPeerMetadata, properties WireguardPeer) *WireguardPeerRead {
	this := WireguardPeerRead{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewWireguardPeerReadWithDefaults instantiates a new WireguardPeerRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardPeerReadWithDefaults() *WireguardPeerRead {
	this := WireguardPeerRead{}
	return &this
}

// GetId returns the Id field value
func (o *WireguardPeerRead) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WireguardPeerRead) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *WireguardPeerRead) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *WireguardPeerRead) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *WireguardPeerRead) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *WireguardPeerRead) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *WireguardPeerRead) GetMetadata() WireguardPeerMetadata {
	if o == nil {
		var ret WireguardPeerMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerRead) GetMetadataOk() (*WireguardPeerMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *WireguardPeerRead) SetMetadata(v WireguardPeerMetadata) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *WireguardPeerRead) GetProperties() WireguardPeer {
	if o == nil {
		var ret WireguardPeer
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerRead) GetPropertiesOk() (*WireguardPeer, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *WireguardPeerRead) SetProperties(v WireguardPeer) {
	o.Properties = v
}

func (o WireguardPeerRead) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WireguardPeerRead) ToMap() (map[string]interface{}, error) {
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
	if !IsZero(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsZero(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableWireguardPeerRead struct {
	value *WireguardPeerRead
	isSet bool
}

func (v NullableWireguardPeerRead) Get() *WireguardPeerRead {
	return v.value
}

func (v *NullableWireguardPeerRead) Set(val *WireguardPeerRead) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardPeerRead) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardPeerRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardPeerRead(val *WireguardPeerRead) *NullableWireguardPeerRead {
	return &NullableWireguardPeerRead{value: val, isSet: true}
}

func (v NullableWireguardPeerRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardPeerRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
