/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.   MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
)

// DayOfTheWeek The week day.
type DayOfTheWeek string

// List of DayOfTheWeek
const (
	DAYOFTHEWEEK_SUNDAY    DayOfTheWeek = "Sunday"
	DAYOFTHEWEEK_MONDAY    DayOfTheWeek = "Monday"
	DAYOFTHEWEEK_TUESDAY   DayOfTheWeek = "Tuesday"
	DAYOFTHEWEEK_WEDNESDAY DayOfTheWeek = "Wednesday"
	DAYOFTHEWEEK_THURSDAY  DayOfTheWeek = "Thursday"
	DAYOFTHEWEEK_FRIDAY    DayOfTheWeek = "Friday"
	DAYOFTHEWEEK_SATURDAY  DayOfTheWeek = "Saturday"
)

func (v *DayOfTheWeek) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DayOfTheWeek(value)
	for _, existing := range []DayOfTheWeek{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DayOfTheWeek", value)
}

// Ptr returns reference to DayOfTheWeek value
func (v DayOfTheWeek) Ptr() *DayOfTheWeek {
	return &v
}

type NullableDayOfTheWeek struct {
	value *DayOfTheWeek
	isSet bool
}

func (v NullableDayOfTheWeek) Get() *DayOfTheWeek {
	return v.value
}

func (v *NullableDayOfTheWeek) Set(val *DayOfTheWeek) {
	v.value = val
	v.isSet = true
}

func (v NullableDayOfTheWeek) IsSet() bool {
	return v.isSet
}

func (v *NullableDayOfTheWeek) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDayOfTheWeek(val *DayOfTheWeek) *NullableDayOfTheWeek {
	return &NullableDayOfTheWeek{value: val, isSet: true}
}

func (v NullableDayOfTheWeek) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDayOfTheWeek) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
