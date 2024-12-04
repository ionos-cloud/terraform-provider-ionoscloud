/*
 * IONOS Cloud VPN Gateway API
 *
 * The Managed VPN Gateway service provides secure and scalable connectivity, enabling encrypted communication between your IONOS cloud resources in a VDC and remote networks (on-premises, multi-cloud, private LANs in other VDCs etc).
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vpn

import (
	"encoding/json"
)

// checks if the WireguardEndpoint type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardEndpoint{}

// WireguardEndpoint Properties with all data needed to create a new WireGuard Gateway endpoint.
type WireguardEndpoint struct {
	// Hostname or IPV4 address that the WireGuard Server will connect to.
	Host string `json:"host"`
	// IP port number
	Port *int32 `json:"port,omitempty"`
}

// NewWireguardEndpoint instantiates a new WireguardEndpoint object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardEndpoint(host string) *WireguardEndpoint {
	this := WireguardEndpoint{}

	this.Host = host
	var port int32 = 51820
	this.Port = &port

	return &this
}

// NewWireguardEndpointWithDefaults instantiates a new WireguardEndpoint object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardEndpointWithDefaults() *WireguardEndpoint {
	this := WireguardEndpoint{}
	var port int32 = 51820
	this.Port = &port
	return &this
}

// GetHost returns the Host field value
func (o *WireguardEndpoint) GetHost() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Host
}

// GetHostOk returns a tuple with the Host field value
// and a boolean to check if the value has been set.
func (o *WireguardEndpoint) GetHostOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Host, true
}

// SetHost sets field value
func (o *WireguardEndpoint) SetHost(v string) {
	o.Host = v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *WireguardEndpoint) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardEndpoint) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *WireguardEndpoint) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *WireguardEndpoint) SetPort(v int32) {
	o.Port = &v
}

func (o WireguardEndpoint) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["host"] = o.Host
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}
	return toSerialize, nil
}

type NullableWireguardEndpoint struct {
	value *WireguardEndpoint
	isSet bool
}

func (v NullableWireguardEndpoint) Get() *WireguardEndpoint {
	return v.value
}

func (v *NullableWireguardEndpoint) Set(val *WireguardEndpoint) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardEndpoint) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardEndpoint) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardEndpoint(val *WireguardEndpoint) *NullableWireguardEndpoint {
	return &NullableWireguardEndpoint{value: val, isSet: true}
}

func (v NullableWireguardEndpoint) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardEndpoint) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
