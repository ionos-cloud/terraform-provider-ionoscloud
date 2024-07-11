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

// checks if the GetObjectTaggingOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetObjectTaggingOutput{}

// GetObjectTaggingOutput struct for GetObjectTaggingOutput
type GetObjectTaggingOutput struct {
	XMLName xml.Name `xml:"GetObjectTaggingOutput"`
	// Contains the tag set.
	TagSet []Tag `json:"TagSet" xml:"TagSet"`
}

// NewGetObjectTaggingOutput instantiates a new GetObjectTaggingOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetObjectTaggingOutput(tagSet []Tag) *GetObjectTaggingOutput {
	this := GetObjectTaggingOutput{}

	this.TagSet = tagSet

	return &this
}

// NewGetObjectTaggingOutputWithDefaults instantiates a new GetObjectTaggingOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetObjectTaggingOutputWithDefaults() *GetObjectTaggingOutput {
	this := GetObjectTaggingOutput{}
	return &this
}

// GetTagSet returns the TagSet field value
func (o *GetObjectTaggingOutput) GetTagSet() []Tag {
	if o == nil {
		var ret []Tag
		return ret
	}

	return o.TagSet
}

// GetTagSetOk returns a tuple with the TagSet field value
// and a boolean to check if the value has been set.
func (o *GetObjectTaggingOutput) GetTagSetOk() ([]Tag, bool) {
	if o == nil {
		return nil, false
	}
	return o.TagSet, true
}

// SetTagSet sets field value
func (o *GetObjectTaggingOutput) SetTagSet(v []Tag) {
	o.TagSet = v
}

func (o GetObjectTaggingOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetObjectTaggingOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.TagSet) {
		toSerialize["TagSet"] = o.TagSet
	}
	return toSerialize, nil
}

type NullableGetObjectTaggingOutput struct {
	value *GetObjectTaggingOutput
	isSet bool
}

func (v NullableGetObjectTaggingOutput) Get() *GetObjectTaggingOutput {
	return v.value
}

func (v *NullableGetObjectTaggingOutput) Set(val *GetObjectTaggingOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableGetObjectTaggingOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableGetObjectTaggingOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetObjectTaggingOutput(val *GetObjectTaggingOutput) *NullableGetObjectTaggingOutput {
	return &NullableGetObjectTaggingOutput{value: val, isSet: true}
}

func (v NullableGetObjectTaggingOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetObjectTaggingOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
