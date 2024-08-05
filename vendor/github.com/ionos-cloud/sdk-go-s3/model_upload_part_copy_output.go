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

// UploadPartCopyOutput struct for UploadPartCopyOutput
type UploadPartCopyOutput struct {
	XMLName        xml.Name        `xml:"UploadPartCopyOutput"`
	CopyPartResult *CopyPartResult `json:"CopyPartResult,omitempty" xml:"CopyPartResult"`
}

// NewUploadPartCopyOutput instantiates a new UploadPartCopyOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUploadPartCopyOutput() *UploadPartCopyOutput {
	this := UploadPartCopyOutput{}

	return &this
}

// NewUploadPartCopyOutputWithDefaults instantiates a new UploadPartCopyOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUploadPartCopyOutputWithDefaults() *UploadPartCopyOutput {
	this := UploadPartCopyOutput{}
	return &this
}

// GetCopyPartResult returns the CopyPartResult field value
// If the value is explicit nil, the zero value for CopyPartResult will be returned
func (o *UploadPartCopyOutput) GetCopyPartResult() *CopyPartResult {
	if o == nil {
		return nil
	}

	return o.CopyPartResult

}

// GetCopyPartResultOk returns a tuple with the CopyPartResult field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *UploadPartCopyOutput) GetCopyPartResultOk() (*CopyPartResult, bool) {
	if o == nil {
		return nil, false
	}

	return o.CopyPartResult, true
}

// SetCopyPartResult sets field value
func (o *UploadPartCopyOutput) SetCopyPartResult(v CopyPartResult) {

	o.CopyPartResult = &v

}

// HasCopyPartResult returns a boolean if a field has been set.
func (o *UploadPartCopyOutput) HasCopyPartResult() bool {
	if o != nil && o.CopyPartResult != nil {
		return true
	}

	return false
}

func (o UploadPartCopyOutput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.CopyPartResult != nil {
		toSerialize["CopyPartResult"] = o.CopyPartResult
	}

	return json.Marshal(toSerialize)
}

type NullableUploadPartCopyOutput struct {
	value *UploadPartCopyOutput
	isSet bool
}

func (v NullableUploadPartCopyOutput) Get() *UploadPartCopyOutput {
	return v.value
}

func (v *NullableUploadPartCopyOutput) Set(val *UploadPartCopyOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableUploadPartCopyOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableUploadPartCopyOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUploadPartCopyOutput(val *UploadPartCopyOutput) *NullableUploadPartCopyOutput {
	return &NullableUploadPartCopyOutput{value: val, isSet: true}
}

func (v NullableUploadPartCopyOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUploadPartCopyOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
