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

// checks if the IPSecTunnelAuth type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecTunnelAuth{}

// IPSecTunnelAuth Properties with all data needed to define IPSec Authentication.
type IPSecTunnelAuth struct {
	// The Authentication Method to use for IPSec Authentication.\\ Options:   - PSK
	Method string    `json:"method"`
	Psk    *IPSecPSK `json:"psk,omitempty"`
}

// NewIPSecTunnelAuth instantiates a new IPSecTunnelAuth object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecTunnelAuth(method string) *IPSecTunnelAuth {
	this := IPSecTunnelAuth{}

	this.Method = method

	return &this
}

// NewIPSecTunnelAuthWithDefaults instantiates a new IPSecTunnelAuth object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecTunnelAuthWithDefaults() *IPSecTunnelAuth {
	this := IPSecTunnelAuth{}
	var method string = "PSK"
	this.Method = method
	return &this
}

// GetMethod returns the Method field value
func (o *IPSecTunnelAuth) GetMethod() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Method
}

// GetMethodOk returns a tuple with the Method field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnelAuth) GetMethodOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Method, true
}

// SetMethod sets field value
func (o *IPSecTunnelAuth) SetMethod(v string) {
	o.Method = v
}

// GetPsk returns the Psk field value if set, zero value otherwise.
func (o *IPSecTunnelAuth) GetPsk() IPSecPSK {
	if o == nil || IsNil(o.Psk) {
		var ret IPSecPSK
		return ret
	}
	return *o.Psk
}

// GetPskOk returns a tuple with the Psk field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecTunnelAuth) GetPskOk() (*IPSecPSK, bool) {
	if o == nil || IsNil(o.Psk) {
		return nil, false
	}
	return o.Psk, true
}

// HasPsk returns a boolean if a field has been set.
func (o *IPSecTunnelAuth) HasPsk() bool {
	if o != nil && !IsNil(o.Psk) {
		return true
	}

	return false
}

// SetPsk gets a reference to the given IPSecPSK and assigns it to the Psk field.
func (o *IPSecTunnelAuth) SetPsk(v IPSecPSK) {
	o.Psk = &v
}

func (o IPSecTunnelAuth) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecTunnelAuth) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Method) {
		toSerialize["method"] = o.Method
	}
	if !IsNil(o.Psk) {
		toSerialize["psk"] = o.Psk
	}
	return toSerialize, nil
}

type NullableIPSecTunnelAuth struct {
	value *IPSecTunnelAuth
	isSet bool
}

func (v NullableIPSecTunnelAuth) Get() *IPSecTunnelAuth {
	return v.value
}

func (v *NullableIPSecTunnelAuth) Set(val *IPSecTunnelAuth) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecTunnelAuth) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecTunnelAuth) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecTunnelAuth(val *IPSecTunnelAuth) *NullableIPSecTunnelAuth {
	return &NullableIPSecTunnelAuth{value: val, isSet: true}
}

func (v NullableIPSecTunnelAuth) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecTunnelAuth) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}