/*
 * Redis DB API
 *
 * Redis Database API
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

// UserPassword - struct for UserPassword
type UserPassword struct {
	HashedPassword *HashedPassword
	string         *string
}

// HashedPasswordAsUserPassword is a convenience function that returns HashedPassword wrapped in UserPassword
func HashedPasswordAsUserPassword(v *HashedPassword) UserPassword {
	return UserPassword{HashedPassword: v}
}

// stringAsUserPassword is a convenience function that returns string wrapped in UserPassword
func stringAsUserPassword(v *string) UserPassword {
	return UserPassword{string: v}
}

// Unmarshal JSON data into one of the pointers in the struct
func (dst *UserPassword) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into HashedPassword
	err = json.Unmarshal(data, &dst.HashedPassword)
	if err == nil {
		jsonHashedPassword, _ := json.Marshal(dst.HashedPassword)
		if string(jsonHashedPassword) == "{}" { // empty struct
			dst.HashedPassword = nil
		} else {
			match++
		}
	} else {
		dst.HashedPassword = nil
	}

	// try to unmarshal data into string
	err = json.Unmarshal(data, &dst.string)
	if err == nil {
		jsonstring, _ := json.Marshal(dst.string)
		if string(jsonstring) == "{}" { // empty struct
			dst.string = nil
		} else {
			match++
		}
	} else {
		dst.string = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.HashedPassword = nil
		dst.string = nil

		return fmt.Errorf("Data matches more than one schema in oneOf(UserPassword)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("Data failed to match schemas in oneOf(UserPassword)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src UserPassword) MarshalJSON() ([]byte, error) {
	if src.HashedPassword != nil {
		return json.Marshal(&src.HashedPassword)
	}

	if src.string != nil {
		return json.Marshal(&src.string)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *UserPassword) GetActualInstance() interface{} {
	if obj.HashedPassword != nil {
		return obj.HashedPassword
	}

	if obj.string != nil {
		return obj.string
	}

	// all schemas are nil
	return nil
}

type NullableUserPassword struct {
	value *UserPassword
	isSet bool
}

func (v NullableUserPassword) Get() *UserPassword {
	return v.value
}

func (v *NullableUserPassword) Set(val *UserPassword) {
	v.value = val
	v.isSet = true
}

func (v NullableUserPassword) IsSet() bool {
	return v.isSet
}

func (v *NullableUserPassword) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserPassword(val *UserPassword) *NullableUserPassword {
	return &NullableUserPassword{value: val, isSet: true}
}

func (v NullableUserPassword) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserPassword) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
