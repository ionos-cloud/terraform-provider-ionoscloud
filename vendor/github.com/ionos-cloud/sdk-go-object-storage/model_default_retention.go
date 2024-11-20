/*
 * IONOS Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

import "encoding/xml"

// DefaultRetention The default Object Lock retention mode and period for new objects placed in the specified bucket. Bucket settings require both a mode and a period. The period can be either `Days` or `Years` but you must select one. You cannot specify `Days` and `Years` at the same time.
type DefaultRetention struct {
	XMLName xml.Name `xml:"DefaultRetention"`
	// The default Object Lock retention mode for new objects placed in the specified bucket. Must be used with either `Days` or `Years`.
	Mode *string `json:"Mode,omitempty" xml:"Mode"`
	// The number of days that you want to specify for the default retention period. Must be used with `Mode`.
	Days *int32 `json:"Days,omitempty" xml:"Days"`
	// The number of years that you want to specify for the default retention period. Must be used with `Mode`.
	Years *int32 `json:"Years,omitempty" xml:"Years"`
}

// NewDefaultRetention instantiates a new DefaultRetention object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDefaultRetention() *DefaultRetention {
	this := DefaultRetention{}

	return &this
}

// NewDefaultRetentionWithDefaults instantiates a new DefaultRetention object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDefaultRetentionWithDefaults() *DefaultRetention {
	this := DefaultRetention{}
	return &this
}

// GetMode returns the Mode field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DefaultRetention) GetMode() *string {
	if o == nil {
		return nil
	}

	return o.Mode

}

// GetModeOk returns a tuple with the Mode field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DefaultRetention) GetModeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Mode, true
}

// SetMode sets field value
func (o *DefaultRetention) SetMode(v string) {

	o.Mode = &v

}

// HasMode returns a boolean if a field has been set.
func (o *DefaultRetention) HasMode() bool {
	if o != nil && o.Mode != nil {
		return true
	}

	return false
}

// GetDays returns the Days field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *DefaultRetention) GetDays() *int32 {
	if o == nil {
		return nil
	}

	return o.Days

}

// GetDaysOk returns a tuple with the Days field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DefaultRetention) GetDaysOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Days, true
}

// SetDays sets field value
func (o *DefaultRetention) SetDays(v int32) {

	o.Days = &v

}

// HasDays returns a boolean if a field has been set.
func (o *DefaultRetention) HasDays() bool {
	if o != nil && o.Days != nil {
		return true
	}

	return false
}

// GetYears returns the Years field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *DefaultRetention) GetYears() *int32 {
	if o == nil {
		return nil
	}

	return o.Years

}

// GetYearsOk returns a tuple with the Years field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DefaultRetention) GetYearsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Years, true
}

// SetYears sets field value
func (o *DefaultRetention) SetYears(v int32) {

	o.Years = &v

}

// HasYears returns a boolean if a field has been set.
func (o *DefaultRetention) HasYears() bool {
	if o != nil && o.Years != nil {
		return true
	}

	return false
}

func (o DefaultRetention) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Mode != nil {
		toSerialize["Mode"] = o.Mode
	}

	if o.Days != nil {
		toSerialize["Days"] = o.Days
	}

	if o.Years != nil {
		toSerialize["Years"] = o.Years
	}

	return json.Marshal(toSerialize)
}

type NullableDefaultRetention struct {
	value *DefaultRetention
	isSet bool
}

func (v NullableDefaultRetention) Get() *DefaultRetention {
	return v.value
}

func (v *NullableDefaultRetention) Set(val *DefaultRetention) {
	v.value = val
	v.isSet = true
}

func (v NullableDefaultRetention) IsSet() bool {
	return v.isSet
}

func (v *NullableDefaultRetention) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDefaultRetention(val *DefaultRetention) *NullableDefaultRetention {
	return &NullableDefaultRetention{value: val, isSet: true}
}

func (v NullableDefaultRetention) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDefaultRetention) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
