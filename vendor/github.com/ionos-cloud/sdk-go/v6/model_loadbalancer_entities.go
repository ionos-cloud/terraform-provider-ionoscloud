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

// LoadbalancerEntities struct for LoadbalancerEntities
type LoadbalancerEntities struct {
	Balancednics *BalancedNics `json:"balancednics,omitempty"`
}



// GetBalancednics returns the Balancednics field value
// If the value is explicit nil, the zero value for BalancedNics will be returned
func (o *LoadbalancerEntities) GetBalancednics() *BalancedNics {
	if o == nil {
		return nil
	}


	return o.Balancednics

}

// GetBalancednicsOk returns a tuple with the Balancednics field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LoadbalancerEntities) GetBalancednicsOk() (*BalancedNics, bool) {
	if o == nil {
		return nil, false
	}


	return o.Balancednics, true
}

// SetBalancednics sets field value
func (o *LoadbalancerEntities) SetBalancednics(v BalancedNics) {


	o.Balancednics = &v

}

// HasBalancednics returns a boolean if a field has been set.
func (o *LoadbalancerEntities) HasBalancednics() bool {
	if o != nil && o.Balancednics != nil {
		return true
	}

	return false
}


func (o LoadbalancerEntities) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Balancednics != nil {
		toSerialize["balancednics"] = o.Balancednics
	}
	
	return json.Marshal(toSerialize)
}

type NullableLoadbalancerEntities struct {
	value *LoadbalancerEntities
	isSet bool
}

func (v NullableLoadbalancerEntities) Get() *LoadbalancerEntities {
	return v.value
}

func (v *NullableLoadbalancerEntities) Set(val *LoadbalancerEntities) {
	v.value = val
	v.isSet = true
}

func (v NullableLoadbalancerEntities) IsSet() bool {
	return v.isSet
}

func (v *NullableLoadbalancerEntities) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLoadbalancerEntities(val *LoadbalancerEntities) *NullableLoadbalancerEntities {
	return &NullableLoadbalancerEntities{value: val, isSet: true}
}

func (v NullableLoadbalancerEntities) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLoadbalancerEntities) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


