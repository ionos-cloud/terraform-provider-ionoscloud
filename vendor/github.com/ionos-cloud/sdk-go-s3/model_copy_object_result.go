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

// CopyObjectResult Container for all response elements.
type CopyObjectResult struct {
	XMLName xml.Name `xml:"CopyObjectResult"`
	// Entity tag that identifies the object's data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.
	ETag *string `json:"ETag,omitempty" xml:"ETag"`
	// Creation date of the object.
	LastModified *IonosTime `json:"LastModified,omitempty" xml:"LastModified"`
}

// NewCopyObjectResult instantiates a new CopyObjectResult object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCopyObjectResult() *CopyObjectResult {
	this := CopyObjectResult{}

	return &this
}

// NewCopyObjectResultWithDefaults instantiates a new CopyObjectResult object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCopyObjectResultWithDefaults() *CopyObjectResult {
	this := CopyObjectResult{}
	return &this
}

// GetETag returns the ETag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CopyObjectResult) GetETag() *string {
	if o == nil {
		return nil
	}

	return o.ETag

}

// GetETagOk returns a tuple with the ETag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CopyObjectResult) GetETagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ETag, true
}

// SetETag sets field value
func (o *CopyObjectResult) SetETag(v string) {

	o.ETag = &v

}

// HasETag returns a boolean if a field has been set.
func (o *CopyObjectResult) HasETag() bool {
	if o != nil && o.ETag != nil {
		return true
	}

	return false
}

// GetLastModified returns the LastModified field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *CopyObjectResult) GetLastModified() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModified == nil {
		return nil
	}
	return &o.LastModified.Time

}

// GetLastModifiedOk returns a tuple with the LastModified field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CopyObjectResult) GetLastModifiedOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModified == nil {
		return nil, false
	}
	return &o.LastModified.Time, true

}

// SetLastModified sets field value
func (o *CopyObjectResult) SetLastModified(v time.Time) {

	o.LastModified = &IonosTime{v}

}

// HasLastModified returns a boolean if a field has been set.
func (o *CopyObjectResult) HasLastModified() bool {
	if o != nil && o.LastModified != nil {
		return true
	}

	return false
}

func (o CopyObjectResult) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ETag != nil {
		toSerialize["ETag"] = o.ETag
	}

	if o.LastModified != nil {
		toSerialize["LastModified"] = o.LastModified
	}

	return json.Marshal(toSerialize)
}

type NullableCopyObjectResult struct {
	value *CopyObjectResult
	isSet bool
}

func (v NullableCopyObjectResult) Get() *CopyObjectResult {
	return v.value
}

func (v *NullableCopyObjectResult) Set(val *CopyObjectResult) {
	v.value = val
	v.isSet = true
}

func (v NullableCopyObjectResult) IsSet() bool {
	return v.isSet
}

func (v *NullableCopyObjectResult) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCopyObjectResult(val *CopyObjectResult) *NullableCopyObjectResult {
	return &NullableCopyObjectResult{value: val, isSet: true}
}

func (v NullableCopyObjectResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCopyObjectResult) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
