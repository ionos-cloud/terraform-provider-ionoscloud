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

	"fmt"
)

// CIDR struct for CIDR
type CIDR struct {
	string *string
}

// Unmarshal JSON data into any of the pointers in the struct
func (dst *CIDR) UnmarshalJSON(data []byte) error {
	var err error
	// try to unmarshal JSON data into string
	err = json.Unmarshal(data, &dst.string)
	if err == nil {
		jsonstring, _ := json.Marshal(dst.string)
		if string(jsonstring) == "{}" { // empty struct
			dst.string = nil
		} else {
			return nil // data stored in dst.string, return on the first match
		}
	} else {
		dst.string = nil
	}

	return fmt.Errorf("Data failed to match schemas in anyOf(CIDR)")
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src *CIDR) MarshalJSON() ([]byte, error) {
	if src.string != nil {
		return json.Marshal(&src.string)
	}

	return nil, nil // no data in anyOf schemas
}

type NullableCIDR struct {
	value *CIDR
	isSet bool
}

func (v NullableCIDR) Get() *CIDR {
	return v.value
}

func (v *NullableCIDR) Set(val *CIDR) {
	v.value = val
	v.isSet = true
}

func (v NullableCIDR) IsSet() bool {
	return v.isSet
}

func (v *NullableCIDR) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCIDR(val *CIDR) *NullableCIDR {
	return &NullableCIDR{value: val, isSet: true}
}

func (v NullableCIDR) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCIDR) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
