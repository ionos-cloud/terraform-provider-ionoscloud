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

// checks if the BlockPublicAccessPayload type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BlockPublicAccessPayload{}

// BlockPublicAccessPayload struct for BlockPublicAccessPayload
type BlockPublicAccessPayload struct {
	XMLName xml.Name `xml:"PublicAccessBlockConfiguration"`
	// Indicates that access to the bucket via Access Control Lists (ACLs) that grant public access is blocked. In other words, ACLs that allow public access are not permitted.
	BlockPublicAcls *bool `json:"BlockPublicAcls,omitempty" xml:"BlockPublicAcls"`
	// Instructs the system to ignore any ACLs that grant public access. Even if ACLs are set to allow public access, they will be disregarded.
	IgnorePublicAcls *bool `json:"IgnorePublicAcls,omitempty" xml:"IgnorePublicAcls"`
	// Blocks public access to the bucket via bucket policies. Bucket policies that grant public access will not be allowed.
	BlockPublicPolicy *bool `json:"BlockPublicPolicy,omitempty" xml:"BlockPublicPolicy"`
	// Restricts access to buckets that have public policies. Buckets with policies that grant public access will have their access restricted.
	RestrictPublicBuckets *bool `json:"RestrictPublicBuckets,omitempty" xml:"RestrictPublicBuckets"`
}

// NewBlockPublicAccessPayload instantiates a new BlockPublicAccessPayload object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBlockPublicAccessPayload() *BlockPublicAccessPayload {
	this := BlockPublicAccessPayload{}

	var blockPublicAcls bool = false
	this.BlockPublicAcls = &blockPublicAcls
	var ignorePublicAcls bool = false
	this.IgnorePublicAcls = &ignorePublicAcls
	var blockPublicPolicy bool = false
	this.BlockPublicPolicy = &blockPublicPolicy
	var restrictPublicBuckets bool = false
	this.RestrictPublicBuckets = &restrictPublicBuckets

	return &this
}

// NewBlockPublicAccessPayloadWithDefaults instantiates a new BlockPublicAccessPayload object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBlockPublicAccessPayloadWithDefaults() *BlockPublicAccessPayload {
	this := BlockPublicAccessPayload{}
	var blockPublicAcls bool = false
	this.BlockPublicAcls = &blockPublicAcls
	var ignorePublicAcls bool = false
	this.IgnorePublicAcls = &ignorePublicAcls
	var blockPublicPolicy bool = false
	this.BlockPublicPolicy = &blockPublicPolicy
	var restrictPublicBuckets bool = false
	this.RestrictPublicBuckets = &restrictPublicBuckets
	return &this
}

// GetBlockPublicAcls returns the BlockPublicAcls field value if set, zero value otherwise.
func (o *BlockPublicAccessPayload) GetBlockPublicAcls() bool {
	if o == nil || IsNil(o.BlockPublicAcls) {
		var ret bool
		return ret
	}
	return *o.BlockPublicAcls
}

// GetBlockPublicAclsOk returns a tuple with the BlockPublicAcls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BlockPublicAccessPayload) GetBlockPublicAclsOk() (*bool, bool) {
	if o == nil || IsNil(o.BlockPublicAcls) {
		return nil, false
	}
	return o.BlockPublicAcls, true
}

// HasBlockPublicAcls returns a boolean if a field has been set.
func (o *BlockPublicAccessPayload) HasBlockPublicAcls() bool {
	if o != nil && !IsNil(o.BlockPublicAcls) {
		return true
	}

	return false
}

// SetBlockPublicAcls gets a reference to the given bool and assigns it to the BlockPublicAcls field.
func (o *BlockPublicAccessPayload) SetBlockPublicAcls(v bool) {
	o.BlockPublicAcls = &v
}

// GetIgnorePublicAcls returns the IgnorePublicAcls field value if set, zero value otherwise.
func (o *BlockPublicAccessPayload) GetIgnorePublicAcls() bool {
	if o == nil || IsNil(o.IgnorePublicAcls) {
		var ret bool
		return ret
	}
	return *o.IgnorePublicAcls
}

// GetIgnorePublicAclsOk returns a tuple with the IgnorePublicAcls field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BlockPublicAccessPayload) GetIgnorePublicAclsOk() (*bool, bool) {
	if o == nil || IsNil(o.IgnorePublicAcls) {
		return nil, false
	}
	return o.IgnorePublicAcls, true
}

