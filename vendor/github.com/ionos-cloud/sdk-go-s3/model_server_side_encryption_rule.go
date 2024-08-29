/*
 * IONOS S3 Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS S3 Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS S3 Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
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

// ServerSideEncryptionRule Specifies the default server-side encryption configuration.
type ServerSideEncryptionRule struct {
	XMLName                            xml.Name                       `xml:"Rule"`
	ApplyServerSideEncryptionByDefault *ServerSideEncryptionByDefault `json:"ApplyServerSideEncryptionByDefault,omitempty" xml:"ApplyServerSideEncryptionByDefault"`
}

// NewServerSideEncryptionRule instantiates a new ServerSideEncryptionRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewServerSideEncryptionRule() *ServerSideEncryptionRule {
	this := ServerSideEncryptionRule{}

	return &this
}

// NewServerSideEncryptionRuleWithDefaults instantiates a new ServerSideEncryptionRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewServerSideEncryptionRuleWithDefaults() *ServerSideEncryptionRule {
	this := ServerSideEncryptionRule{}
	return &this
}

// GetApplyServerSideEncryptionByDefault returns the ApplyServerSideEncryptionByDefault field value
// If the value is explicit nil, the zero value for ServerSideEncryptionByDefault will be returned
func (o *ServerSideEncryptionRule) GetApplyServerSideEncryptionByDefault() *ServerSideEncryptionByDefault {
	if o == nil {
		return nil
	}

	return o.ApplyServerSideEncryptionByDefault

}

// GetApplyServerSideEncryptionByDefaultOk returns a tuple with the ApplyServerSideEncryptionByDefault field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ServerSideEncryptionRule) GetApplyServerSideEncryptionByDefaultOk() (*ServerSideEncryptionByDefault, bool) {
	if o == nil {
		return nil, false
	}

	return o.ApplyServerSideEncryptionByDefault, true
}

// SetApplyServerSideEncryptionByDefault sets field value
func (o *ServerSideEncryptionRule) SetApplyServerSideEncryptionByDefault(v ServerSideEncryptionByDefault) {

	o.ApplyServerSideEncryptionByDefault = &v

}

// HasApplyServerSideEncryptionByDefault returns a boolean if a field has been set.
func (o *ServerSideEncryptionRule) HasApplyServerSideEncryptionByDefault() bool {
	if o != nil && o.ApplyServerSideEncryptionByDefault != nil {
		return true
	}

	return false
}

func (o ServerSideEncryptionRule) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ApplyServerSideEncryptionByDefault != nil {
		toSerialize["ApplyServerSideEncryptionByDefault"] = o.ApplyServerSideEncryptionByDefault
	}

	return json.Marshal(toSerialize)
}

type NullableServerSideEncryptionRule struct {
	value *ServerSideEncryptionRule
	isSet bool
}

func (v NullableServerSideEncryptionRule) Get() *ServerSideEncryptionRule {
	return v.value
}

func (v *NullableServerSideEncryptionRule) Set(val *ServerSideEncryptionRule) {
	v.value = val
	v.isSet = true
}

func (v NullableServerSideEncryptionRule) IsSet() bool {
	return v.isSet
}

func (v *NullableServerSideEncryptionRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableServerSideEncryptionRule(val *ServerSideEncryptionRule) *NullableServerSideEncryptionRule {
	return &NullableServerSideEncryptionRule{value: val, isSet: true}
}

func (v NullableServerSideEncryptionRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableServerSideEncryptionRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}