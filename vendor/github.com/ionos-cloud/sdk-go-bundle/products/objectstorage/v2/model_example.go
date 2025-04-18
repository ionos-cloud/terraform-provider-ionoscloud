/*
 * IONOS Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package objectstorage

import (
	"encoding/json"
)

import "encoding/xml"

// checks if the Example type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Example{}

// Example struct for Example
type Example struct {
	XMLName                 xml.Name                        `xml:"Example"`
	CompleteMultipartUpload *ExampleCompleteMultipartUpload `json:"CompleteMultipartUpload,omitempty" xml:"CompleteMultipartUpload"`
}

// NewExample instantiates a new Example object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExample() *Example {
	this := Example{}

	return &this
}

// NewExampleWithDefaults instantiates a new Example object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExampleWithDefaults() *Example {
	this := Example{}
	return &this
}

// GetCompleteMultipartUpload returns the CompleteMultipartUpload field value if set, zero value otherwise.
func (o *Example) GetCompleteMultipartUpload() ExampleCompleteMultipartUpload {
	if o == nil || IsNil(o.CompleteMultipartUpload) {
		var ret ExampleCompleteMultipartUpload
		return ret
	}
	return *o.CompleteMultipartUpload
}

// GetCompleteMultipartUploadOk returns a tuple with the CompleteMultipartUpload field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Example) GetCompleteMultipartUploadOk() (*ExampleCompleteMultipartUpload, bool) {
	if o == nil || IsNil(o.CompleteMultipartUpload) {
		return nil, false
	}
	return o.CompleteMultipartUpload, true
}

// HasCompleteMultipartUpload returns a boolean if a field has been set.
func (o *Example) HasCompleteMultipartUpload() bool {
	if o != nil && !IsNil(o.CompleteMultipartUpload) {
		return true
	}

	return false
}

// SetCompleteMultipartUpload gets a reference to the given ExampleCompleteMultipartUpload and assigns it to the CompleteMultipartUpload field.
func (o *Example) SetCompleteMultipartUpload(v ExampleCompleteMultipartUpload) {
	o.CompleteMultipartUpload = &v
}

func (o Example) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Example) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CompleteMultipartUpload) {
		toSerialize["CompleteMultipartUpload"] = o.CompleteMultipartUpload
	}
	return toSerialize, nil
}

type NullableExample struct {
	value *Example
	isSet bool
}

func (v NullableExample) Get() *Example {
	return v.value
}

func (v *NullableExample) Set(val *Example) {
	v.value = val
	v.isSet = true
}

func (v NullableExample) IsSet() bool {
	return v.isSet
}

func (v *NullableExample) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExample(val *Example) *NullableExample {
	return &NullableExample{value: val, isSet: true}
}

func (v NullableExample) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExample) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