// HasIgnorePublicAcls returns a boolean if a field has been set.
func (o *BlockPublicAccessPayload) HasIgnorePublicAcls() bool {
	if o != nil && !IsNil(o.IgnorePublicAcls) {
		return true
	}

	return false
}

// SetIgnorePublicAcls gets a reference to the given bool and assigns it to the IgnorePublicAcls field.
func (o *BlockPublicAccessPayload) SetIgnorePublicAcls(v bool) {
	o.IgnorePublicAcls = &v
}

// GetBlockPublicPolicy returns the BlockPublicPolicy field value if set, zero value otherwise.
func (o *BlockPublicAccessPayload) GetBlockPublicPolicy() bool {
	if o == nil || IsNil(o.BlockPublicPolicy) {
		var ret bool
		return ret
	}
	return *o.BlockPublicPolicy
}

// GetBlockPublicPolicyOk returns a tuple with the BlockPublicPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BlockPublicAccessPayload) GetBlockPublicPolicyOk() (*bool, bool) {
	if o == nil || IsNil(o.BlockPublicPolicy) {
		return nil, false
	}
	return o.BlockPublicPolicy, true
}

// HasBlockPublicPolicy returns a boolean if a field has been set.
func (o *BlockPublicAccessPayload) HasBlockPublicPolicy() bool {
	if o != nil && !IsNil(o.BlockPublicPolicy) {
		return true
	}

	return false
}

// SetBlockPublicPolicy gets a reference to the given bool and assigns it to the BlockPublicPolicy field.
func (o *BlockPublicAccessPayload) SetBlockPublicPolicy(v bool) {
	o.BlockPublicPolicy = &v
}

// GetRestrictPublicBuckets returns the RestrictPublicBuckets field value if set, zero value otherwise.
func (o *BlockPublicAccessPayload) GetRestrictPublicBuckets() bool {
	if o == nil || IsNil(o.RestrictPublicBuckets) {
		var ret bool
		return ret
	}
	return *o.RestrictPublicBuckets
}

// GetRestrictPublicBucketsOk returns a tuple with the RestrictPublicBuckets field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BlockPublicAccessPayload) GetRestrictPublicBucketsOk() (*bool, bool) {
	if o == nil || IsNil(o.RestrictPublicBuckets) {
		return nil, false
	}
	return o.RestrictPublicBuckets, true
}

// HasRestrictPublicBuckets returns a boolean if a field has been set.
func (o *BlockPublicAccessPayload) HasRestrictPublicBuckets() bool {
	if o != nil && !IsNil(o.RestrictPublicBuckets) {
		return true
	}

	return false
}

// SetRestrictPublicBuckets gets a reference to the given bool and assigns it to the RestrictPublicBuckets field.
func (o *BlockPublicAccessPayload) SetRestrictPublicBuckets(v bool) {
	o.RestrictPublicBuckets = &v
}

func (o BlockPublicAccessPayload) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BlockPublicAccessPayload) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.BlockPublicAcls) {
		toSerialize["BlockPublicAcls"] = o.BlockPublicAcls
	}
	if !IsNil(o.IgnorePublicAcls) {
		toSerialize["IgnorePublicAcls"] = o.IgnorePublicAcls
	}
	if !IsNil(o.BlockPublicPolicy) {
		toSerialize["BlockPublicPolicy"] = o.BlockPublicPolicy
	}
	if !IsNil(o.RestrictPublicBuckets) {
		toSerialize["RestrictPublicBuckets"] = o.RestrictPublicBuckets
	}
	return toSerialize, nil
}

type NullableBlockPublicAccessPayload struct {
	value *BlockPublicAccessPayload
	isSet bool
}

func (v NullableBlockPublicAccessPayload) Get() *BlockPublicAccessPayload {
	return v.value
}

func (v *NullableBlockPublicAccessPayload) Set(val *BlockPublicAccessPayload) {
	v.value = val
	v.isSet = true
}

func (v NullableBlockPublicAccessPayload) IsSet() bool {
	return v.isSet
}

func (v *NullableBlockPublicAccessPayload) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBlockPublicAccessPayload(val *BlockPublicAccessPayload) *NullableBlockPublicAccessPayload {
	return &NullableBlockPublicAccessPayload{value: val, isSet: true}
}

func (v NullableBlockPublicAccessPayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBlockPublicAccessPayload) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
