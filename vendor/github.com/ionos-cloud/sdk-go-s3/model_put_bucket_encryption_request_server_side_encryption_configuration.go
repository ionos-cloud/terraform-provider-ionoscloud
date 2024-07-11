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

// checks if the PutBucketEncryptionRequestServerSideEncryptionConfiguration type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PutBucketEncryptionRequestServerSideEncryptionConfiguration{}

// PutBucketEncryptionRequestServerSideEncryptionConfiguration Specifies the default server-side-encryption configuration. The valid value is AES256.
type PutBucketEncryptionRequestServerSideEncryptionConfiguration struct {
	XMLName xml.Name                   `xml:"PutBucketEncryptionRequestServerSideEncryptionConfiguration"`
	Rules   []ServerSideEncryptionRule `json:"Rules,omitempty" xml:"Rules"`
}

// NewPutBucketEncryptionRequestServerSideEncryptionConfiguration instantiates a new PutBucketEncryptionRequestServerSideEncryptionConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPutBucketEncryptionRequestServerSideEncryptionConfiguration() *PutBucketEncryptionRequestServerSideEncryptionConfiguration {
	this := PutBucketEncryptionRequestServerSideEncryptionConfiguration{}

	return &this
}

// NewPutBucketEncryptionRequestServerSideEncryptionConfigurationWithDefaults instantiates a new PutBucketEncryptionRequestServerSideEncryptionConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPutBucketEncryptionRequestServerSideEncryptionConfigurationWithDefaults() *PutBucketEncryptionRequestServerSideEncryptionConfiguration {
	this := PutBucketEncryptionRequestServerSideEncryptionConfiguration{}
	return &this
}

// GetRules returns the Rules field value if set, zero value otherwise.
func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) GetRules() []ServerSideEncryptionRule {
	if o == nil || IsNil(o.Rules) {
		var ret []ServerSideEncryptionRule
		return ret
	}
	return o.Rules
}

// GetRulesOk returns a tuple with the Rules field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) GetRulesOk() ([]ServerSideEncryptionRule, bool) {
	if o == nil || IsNil(o.Rules) {
		return nil, false
	}
	return o.Rules, true
}

// HasRules returns a boolean if a field has been set.
func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) HasRules() bool {
	if o != nil && !IsNil(o.Rules) {
		return true
	}

	return false
}

// SetRules gets a reference to the given []ServerSideEncryptionRule and assigns it to the Rules field.
func (o *PutBucketEncryptionRequestServerSideEncryptionConfiguration) SetRules(v []ServerSideEncryptionRule) {
	o.Rules = v
}

func (o PutBucketEncryptionRequestServerSideEncryptionConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PutBucketEncryptionRequestServerSideEncryptionConfiguration) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Rules) {
		toSerialize["Rules"] = o.Rules
	}
	return toSerialize, nil
}

type NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration struct {
	value *PutBucketEncryptionRequestServerSideEncryptionConfiguration
	isSet bool
}

func (v NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) Get() *PutBucketEncryptionRequestServerSideEncryptionConfiguration {
	return v.value
}

func (v *NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) Set(val *PutBucketEncryptionRequestServerSideEncryptionConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePutBucketEncryptionRequestServerSideEncryptionConfiguration(val *PutBucketEncryptionRequestServerSideEncryptionConfiguration) *NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration {
	return &NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration{value: val, isSet: true}
}

func (v NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePutBucketEncryptionRequestServerSideEncryptionConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
