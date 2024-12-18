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

// NoncurrentVersionExpiration Specifies when noncurrent object versions expire. Upon expiration, IONOS Object Storage permanently deletes the noncurrent object versions. You set this lifecycle configuration operation on a bucket that has versioning enabled (or suspended) to request that IONOS Object Storage delete noncurrent object versions at a specific period in the object's lifetime.
type NoncurrentVersionExpiration struct {
	XMLName xml.Name `xml:"NoncurrentVersionExpiration"`
	// Specifies the number of days an object is noncurrent before IONOS Object Storage can perform the associated operation.
	NoncurrentDays *int32 `json:"NoncurrentDays,omitempty" xml:"NoncurrentDays"`
}

// NewNoncurrentVersionExpiration instantiates a new NoncurrentVersionExpiration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNoncurrentVersionExpiration() *NoncurrentVersionExpiration {
	this := NoncurrentVersionExpiration{}

	return &this
}

// NewNoncurrentVersionExpirationWithDefaults instantiates a new NoncurrentVersionExpiration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNoncurrentVersionExpirationWithDefaults() *NoncurrentVersionExpiration {
	this := NoncurrentVersionExpiration{}
	return &this
}

// GetNoncurrentDays returns the NoncurrentDays field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *NoncurrentVersionExpiration) GetNoncurrentDays() *int32 {
	if o == nil {
		return nil
	}

	return o.NoncurrentDays

}

// GetNoncurrentDaysOk returns a tuple with the NoncurrentDays field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoncurrentVersionExpiration) GetNoncurrentDaysOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.NoncurrentDays, true
}

// SetNoncurrentDays sets field value
func (o *NoncurrentVersionExpiration) SetNoncurrentDays(v int32) {

	o.NoncurrentDays = &v

}

// HasNoncurrentDays returns a boolean if a field has been set.
func (o *NoncurrentVersionExpiration) HasNoncurrentDays() bool {
	if o != nil && o.NoncurrentDays != nil {
		return true
	}

	return false
}

func (o NoncurrentVersionExpiration) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.NoncurrentDays != nil {
		toSerialize["NoncurrentDays"] = o.NoncurrentDays
	}

	return json.Marshal(toSerialize)
}

type NullableNoncurrentVersionExpiration struct {
	value *NoncurrentVersionExpiration
	isSet bool
}

func (v NullableNoncurrentVersionExpiration) Get() *NoncurrentVersionExpiration {
	return v.value
}

func (v *NullableNoncurrentVersionExpiration) Set(val *NoncurrentVersionExpiration) {
	v.value = val
	v.isSet = true
}

func (v NullableNoncurrentVersionExpiration) IsSet() bool {
	return v.isSet
}

func (v *NullableNoncurrentVersionExpiration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNoncurrentVersionExpiration(val *NoncurrentVersionExpiration) *NullableNoncurrentVersionExpiration {
	return &NullableNoncurrentVersionExpiration{value: val, isSet: true}
}

func (v NullableNoncurrentVersionExpiration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNoncurrentVersionExpiration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
