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

// checks if the GetBucketCorsOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetBucketCorsOutput{}

// GetBucketCorsOutput struct for GetBucketCorsOutput
type GetBucketCorsOutput struct {
	XMLName xml.Name `xml:"CORSConfiguration"`
	// A set of origins and methods (cross-origin access that you want to allow). You can add up to 100 rules to the configuration.
	CORSRules []CORSRule `json:"CORSRules,omitempty" xml:"CORSRule"`
}

// NewGetBucketCorsOutput instantiates a new GetBucketCorsOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetBucketCorsOutput() *GetBucketCorsOutput {
	this := GetBucketCorsOutput{}

	return &this
}

// NewGetBucketCorsOutputWithDefaults instantiates a new GetBucketCorsOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetBucketCorsOutputWithDefaults() *GetBucketCorsOutput {
	this := GetBucketCorsOutput{}
	return &this
}

// GetCORSRules returns the CORSRules field value if set, zero value otherwise.
func (o *GetBucketCorsOutput) GetCORSRules() []CORSRule {
	if o == nil || IsNil(o.CORSRules) {
		var ret []CORSRule
		return ret
	}
	return o.CORSRules
}

// GetCORSRulesOk returns a tuple with the CORSRules field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetBucketCorsOutput) GetCORSRulesOk() ([]CORSRule, bool) {
	if o == nil || IsNil(o.CORSRules) {
		return nil, false
	}
	return o.CORSRules, true
}

// HasCORSRules returns a boolean if a field has been set.
func (o *GetBucketCorsOutput) HasCORSRules() bool {
	if o != nil && !IsNil(o.CORSRules) {
		return true
	}

	return false
}

// SetCORSRules gets a reference to the given []CORSRule and assigns it to the CORSRules field.
func (o *GetBucketCorsOutput) SetCORSRules(v []CORSRule) {
	o.CORSRules = v
}

func (o GetBucketCorsOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetBucketCorsOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CORSRules) {
		toSerialize["CORSRules"] = o.CORSRules
	}
	return toSerialize, nil
}

type NullableGetBucketCorsOutput struct {
	value *GetBucketCorsOutput
	isSet bool
}

func (v NullableGetBucketCorsOutput) Get() *GetBucketCorsOutput {
	return v.value
}

func (v *NullableGetBucketCorsOutput) Set(val *GetBucketCorsOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableGetBucketCorsOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableGetBucketCorsOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetBucketCorsOutput(val *GetBucketCorsOutput) *NullableGetBucketCorsOutput {
	return &NullableGetBucketCorsOutput{value: val, isSet: true}
}

func (v NullableGetBucketCorsOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetBucketCorsOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
