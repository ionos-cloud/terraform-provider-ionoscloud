/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// NatGatewayLanProperties struct for NatGatewayLanProperties
type NatGatewayLanProperties struct {
	// Id for the LAN connected to the NAT gateway
	Id *int32 `json:"id"`
	// Collection of gateway IP addresses of the NAT gateway. Will be auto-generated if not provided. Should ideally be an IP belonging to the same subnet as the LAN
	GatewayIps *[]string `json:"gatewayIps,omitempty"`
}


// GetId returns the Id field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *NatGatewayLanProperties) GetId() *int32 {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NatGatewayLanProperties) GetIdOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *NatGatewayLanProperties) SetId(v int32) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *NatGatewayLanProperties) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetGatewayIps returns the GatewayIps field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *NatGatewayLanProperties) GetGatewayIps() *[]string {
	if o == nil {
		return nil
	}


	return o.GatewayIps

}

// GetGatewayIpsOk returns a tuple with the GatewayIps field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NatGatewayLanProperties) GetGatewayIpsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}


	return o.GatewayIps, true
}

// SetGatewayIps sets field value
func (o *NatGatewayLanProperties) SetGatewayIps(v []string) {


	o.GatewayIps = &v

}

// HasGatewayIps returns a boolean if a field has been set.
func (o *NatGatewayLanProperties) HasGatewayIps() bool {
	if o != nil && o.GatewayIps != nil {
		return true
	}

	return false
}

func (o NatGatewayLanProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.GatewayIps != nil {
		toSerialize["gatewayIps"] = o.GatewayIps
	}
	return json.Marshal(toSerialize)
}
type NullableNatGatewayLanProperties struct {
	value *NatGatewayLanProperties
	isSet bool
}

func (v NullableNatGatewayLanProperties) Get() *NatGatewayLanProperties {
	return v.value
}

func (v *NullableNatGatewayLanProperties) Set(val *NatGatewayLanProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableNatGatewayLanProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableNatGatewayLanProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNatGatewayLanProperties(val *NatGatewayLanProperties) *NullableNatGatewayLanProperties {
	return &NullableNatGatewayLanProperties{value: val, isSet: true}
}

func (v NullableNatGatewayLanProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNatGatewayLanProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


