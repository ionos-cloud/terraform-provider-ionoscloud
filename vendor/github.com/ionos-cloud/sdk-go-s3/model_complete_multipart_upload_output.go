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

// checks if the CompleteMultipartUploadOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CompleteMultipartUploadOutput{}

// CompleteMultipartUploadOutput struct for CompleteMultipartUploadOutput
type CompleteMultipartUploadOutput struct {
	XMLName xml.Name `xml:"CompleteMultipartUploadOutput"`
	// The URI that identifies the newly created object.
	Location *string `json:"Location,omitempty" xml:"Location"`
	// The bucket name.
	Bucket *string `json:"Bucket,omitempty" xml:"Name"`
	// The object key.
	Key *string `json:"Key,omitempty" xml:"Key"`
	// Entity tag that identifies the object's data. Objects with different object data will have different entity tags. The entity tag is an opaque string. The entity tag may or may not be an MD5 digest of the object data. If the entity tag is not an MD5 digest of the object data, it will contain one or more nonhexadecimal characters and/or will consist of less than 32 or more than 32 hexadecimal digits.
	ETag *string `json:"ETag,omitempty" xml:"ETag"`
}

// NewCompleteMultipartUploadOutput instantiates a new CompleteMultipartUploadOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCompleteMultipartUploadOutput() *CompleteMultipartUploadOutput {
	this := CompleteMultipartUploadOutput{}

	return &this
}

// NewCompleteMultipartUploadOutputWithDefaults instantiates a new CompleteMultipartUploadOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCompleteMultipartUploadOutputWithDefaults() *CompleteMultipartUploadOutput {
	this := CompleteMultipartUploadOutput{}
	return &this
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *CompleteMultipartUploadOutput) GetLocation() string {
	if o == nil || IsNil(o.Location) {
		var ret string
		return ret
	}
	return *o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CompleteMultipartUploadOutput) GetLocationOk() (*string, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *CompleteMultipartUploadOutput) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given string and assigns it to the Location field.
func (o *CompleteMultipartUploadOutput) SetLocation(v string) {
	o.Location = &v
}

// GetBucket returns the Bucket field value if set, zero value otherwise.
func (o *CompleteMultipartUploadOutput) GetBucket() string {
	if o == nil || IsNil(o.Bucket) {
		var ret string
		return ret
	}
	return *o.Bucket
}

// GetBucketOk returns a tuple with the Bucket field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CompleteMultipartUploadOutput) GetBucketOk() (*string, bool) {
	if o == nil || IsNil(o.Bucket) {
		return nil, false
	}
	return o.Bucket, true
}

// HasBucket returns a boolean if a field has been set.
func (o *CompleteMultipartUploadOutput) HasBucket() bool {
	if o != nil && !IsNil(o.Bucket) {
		return true
	}

	return false
}

// SetBucket gets a reference to the given string and assigns it to the Bucket field.
func (o *CompleteMultipartUploadOutput) SetBucket(v string) {
	o.Bucket = &v
}

// GetKey returns the Key field value if set, zero value otherwise.
func (o *CompleteMultipartUploadOutput) GetKey() string {
	if o == nil || IsNil(o.Key) {
		var ret string
		return ret
	}
	return *o.Key
}

// GetKeyOk returns a tuple with the Key field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CompleteMultipartUploadOutput) GetKeyOk() (*string, bool) {
	if o == nil || IsNil(o.Key) {
		return nil, false
	}
	return o.Key, true
}

// HasKey returns a boolean if a field has been set.
func (o *CompleteMultipartUploadOutput) HasKey() bool {
	if o != nil && !IsNil(o.Key) {
		return true
	}

	return false
}

// SetKey gets a reference to the given string and assigns it to the Key field.
func (o *CompleteMultipartUploadOutput) SetKey(v string) {
	o.Key = &v
}

// GetETag returns the ETag field value if set, zero value otherwise.
func (o *CompleteMultipartUploadOutput) GetETag() string {
	if o == nil || IsNil(o.ETag) {
		var ret string
		return ret
	}
	return *o.ETag
}

// GetETagOk returns a tuple with the ETag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CompleteMultipartUploadOutput) GetETagOk() (*string, bool) {
	if o == nil || IsNil(o.ETag) {
		return nil, false
	}
	return o.ETag, true
}

// HasETag returns a boolean if a field has been set.
func (o *CompleteMultipartUploadOutput) HasETag() bool {
	if o != nil && !IsNil(o.ETag) {
		return true
	}

	return false
}

// SetETag gets a reference to the given string and assigns it to the ETag field.
func (o *CompleteMultipartUploadOutput) SetETag(v string) {
	o.ETag = &v
}

func (o CompleteMultipartUploadOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CompleteMultipartUploadOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Location) {
		toSerialize["Location"] = o.Location
	}
	if !IsNil(o.Bucket) {
		toSerialize["Bucket"] = o.Bucket
	}
	if !IsNil(o.Key) {
		toSerialize["Key"] = o.Key
	}
	if !IsNil(o.ETag) {
		toSerialize["ETag"] = o.ETag
	}
	return toSerialize, nil
}

type NullableCompleteMultipartUploadOutput struct {
	value *CompleteMultipartUploadOutput
	isSet bool
}

func (v NullableCompleteMultipartUploadOutput) Get() *CompleteMultipartUploadOutput {
	return v.value
}

func (v *NullableCompleteMultipartUploadOutput) Set(val *CompleteMultipartUploadOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableCompleteMultipartUploadOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableCompleteMultipartUploadOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCompleteMultipartUploadOutput(val *CompleteMultipartUploadOutput) *NullableCompleteMultipartUploadOutput {
	return &NullableCompleteMultipartUploadOutput{value: val, isSet: true}
}

func (v NullableCompleteMultipartUploadOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCompleteMultipartUploadOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
