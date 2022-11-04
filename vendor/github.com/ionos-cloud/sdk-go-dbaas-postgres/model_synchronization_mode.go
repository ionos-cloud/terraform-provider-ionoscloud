/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
)

// SynchronizationMode Represents different modes of replication.
type SynchronizationMode string

// List of SynchronizationMode
const (
	ASYNCHRONOUS         SynchronizationMode = "ASYNCHRONOUS"
	SYNCHRONOUS          SynchronizationMode = "SYNCHRONOUS"
	STRICTLY_SYNCHRONOUS SynchronizationMode = "STRICTLY_SYNCHRONOUS"
)

func (v *SynchronizationMode) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := SynchronizationMode(value)
	for _, existing := range []SynchronizationMode{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid SynchronizationMode", value)
}

// Ptr returns reference to SynchronizationMode value
func (v SynchronizationMode) Ptr() *SynchronizationMode {
	return &v
}

type NullableSynchronizationMode struct {
	value *SynchronizationMode
	isSet bool
}

func (v NullableSynchronizationMode) Get() *SynchronizationMode {
	return v.value
}

func (v *NullableSynchronizationMode) Set(val *SynchronizationMode) {
	v.value = val
	v.isSet = true
}

func (v NullableSynchronizationMode) IsSet() bool {
	return v.isSet
}

func (v *NullableSynchronizationMode) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSynchronizationMode(val *SynchronizationMode) *NullableSynchronizationMode {
	return &NullableSynchronizationMode{value: val, isSet: true}
}

func (v NullableSynchronizationMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSynchronizationMode) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
