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

// ObjectLockRule The container element for an Object Lock rule.
type ObjectLockRule struct {
	XMLName          xml.Name          `xml:"Rule"`
	DefaultRetention *DefaultRetention `json:"DefaultRetention,omitempty" xml:"DefaultRetention"`
}

// NewObjectLockRule instantiates a new ObjectLockRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewObjectLockRule() *ObjectLockRule {
	this := ObjectLockRule{}

	return &this
}

// NewObjectLockRuleWithDefaults instantiates a new ObjectLockRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewObjectLockRuleWithDefaults() *ObjectLockRule {
	this := ObjectLockRule{}
	return &this
}

// GetDefaultRetention returns the DefaultRetention field value
// If the value is explicit nil, the zero value for DefaultRetention will be returned
func (o *ObjectLockRule) GetDefaultRetention() *DefaultRetention {
	if o == nil {
		return nil
	}

	return o.DefaultRetention

}

// GetDefaultRetentionOk returns a tuple with the DefaultRetention field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ObjectLockRule) GetDefaultRetentionOk() (*DefaultRetention, bool) {
	if o == nil {
		return nil, false
	}

	return o.DefaultRetention, true
}

// SetDefaultRetention sets field value
func (o *ObjectLockRule) SetDefaultRetention(v DefaultRetention) {

	o.DefaultRetention = &v

}

// HasDefaultRetention returns a boolean if a field has been set.
func (o *ObjectLockRule) HasDefaultRetention() bool {
	if o != nil && o.DefaultRetention != nil {
		return true
	}

	return false
}

func (o ObjectLockRule) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.DefaultRetention != nil {
		toSerialize["DefaultRetention"] = o.DefaultRetention
	}

	return json.Marshal(toSerialize)
}

type NullableObjectLockRule struct {
	value *ObjectLockRule
	isSet bool
}

func (v NullableObjectLockRule) Get() *ObjectLockRule {
	return v.value
}

func (v *NullableObjectLockRule) Set(val *ObjectLockRule) {
	v.value = val
	v.isSet = true
}

func (v NullableObjectLockRule) IsSet() bool {
	return v.isSet
}

func (v *NullableObjectLockRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableObjectLockRule(val *ObjectLockRule) *NullableObjectLockRule {
	return &NullableObjectLockRule{value: val, isSet: true}
}

func (v NullableObjectLockRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableObjectLockRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
