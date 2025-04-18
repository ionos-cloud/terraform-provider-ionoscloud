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
	"time"
)

import "encoding/xml"

// checks if the Bucket type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Bucket{}

// Bucket A bucket in object storage is a flat container used to store an unlimited number of objects (files).
type Bucket struct {
	XMLName xml.Name `xml:"Bucket"`
	// The bucket name.
	Name *string `json:"Name,omitempty" xml:"Name"`
	// Represents the UTC date and time of bucket creation.
	CreationDate *IonosTime `json:"CreationDate,omitempty" xml:"CreationDate"`
}

// NewBucket instantiates a new Bucket object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBucket() *Bucket {
	this := Bucket{}

	return &this
}

// NewBucketWithDefaults instantiates a new Bucket object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBucketWithDefaults() *Bucket {
	this := Bucket{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Bucket) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Bucket) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Bucket) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Bucket) SetName(v string) {
	o.Name = &v
}

// GetCreationDate returns the CreationDate field value if set, zero value otherwise.
func (o *Bucket) GetCreationDate() time.Time {
	if o == nil || IsNil(o.CreationDate) {
		var ret time.Time
		return ret
	}
	return o.CreationDate.Time
}

// GetCreationDateOk returns a tuple with the CreationDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Bucket) GetCreationDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreationDate) {
		return nil, false
	}
	return &o.CreationDate.Time, true
}

// HasCreationDate returns a boolean if a field has been set.
func (o *Bucket) HasCreationDate() bool {
	if o != nil && !IsNil(o.CreationDate) {
		return true
	}

	return false
}

// SetCreationDate gets a reference to the given time.Time and assigns it to the CreationDate field.
func (o *Bucket) SetCreationDate(v time.Time) {
	o.CreationDate = &IonosTime{v}
}

func (o Bucket) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Bucket) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["Name"] = o.Name
	}
	if !IsNil(o.CreationDate) {
		toSerialize["CreationDate"] = o.CreationDate
	}
	return toSerialize, nil
}

type NullableBucket struct {
	value *Bucket
	isSet bool
}

func (v NullableBucket) Get() *Bucket {
	return v.value
}

func (v *NullableBucket) Set(val *Bucket) {
	v.value = val
	v.isSet = true
}

func (v NullableBucket) IsSet() bool {
	return v.isSet
}

func (v *NullableBucket) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBucket(val *Bucket) *NullableBucket {
	return &NullableBucket{value: val, isSet: true}
}

func (v NullableBucket) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBucket) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
