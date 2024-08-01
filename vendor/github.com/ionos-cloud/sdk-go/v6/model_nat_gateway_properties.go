/*
 * CLOUD API
 *
 *  IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// NatGatewayProperties struct for NatGatewayProperties
type NatGatewayProperties struct {
	// Name of the NAT Gateway.
	Name *string `json:"name"`
	// Collection of public IP addresses of the NAT Gateway. Should be customer reserved IP addresses in that location.
	PublicIps *[]string `json:"publicIps"`
	// Collection of LANs connected to the NAT Gateway. IPs must contain a valid subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
	Lans *[]NatGatewayLanProperties `json:"lans,omitempty"`
}

// NewNatGatewayProperties instantiates a new NatGatewayProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNatGatewayProperties(name string, publicIps []string) *NatGatewayProperties {
	this := NatGatewayProperties{}

	this.Name = &name
	this.PublicIps = &publicIps

	return &this
}

// NewNatGatewayPropertiesWithDefaults instantiates a new NatGatewayProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNatGatewayPropertiesWithDefaults() *NatGatewayProperties {
	this := NatGatewayProperties{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, nil is returned
func (o *NatGatewayProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NatGatewayProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *NatGatewayProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *NatGatewayProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetPublicIps returns the PublicIps field value
// If the value is explicit nil, nil is returned
func (o *NatGatewayProperties) GetPublicIps() *[]string {
	if o == nil {
		return nil
	}

	return o.PublicIps

}

// GetPublicIpsOk returns a tuple with the PublicIps field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NatGatewayProperties) GetPublicIpsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.PublicIps, true
}

// SetPublicIps sets field value
func (o *NatGatewayProperties) SetPublicIps(v []string) {

	o.PublicIps = &v

}

// HasPublicIps returns a boolean if a field has been set.
func (o *NatGatewayProperties) HasPublicIps() bool {
	if o != nil && o.PublicIps != nil {
		return true
	}

	return false
}

// GetLans returns the Lans field value
// If the value is explicit nil, nil is returned
func (o *NatGatewayProperties) GetLans() *[]NatGatewayLanProperties {
	if o == nil {
		return nil
	}

	return o.Lans

}

// GetLansOk returns a tuple with the Lans field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NatGatewayProperties) GetLansOk() (*[]NatGatewayLanProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Lans, true
}

// SetLans sets field value
func (o *NatGatewayProperties) SetLans(v []NatGatewayLanProperties) {

	o.Lans = &v

}

// HasLans returns a boolean if a field has been set.
func (o *NatGatewayProperties) HasLans() bool {
	if o != nil && o.Lans != nil {
		return true
	}

	return false
}

func (o NatGatewayProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.PublicIps != nil {
		toSerialize["publicIps"] = o.PublicIps
	}

	if o.Lans != nil {
		toSerialize["lans"] = o.Lans
	}

	return json.Marshal(toSerialize)
}

type NullableNatGatewayProperties struct {
	value *NatGatewayProperties
	isSet bool
}

func (v NullableNatGatewayProperties) Get() *NatGatewayProperties {
	return v.value
}

func (v *NullableNatGatewayProperties) Set(val *NatGatewayProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableNatGatewayProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableNatGatewayProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNatGatewayProperties(val *NatGatewayProperties) *NullableNatGatewayProperties {
	return &NullableNatGatewayProperties{value: val, isSet: true}
}

func (v NullableNatGatewayProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNatGatewayProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
