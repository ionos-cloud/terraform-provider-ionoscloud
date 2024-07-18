/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// checks if the WireguardPeerEnsure type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardPeerEnsure{}

// WireguardPeerEnsure struct for WireguardPeerEnsure
type WireguardPeerEnsure struct {
	// The ID (UUID) of the WireguardPeer.
	Id string `json:"id"`
	// Metadata
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Properties WireguardPeer          `json:"properties"`
}

// NewWireguardPeerEnsure instantiates a new WireguardPeerEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardPeerEnsure(id string, properties WireguardPeer) *WireguardPeerEnsure {
	this := WireguardPeerEnsure{}

	this.Id = id
	this.Properties = properties

	return &this
}

// NewWireguardPeerEnsureWithDefaults instantiates a new WireguardPeerEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardPeerEnsureWithDefaults() *WireguardPeerEnsure {
	this := WireguardPeerEnsure{}
	return &this
}

// GetId returns the Id field value
func (o *WireguardPeerEnsure) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *WireguardPeerEnsure) SetId(v string) {
	o.Id = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *WireguardPeerEnsure) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardPeerEnsure) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *WireguardPeerEnsure) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *WireguardPeerEnsure) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *WireguardPeerEnsure) GetProperties() WireguardPeer {
	if o == nil {
		var ret WireguardPeer
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *WireguardPeerEnsure) GetPropertiesOk() (*WireguardPeer, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *WireguardPeerEnsure) SetProperties(v WireguardPeer) {
	o.Properties = v
}

func (o WireguardPeerEnsure) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WireguardPeerEnsure) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsZero(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableWireguardPeerEnsure struct {
	value *WireguardPeerEnsure
	isSet bool
}

func (v NullableWireguardPeerEnsure) Get() *WireguardPeerEnsure {
	return v.value
}

func (v *NullableWireguardPeerEnsure) Set(val *WireguardPeerEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardPeerEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardPeerEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardPeerEnsure(val *WireguardPeerEnsure) *NullableWireguardPeerEnsure {
	return &NullableWireguardPeerEnsure{value: val, isSet: true}
}

func (v NullableWireguardPeerEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardPeerEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
