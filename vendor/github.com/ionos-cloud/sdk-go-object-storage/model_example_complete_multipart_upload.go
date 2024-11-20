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

// ExampleCompleteMultipartUpload The container for the completed multipart upload details.
type ExampleCompleteMultipartUpload struct {
	XMLName xml.Name `xml:"ExampleCompleteMultipartUpload"`
	// Array of CompletedPart data types.
	Parts *[]CompletedPart `json:"Parts,omitempty" xml:"Parts"`
}

// NewExampleCompleteMultipartUpload instantiates a new ExampleCompleteMultipartUpload object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExampleCompleteMultipartUpload() *ExampleCompleteMultipartUpload {
	this := ExampleCompleteMultipartUpload{}

	return &this
}

// NewExampleCompleteMultipartUploadWithDefaults instantiates a new ExampleCompleteMultipartUpload object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExampleCompleteMultipartUploadWithDefaults() *ExampleCompleteMultipartUpload {
	this := ExampleCompleteMultipartUpload{}
	return &this
}

// GetParts returns the Parts field value
// If the value is explicit nil, the zero value for []CompletedPart will be returned
func (o *ExampleCompleteMultipartUpload) GetParts() *[]CompletedPart {
	if o == nil {
		return nil
	}

	return o.Parts

}

// GetPartsOk returns a tuple with the Parts field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ExampleCompleteMultipartUpload) GetPartsOk() (*[]CompletedPart, bool) {
	if o == nil {
		return nil, false
	}

	return o.Parts, true
}

// SetParts sets field value
func (o *ExampleCompleteMultipartUpload) SetParts(v []CompletedPart) {

	o.Parts = &v

}

// HasParts returns a boolean if a field has been set.
func (o *ExampleCompleteMultipartUpload) HasParts() bool {
	if o != nil && o.Parts != nil {
		return true
	}

	return false
}

func (o ExampleCompleteMultipartUpload) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Parts != nil {
		toSerialize["Parts"] = o.Parts
	}

	return json.Marshal(toSerialize)
}

type NullableExampleCompleteMultipartUpload struct {
	value *ExampleCompleteMultipartUpload
	isSet bool
}

func (v NullableExampleCompleteMultipartUpload) Get() *ExampleCompleteMultipartUpload {
	return v.value
}

func (v *NullableExampleCompleteMultipartUpload) Set(val *ExampleCompleteMultipartUpload) {
	v.value = val
	v.isSet = true
}

func (v NullableExampleCompleteMultipartUpload) IsSet() bool {
	return v.isSet
}

func (v *NullableExampleCompleteMultipartUpload) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExampleCompleteMultipartUpload(val *ExampleCompleteMultipartUpload) *NullableExampleCompleteMultipartUpload {
	return &NullableExampleCompleteMultipartUpload{value: val, isSet: true}
}

func (v NullableExampleCompleteMultipartUpload) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExampleCompleteMultipartUpload) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
