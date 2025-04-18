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

// checks if the HeadObjectOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &HeadObjectOutput{}

// HeadObjectOutput struct for HeadObjectOutput
type HeadObjectOutput struct {
	XMLName xml.Name `xml:"HeadObjectOutput"`
	// A map of metadata to store with the object. Each key must start with  `x-amz-meta-` prefix.
	Metadata *map[string]string `json:"Metadata,omitempty" xml:"Metadata"`
}

// NewHeadObjectOutput instantiates a new HeadObjectOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewHeadObjectOutput() *HeadObjectOutput {
	this := HeadObjectOutput{}

	return &this
}

// NewHeadObjectOutputWithDefaults instantiates a new HeadObjectOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewHeadObjectOutputWithDefaults() *HeadObjectOutput {
	this := HeadObjectOutput{}
	return &this
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *HeadObjectOutput) GetMetadata() map[string]string {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]string
		return ret
	}
	return *o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *HeadObjectOutput) GetMetadataOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Metadata) {
		return nil, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *HeadObjectOutput) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]string and assigns it to the Metadata field.
func (o *HeadObjectOutput) SetMetadata(v map[string]string) {
	o.Metadata = &v
}

func (o HeadObjectOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o HeadObjectOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Metadata) {
		toSerialize["Metadata"] = o.Metadata
	}
	return toSerialize, nil
}

type NullableHeadObjectOutput struct {
	value *HeadObjectOutput
	isSet bool
}

func (v NullableHeadObjectOutput) Get() *HeadObjectOutput {
	return v.value
}

func (v *NullableHeadObjectOutput) Set(val *HeadObjectOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableHeadObjectOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableHeadObjectOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableHeadObjectOutput(val *HeadObjectOutput) *NullableHeadObjectOutput {
	return &NullableHeadObjectOutput{value: val, isSet: true}
}

func (v NullableHeadObjectOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableHeadObjectOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
