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

	"time"
)

import "encoding/xml"

// checks if the PutObjectRetentionRequestRetention type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PutObjectRetentionRequestRetention{}

// PutObjectRetentionRequestRetention A Retention configuration for an object.
type PutObjectRetentionRequestRetention struct {
	XMLName xml.Name `xml:"PutObjectRetentionRequestRetention"`
	// Indicates the Retention mode for the specified object.
	Mode *string `json:"Mode,omitempty" xml:"Mode"`
	// The date on which this Object Lock Retention will expire.
	RetainUntilDate *IonosTime `json:"RetainUntilDate,omitempty" xml:"RetainUntilDate"`
}

// NewPutObjectRetentionRequestRetention instantiates a new PutObjectRetentionRequestRetention object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPutObjectRetentionRequestRetention() *PutObjectRetentionRequestRetention {
	this := PutObjectRetentionRequestRetention{}

	return &this
}

// NewPutObjectRetentionRequestRetentionWithDefaults instantiates a new PutObjectRetentionRequestRetention object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPutObjectRetentionRequestRetentionWithDefaults() *PutObjectRetentionRequestRetention {
	this := PutObjectRetentionRequestRetention{}
	return &this
}

// GetMode returns the Mode field value if set, zero value otherwise.
func (o *PutObjectRetentionRequestRetention) GetMode() string {
	if o == nil || IsNil(o.Mode) {
		var ret string
		return ret
	}
	return *o.Mode
}

// GetModeOk returns a tuple with the Mode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutObjectRetentionRequestRetention) GetModeOk() (*string, bool) {
	if o == nil || IsNil(o.Mode) {
		return nil, false
	}
	return o.Mode, true
}

// HasMode returns a boolean if a field has been set.
func (o *PutObjectRetentionRequestRetention) HasMode() bool {
	if o != nil && !IsNil(o.Mode) {
		return true
	}

	return false
}

// SetMode gets a reference to the given string and assigns it to the Mode field.
func (o *PutObjectRetentionRequestRetention) SetMode(v string) {
	o.Mode = &v
}

// GetRetainUntilDate returns the RetainUntilDate field value if set, zero value otherwise.
func (o *PutObjectRetentionRequestRetention) GetRetainUntilDate() time.Time {
	if o == nil || IsNil(o.RetainUntilDate) {
		var ret time.Time
		return ret
	}
	return o.RetainUntilDate.Time
}

// GetRetainUntilDateOk returns a tuple with the RetainUntilDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutObjectRetentionRequestRetention) GetRetainUntilDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.RetainUntilDate) {
		return nil, false
	}
	return &o.RetainUntilDate.Time, true
}

// HasRetainUntilDate returns a boolean if a field has been set.
func (o *PutObjectRetentionRequestRetention) HasRetainUntilDate() bool {
	if o != nil && !IsNil(o.RetainUntilDate) {
		return true
	}

	return false
}

// SetRetainUntilDate gets a reference to the given time.Time and assigns it to the RetainUntilDate field.
func (o *PutObjectRetentionRequestRetention) SetRetainUntilDate(v time.Time) {
	o.RetainUntilDate = &IonosTime{v}
}

func (o PutObjectRetentionRequestRetention) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PutObjectRetentionRequestRetention) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Mode) {
		toSerialize["Mode"] = o.Mode
	}
	if !IsNil(o.RetainUntilDate) {
		toSerialize["RetainUntilDate"] = o.RetainUntilDate
	}
	return toSerialize, nil
}

type NullablePutObjectRetentionRequestRetention struct {
	value *PutObjectRetentionRequestRetention
	isSet bool
}

func (v NullablePutObjectRetentionRequestRetention) Get() *PutObjectRetentionRequestRetention {
	return v.value
}

func (v *NullablePutObjectRetentionRequestRetention) Set(val *PutObjectRetentionRequestRetention) {
	v.value = val
	v.isSet = true
}

func (v NullablePutObjectRetentionRequestRetention) IsSet() bool {
	return v.isSet
}

func (v *NullablePutObjectRetentionRequestRetention) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePutObjectRetentionRequestRetention(val *PutObjectRetentionRequestRetention) *NullablePutObjectRetentionRequestRetention {
	return &NullablePutObjectRetentionRequestRetention{value: val, isSet: true}
}

func (v NullablePutObjectRetentionRequestRetention) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePutObjectRetentionRequestRetention) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}