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

// Destination A container for information about the replication destination.
type Destination struct {
	XMLName xml.Name `xml:"Destination"`
	// Use the same \"Bucket\" value formatting as in the S3 API specification, that is, `arn:aws:s3:::{Bucket}`.
	Bucket       *string       `json:"Bucket" xml:"Bucket"`
	StorageClass *StorageClass `json:"StorageClass,omitempty" xml:"StorageClass"`
}

// NewDestination instantiates a new Destination object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDestination(bucket string) *Destination {
	this := Destination{}

	this.Bucket = &bucket

	return &this
}

// NewDestinationWithDefaults instantiates a new Destination object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDestinationWithDefaults() *Destination {
	this := Destination{}
	return &this
}

// GetBucket returns the Bucket field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Destination) GetBucket() *string {
	if o == nil {
		return nil
	}

	return o.Bucket

}

// GetBucketOk returns a tuple with the Bucket field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Destination) GetBucketOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Bucket, true
}

// SetBucket sets field value
func (o *Destination) SetBucket(v string) {

	o.Bucket = &v

}

// HasBucket returns a boolean if a field has been set.
func (o *Destination) HasBucket() bool {
	if o != nil && o.Bucket != nil {
		return true
	}

	return false
}

// GetStorageClass returns the StorageClass field value
// If the value is explicit nil, the zero value for StorageClass will be returned
func (o *Destination) GetStorageClass() *StorageClass {
	if o == nil {
		return nil
	}

	return o.StorageClass

}

// GetStorageClassOk returns a tuple with the StorageClass field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Destination) GetStorageClassOk() (*StorageClass, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageClass, true
}

// SetStorageClass sets field value
func (o *Destination) SetStorageClass(v StorageClass) {

	o.StorageClass = &v

}

// HasStorageClass returns a boolean if a field has been set.
func (o *Destination) HasStorageClass() bool {
	if o != nil && o.StorageClass != nil {
		return true
	}

	return false
}

func (o Destination) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Bucket != nil {
		toSerialize["Bucket"] = o.Bucket
	}

	if o.StorageClass != nil {
		toSerialize["StorageClass"] = o.StorageClass
	}

	return json.Marshal(toSerialize)
}

type NullableDestination struct {
	value *Destination
	isSet bool
}

func (v NullableDestination) Get() *Destination {
	return v.value
}

func (v *NullableDestination) Set(val *Destination) {
	v.value = val
	v.isSet = true
}

func (v NullableDestination) IsSet() bool {
	return v.isSet
}

func (v *NullableDestination) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDestination(val *Destination) *NullableDestination {
	return &NullableDestination{value: val, isSet: true}
}

func (v NullableDestination) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDestination) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
