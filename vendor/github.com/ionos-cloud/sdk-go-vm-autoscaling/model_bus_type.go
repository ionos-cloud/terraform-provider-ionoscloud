/*
 * VM Auto Scaling API
 *
 * The VM Auto Scaling Service enables IONOS clients to horizontally scale the number of VM replicas based on configured rules. You can use VM Auto Scaling to ensure that you have a sufficient number of replicas to handle your application loads at all times.  For this purpose, create a VM Auto Scaling Group that contains the server replicas. The VM Auto Scaling Service ensures that the number of replicas in the group is always within the defined limits.   When scaling policies are set, VM Auto Scaling creates or deletes replicas according to the requirements of your applications. For each policy, specified 'scale-in' and 'scale-out' actions are performed when the corresponding thresholds are reached.
 *
 * API version: 1-SDK.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
)

// BusType The bus type of the volume. Default setting is 'VIRTIO'. The bus type 'IDE' is also supported.
type BusType string

// List of BusType
const (
	BUSTYPE_VIRTIO BusType = "VIRTIO"
	BUSTYPE_IDE    BusType = "IDE"
)

func (v *BusType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := BusType(value)
	for _, existing := range []BusType{"VIRTIO", "IDE"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid BusType", value)
}

// Ptr returns reference to BusType value
func (v BusType) Ptr() *BusType {
	return &v
}

type NullableBusType struct {
	value *BusType
	isSet bool
}

func (v NullableBusType) Get() *BusType {
	return v.value
}

func (v *NullableBusType) Set(val *BusType) {
	v.value = val
	v.isSet = true
}

func (v NullableBusType) IsSet() bool {
	return v.isSet
}

func (v *NullableBusType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBusType(val *BusType) *NullableBusType {
	return &NullableBusType{value: val, isSet: true}
}

func (v NullableBusType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBusType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
