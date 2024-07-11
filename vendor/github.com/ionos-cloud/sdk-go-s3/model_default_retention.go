/*
 * IONOS S3 Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS S3 Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 Management API Reference](https://api.ionos.com/docs/s3-management/v1/) for managing Access Keys - S3 API Reference for contract-owned buckets - current document - [S3 API Reference for user-owned buckets](https://api.ionos.com/docs/s3-user-owned-buckets/v2/)  ### User documentation [IONOS S3 Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
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

// checks if the DefaultRetention type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DefaultRetention{}

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

// GetMode returns the Mode field value if set, zero value otherwise.
func (o *DefaultRetention) GetMode() string {
	if o == nil || IsNil(o.Mode) {
		var ret string
		return ret
	}
	return *o.Mode
}

// GetModeOk returns a tuple with the Mode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DefaultRetention) GetModeOk() (*string, bool) {
	if o == nil || IsNil(o.Mode) {
		return nil, false
	}
	return o.Mode, true
}

// HasMode returns a boolean if a field has been set.
func (o *DefaultRetention) HasMode() bool {
	if o != nil && !IsNil(o.Mode) {
		return true
	}

	return false
}

// SetMode gets a reference to the given string and assigns it to the Mode field.
func (o *DefaultRetention) SetMode(v string) {
	o.Mode = &v
}

// GetDays returns the Days field value if set, zero value otherwise.
func (o *DefaultRetention) GetDays() int32 {
	if o == nil || IsNil(o.Days) {
		var ret int32
		return ret
	}
	return *o.Days
}

// GetDaysOk returns a tuple with the Days field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DefaultRetention) GetDaysOk() (*int32, bool) {
	if o == nil || IsNil(o.Days) {
		return nil, false
	}
	return o.Days, true
}

// HasDays returns a boolean if a field has been set.
func (o *DefaultRetention) HasDays() bool {
	if o != nil && !IsNil(o.Days) {
		return true
	}

	return false
}

// SetDays gets a reference to the given int32 and assigns it to the Days field.
func (o *DefaultRetention) SetDays(v int32) {
	o.Days = &v
}

// GetYears returns the Years field value if set, zero value otherwise.
func (o *DefaultRetention) GetYears() int32 {
	if o == nil || IsNil(o.Years) {
		var ret int32
		return ret
	}
	return *o.Years
}

// GetYearsOk returns a tuple with the Years field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DefaultRetention) GetYearsOk() (*int32, bool) {
	if o == nil || IsNil(o.Years) {
		return nil, false
	}
	return o.Years, true
}

// HasYears returns a boolean if a field has been set.
func (o *DefaultRetention) HasYears() bool {
	if o != nil && !IsNil(o.Years) {
		return true
	}

	return false
}

// SetYears gets a reference to the given int32 and assigns it to the Years field.
func (o *DefaultRetention) SetYears(v int32) {
	o.Years = &v
}

func (o DefaultRetention) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DefaultRetention) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Mode) {
		toSerialize["Mode"] = o.Mode
	}
	if !IsNil(o.Days) {
		toSerialize["Days"] = o.Days
	}
	if !IsNil(o.Years) {
		toSerialize["Years"] = o.Years
	}
	return toSerialize, nil
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
