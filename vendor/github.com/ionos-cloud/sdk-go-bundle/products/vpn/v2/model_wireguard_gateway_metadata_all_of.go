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

// checks if the WireguardGatewayMetadataAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardGatewayMetadataAllOf{}

// WireguardGatewayMetadataAllOf struct for WireguardGatewayMetadataAllOf
type WireguardGatewayMetadataAllOf struct {
	// Public key correspondng to the WireGuard Server private key.
	PublicKey string `json:"publicKey"`
}

// NewWireguardGatewayMetadataAllOf instantiates a new WireguardGatewayMetadataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardGatewayMetadataAllOf(publicKey string) *WireguardGatewayMetadataAllOf {
	this := WireguardGatewayMetadataAllOf{}

	this.PublicKey = publicKey

	return &this
}

// NewWireguardGatewayMetadataAllOfWithDefaults instantiates a new WireguardGatewayMetadataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardGatewayMetadataAllOfWithDefaults() *WireguardGatewayMetadataAllOf {
	this := WireguardGatewayMetadataAllOf{}
	return &this
}

// GetPublicKey returns the PublicKey field value
func (o *WireguardGatewayMetadataAllOf) GetPublicKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PublicKey
}

// GetPublicKeyOk returns a tuple with the PublicKey field value
// and a boolean to check if the value has been set.
func (o *WireguardGatewayMetadataAllOf) GetPublicKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PublicKey, true
}

// SetPublicKey sets field value
func (o *WireguardGatewayMetadataAllOf) SetPublicKey(v string) {
	o.PublicKey = v
}

func (o WireguardGatewayMetadataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WireguardGatewayMetadataAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.PublicKey) {
		toSerialize["publicKey"] = o.PublicKey
	}
	return toSerialize, nil
}

type NullableWireguardGatewayMetadataAllOf struct {
	value *WireguardGatewayMetadataAllOf
	isSet bool
}

func (v NullableWireguardGatewayMetadataAllOf) Get() *WireguardGatewayMetadataAllOf {
	return v.value
}

func (v *NullableWireguardGatewayMetadataAllOf) Set(val *WireguardGatewayMetadataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardGatewayMetadataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardGatewayMetadataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardGatewayMetadataAllOf(val *WireguardGatewayMetadataAllOf) *NullableWireguardGatewayMetadataAllOf {
	return &NullableWireguardGatewayMetadataAllOf{value: val, isSet: true}
}

func (v NullableWireguardGatewayMetadataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardGatewayMetadataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}