/*
 * VM Auto Scaling service (CloudAPI)
 *
 * VM Auto Scaling service enables IONOS clients to horizontally scale the number of VM instances, based on configured rules. Use Auto Scaling to ensure you will have a sufficient number of instances to handle your application loads at all times.  Create an Auto Scaling group that contains the server instances; Auto Scaling service will ensure that the number of instances in the group is always within these limits.  When target replica count is specified, Auto Scaling will maintain the set number on instances.  When scaling policies are specified, Auto Scaling will create or delete instances based on the demands of your applications. For each policy, specified scale-in and scale-out actions are performed whenever the corresponding thresholds are met.
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

// QueryUnit Units of the applied Metric.
type QueryUnit string

// List of QueryUnit
const (
	PER_HOUR   QueryUnit = "PER_HOUR"
	PER_MINUTE QueryUnit = "PER_MINUTE"
	PER_SECOND QueryUnit = "PER_SECOND"
	TOTAL      QueryUnit = "TOTAL"
)

func (v *QueryUnit) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := QueryUnit(value)
	for _, existing := range []QueryUnit{"PER_HOUR", "PER_MINUTE", "PER_SECOND", "TOTAL"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid QueryUnit", value)
}

// Ptr returns reference to QueryUnit value
func (v QueryUnit) Ptr() *QueryUnit {
	return &v
}

type NullableQueryUnit struct {
	value *QueryUnit
	isSet bool
}

func (v NullableQueryUnit) Get() *QueryUnit {
	return v.value
}

func (v *NullableQueryUnit) Set(val *QueryUnit) {
	v.value = val
	v.isSet = true
}

func (v NullableQueryUnit) IsSet() bool {
	return v.isSet
}

func (v *NullableQueryUnit) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQueryUnit(val *QueryUnit) *NullableQueryUnit {
	return &NullableQueryUnit{value: val, isSet: true}
}

func (v NullableQueryUnit) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQueryUnit) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
