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

// checks if the BucketPolicyStatementConditionDateLessThanOneOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BucketPolicyStatementConditionDateLessThanOneOf{}

// BucketPolicyStatementConditionDateLessThanOneOf struct for BucketPolicyStatementConditionDateLessThanOneOf
type BucketPolicyStatementConditionDateLessThanOneOf struct {
	XMLName      xml.Name `xml:"BucketPolicyStatementConditionDateLessThanOneOf"`
	AwsEpochTime *int32   `json:"aws:EpochTime,omitempty" xml:"aws:EpochTime"`
}

// NewBucketPolicyStatementConditionDateLessThanOneOf instantiates a new BucketPolicyStatementConditionDateLessThanOneOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBucketPolicyStatementConditionDateLessThanOneOf() *BucketPolicyStatementConditionDateLessThanOneOf {
	this := BucketPolicyStatementConditionDateLessThanOneOf{}

	return &this
}

// NewBucketPolicyStatementConditionDateLessThanOneOfWithDefaults instantiates a new BucketPolicyStatementConditionDateLessThanOneOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBucketPolicyStatementConditionDateLessThanOneOfWithDefaults() *BucketPolicyStatementConditionDateLessThanOneOf {
	this := BucketPolicyStatementConditionDateLessThanOneOf{}
	return &this
}

// GetAwsEpochTime returns the AwsEpochTime field value if set, zero value otherwise.
func (o *BucketPolicyStatementConditionDateLessThanOneOf) GetAwsEpochTime() int32 {
	if o == nil || IsNil(o.AwsEpochTime) {
		var ret int32
		return ret
	}
	return *o.AwsEpochTime
}

// GetAwsEpochTimeOk returns a tuple with the AwsEpochTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyStatementConditionDateLessThanOneOf) GetAwsEpochTimeOk() (*int32, bool) {
	if o == nil || IsNil(o.AwsEpochTime) {
		return nil, false
	}
	return o.AwsEpochTime, true
}

// HasAwsEpochTime returns a boolean if a field has been set.
func (o *BucketPolicyStatementConditionDateLessThanOneOf) HasAwsEpochTime() bool {
	if o != nil && !IsNil(o.AwsEpochTime) {
		return true
	}

	return false
}

// SetAwsEpochTime gets a reference to the given int32 and assigns it to the AwsEpochTime field.
func (o *BucketPolicyStatementConditionDateLessThanOneOf) SetAwsEpochTime(v int32) {
	o.AwsEpochTime = &v
}

func (o BucketPolicyStatementConditionDateLessThanOneOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BucketPolicyStatementConditionDateLessThanOneOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AwsEpochTime) {
		toSerialize["aws:EpochTime"] = o.AwsEpochTime
	}
	return toSerialize, nil
}

type NullableBucketPolicyStatementConditionDateLessThanOneOf struct {
	value *BucketPolicyStatementConditionDateLessThanOneOf
	isSet bool
}

func (v NullableBucketPolicyStatementConditionDateLessThanOneOf) Get() *BucketPolicyStatementConditionDateLessThanOneOf {
	return v.value
}

func (v *NullableBucketPolicyStatementConditionDateLessThanOneOf) Set(val *BucketPolicyStatementConditionDateLessThanOneOf) {
	v.value = val
	v.isSet = true
}

func (v NullableBucketPolicyStatementConditionDateLessThanOneOf) IsSet() bool {
	return v.isSet
}

func (v *NullableBucketPolicyStatementConditionDateLessThanOneOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBucketPolicyStatementConditionDateLessThanOneOf(val *BucketPolicyStatementConditionDateLessThanOneOf) *NullableBucketPolicyStatementConditionDateLessThanOneOf {
	return &NullableBucketPolicyStatementConditionDateLessThanOneOf{value: val, isSet: true}
}

func (v NullableBucketPolicyStatementConditionDateLessThanOneOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBucketPolicyStatementConditionDateLessThanOneOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}