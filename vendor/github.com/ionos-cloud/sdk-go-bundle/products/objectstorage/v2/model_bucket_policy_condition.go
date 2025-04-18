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

// checks if the BucketPolicyCondition type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BucketPolicyCondition{}

// BucketPolicyCondition Conditions for when a policy is in effect.   IONOS Object Storage supports only the following condition operators and keys. Bucket policy does not yet support string interpolation.  **Condition Operators** - `IpAddress` - `NotIpAddress` - `DateGreaterThan` - `DateLessThan`  **Condition Keys** - `aws:SourceIp` - `aws:CurrentTime` - `aws:EpochTime`  Only the following condition keys are supported for the `ListBucket` action: - `s3:prefix` - `s3:delimiter` - `s3:max-keys`
type BucketPolicyCondition struct {
	XMLName         xml.Name                        `xml:"BucketPolicyCondition"`
	IpAddress       *BucketPolicyConditionIpAddress `json:"IpAddress,omitempty" xml:"IpAddress"`
	NotIpAddress    *BucketPolicyConditionIpAddress `json:"NotIpAddress,omitempty" xml:"NotIpAddress"`
	DateGreaterThan *BucketPolicyConditionDate      `json:"DateGreaterThan,omitempty" xml:"DateGreaterThan"`
	DateLessThan    *BucketPolicyConditionDate      `json:"DateLessThan,omitempty" xml:"DateLessThan"`
}

// NewBucketPolicyCondition instantiates a new BucketPolicyCondition object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBucketPolicyCondition() *BucketPolicyCondition {
	this := BucketPolicyCondition{}

	return &this
}

// NewBucketPolicyConditionWithDefaults instantiates a new BucketPolicyCondition object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBucketPolicyConditionWithDefaults() *BucketPolicyCondition {
	this := BucketPolicyCondition{}
	return &this
}

// GetIpAddress returns the IpAddress field value if set, zero value otherwise.
func (o *BucketPolicyCondition) GetIpAddress() BucketPolicyConditionIpAddress {
	if o == nil || IsNil(o.IpAddress) {
		var ret BucketPolicyConditionIpAddress
		return ret
	}
	return *o.IpAddress
}

// GetIpAddressOk returns a tuple with the IpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyCondition) GetIpAddressOk() (*BucketPolicyConditionIpAddress, bool) {
	if o == nil || IsNil(o.IpAddress) {
		return nil, false
	}
	return o.IpAddress, true
}

// HasIpAddress returns a boolean if a field has been set.
func (o *BucketPolicyCondition) HasIpAddress() bool {
	if o != nil && !IsNil(o.IpAddress) {
		return true
	}

	return false
}

// SetIpAddress gets a reference to the given BucketPolicyConditionIpAddress and assigns it to the IpAddress field.
func (o *BucketPolicyCondition) SetIpAddress(v BucketPolicyConditionIpAddress) {
	o.IpAddress = &v
}

// GetNotIpAddress returns the NotIpAddress field value if set, zero value otherwise.
func (o *BucketPolicyCondition) GetNotIpAddress() BucketPolicyConditionIpAddress {
	if o == nil || IsNil(o.NotIpAddress) {
		var ret BucketPolicyConditionIpAddress
		return ret
	}
	return *o.NotIpAddress
}

// GetNotIpAddressOk returns a tuple with the NotIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyCondition) GetNotIpAddressOk() (*BucketPolicyConditionIpAddress, bool) {
	if o == nil || IsNil(o.NotIpAddress) {
		return nil, false
	}
	return o.NotIpAddress, true
}

// HasNotIpAddress returns a boolean if a field has been set.
func (o *BucketPolicyCondition) HasNotIpAddress() bool {
	if o != nil && !IsNil(o.NotIpAddress) {
		return true
	}

	return false
}

// SetNotIpAddress gets a reference to the given BucketPolicyConditionIpAddress and assigns it to the NotIpAddress field.
func (o *BucketPolicyCondition) SetNotIpAddress(v BucketPolicyConditionIpAddress) {
	o.NotIpAddress = &v
}

// GetDateGreaterThan returns the DateGreaterThan field value if set, zero value otherwise.
func (o *BucketPolicyCondition) GetDateGreaterThan() BucketPolicyConditionDate {
	if o == nil || IsNil(o.DateGreaterThan) {
		var ret BucketPolicyConditionDate
		return ret
	}
	return *o.DateGreaterThan
}

// GetDateGreaterThanOk returns a tuple with the DateGreaterThan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyCondition) GetDateGreaterThanOk() (*BucketPolicyConditionDate, bool) {
	if o == nil || IsNil(o.DateGreaterThan) {
		return nil, false
	}
	return o.DateGreaterThan, true
}

// HasDateGreaterThan returns a boolean if a field has been set.
func (o *BucketPolicyCondition) HasDateGreaterThan() bool {
	if o != nil && !IsNil(o.DateGreaterThan) {
		return true
	}

	return false
}

// SetDateGreaterThan gets a reference to the given BucketPolicyConditionDate and assigns it to the DateGreaterThan field.
func (o *BucketPolicyCondition) SetDateGreaterThan(v BucketPolicyConditionDate) {
	o.DateGreaterThan = &v
}

// GetDateLessThan returns the DateLessThan field value if set, zero value otherwise.
func (o *BucketPolicyCondition) GetDateLessThan() BucketPolicyConditionDate {
	if o == nil || IsNil(o.DateLessThan) {
		var ret BucketPolicyConditionDate
		return ret
	}
	return *o.DateLessThan
}

// GetDateLessThanOk returns a tuple with the DateLessThan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BucketPolicyCondition) GetDateLessThanOk() (*BucketPolicyConditionDate, bool) {
	if o == nil || IsNil(o.DateLessThan) {
		return nil, false
	}
	return o.DateLessThan, true
}

// HasDateLessThan returns a boolean if a field has been set.
func (o *BucketPolicyCondition) HasDateLessThan() bool {
	if o != nil && !IsNil(o.DateLessThan) {
		return true
	}

	return false
}

// SetDateLessThan gets a reference to the given BucketPolicyConditionDate and assigns it to the DateLessThan field.
func (o *BucketPolicyCondition) SetDateLessThan(v BucketPolicyConditionDate) {
	o.DateLessThan = &v
}

func (o BucketPolicyCondition) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BucketPolicyCondition) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.IpAddress) {
		toSerialize["IpAddress"] = o.IpAddress
	}
	if !IsNil(o.NotIpAddress) {
		toSerialize["NotIpAddress"] = o.NotIpAddress
	}
	if !IsNil(o.DateGreaterThan) {
		toSerialize["DateGreaterThan"] = o.DateGreaterThan
	}
	if !IsNil(o.DateLessThan) {
		toSerialize["DateLessThan"] = o.DateLessThan
	}
	return toSerialize, nil
}

type NullableBucketPolicyCondition struct {
	value *BucketPolicyCondition
	isSet bool
}

func (v NullableBucketPolicyCondition) Get() *BucketPolicyCondition {
	return v.value
}

func (v *NullableBucketPolicyCondition) Set(val *BucketPolicyCondition) {
	v.value = val
	v.isSet = true
}

func (v NullableBucketPolicyCondition) IsSet() bool {
	return v.isSet
}

func (v *NullableBucketPolicyCondition) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBucketPolicyCondition(val *BucketPolicyCondition) *NullableBucketPolicyCondition {
	return &NullableBucketPolicyCondition{value: val, isSet: true}
}

func (v NullableBucketPolicyCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBucketPolicyCondition) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
