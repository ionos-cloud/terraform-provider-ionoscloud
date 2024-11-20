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

// Rule Specifies lifecycle rules for an IONOS Object Storage bucket.
type Rule struct {
	XMLName xml.Name `xml:"Rule"`
	// Unique identifier for the rule. The value can't be longer than 255 characters.
	ID *string `json:"ID,omitempty" xml:"ID"`
	// Object key prefix that identifies one or more objects to which this rule applies. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.
	Prefix                         *string                         `json:"Prefix" xml:"Prefix"`
	Status                         *ExpirationStatus               `json:"Status" xml:"Status"`
	Expiration                     *LifecycleExpiration            `json:"Expiration,omitempty" xml:"Expiration"`
	NoncurrentVersionExpiration    *NoncurrentVersionExpiration    `json:"NoncurrentVersionExpiration,omitempty" xml:"NoncurrentVersionExpiration"`
	AbortIncompleteMultipartUpload *AbortIncompleteMultipartUpload `json:"AbortIncompleteMultipartUpload,omitempty" xml:"AbortIncompleteMultipartUpload"`
}

// NewRule instantiates a new Rule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRule(prefix string, status ExpirationStatus) *Rule {
	this := Rule{}

	this.Prefix = &prefix
	this.Status = &status

	return &this
}

// NewRuleWithDefaults instantiates a new Rule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRuleWithDefaults() *Rule {
	this := Rule{}
	return &this
}

// GetID returns the ID field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Rule) GetID() *string {
	if o == nil {
		return nil
	}

	return o.ID

}

// GetIDOk returns a tuple with the ID field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetIDOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ID, true
}

// SetID sets field value
func (o *Rule) SetID(v string) {

	o.ID = &v

}

// HasID returns a boolean if a field has been set.
func (o *Rule) HasID() bool {
	if o != nil && o.ID != nil {
		return true
	}

	return false
}

// GetPrefix returns the Prefix field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Rule) GetPrefix() *string {
	if o == nil {
		return nil
	}

	return o.Prefix

}

// GetPrefixOk returns a tuple with the Prefix field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetPrefixOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Prefix, true
}

// SetPrefix sets field value
func (o *Rule) SetPrefix(v string) {

	o.Prefix = &v

}

// HasPrefix returns a boolean if a field has been set.
func (o *Rule) HasPrefix() bool {
	if o != nil && o.Prefix != nil {
		return true
	}

	return false
}

// GetStatus returns the Status field value
// If the value is explicit nil, the zero value for ExpirationStatus will be returned
func (o *Rule) GetStatus() *ExpirationStatus {
	if o == nil {
		return nil
	}

	return o.Status

}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetStatusOk() (*ExpirationStatus, bool) {
	if o == nil {
		return nil, false
	}

	return o.Status, true
}

// SetStatus sets field value
func (o *Rule) SetStatus(v ExpirationStatus) {

	o.Status = &v

}

// HasStatus returns a boolean if a field has been set.
func (o *Rule) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

// GetExpiration returns the Expiration field value
// If the value is explicit nil, the zero value for LifecycleExpiration will be returned
func (o *Rule) GetExpiration() *LifecycleExpiration {
	if o == nil {
		return nil
	}

	return o.Expiration

}

// GetExpirationOk returns a tuple with the Expiration field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetExpirationOk() (*LifecycleExpiration, bool) {
	if o == nil {
		return nil, false
	}

	return o.Expiration, true
}

// SetExpiration sets field value
func (o *Rule) SetExpiration(v LifecycleExpiration) {

	o.Expiration = &v

}

// HasExpiration returns a boolean if a field has been set.
func (o *Rule) HasExpiration() bool {
	if o != nil && o.Expiration != nil {
		return true
	}

	return false
}

// GetNoncurrentVersionExpiration returns the NoncurrentVersionExpiration field value
// If the value is explicit nil, the zero value for NoncurrentVersionExpiration will be returned
func (o *Rule) GetNoncurrentVersionExpiration() *NoncurrentVersionExpiration {
	if o == nil {
		return nil
	}

	return o.NoncurrentVersionExpiration

}

// GetNoncurrentVersionExpirationOk returns a tuple with the NoncurrentVersionExpiration field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetNoncurrentVersionExpirationOk() (*NoncurrentVersionExpiration, bool) {
	if o == nil {
		return nil, false
	}

	return o.NoncurrentVersionExpiration, true
}

// SetNoncurrentVersionExpiration sets field value
func (o *Rule) SetNoncurrentVersionExpiration(v NoncurrentVersionExpiration) {

	o.NoncurrentVersionExpiration = &v

}

// HasNoncurrentVersionExpiration returns a boolean if a field has been set.
func (o *Rule) HasNoncurrentVersionExpiration() bool {
	if o != nil && o.NoncurrentVersionExpiration != nil {
		return true
	}

	return false
}

// GetAbortIncompleteMultipartUpload returns the AbortIncompleteMultipartUpload field value
// If the value is explicit nil, the zero value for AbortIncompleteMultipartUpload will be returned
func (o *Rule) GetAbortIncompleteMultipartUpload() *AbortIncompleteMultipartUpload {
	if o == nil {
		return nil
	}

	return o.AbortIncompleteMultipartUpload

}

// GetAbortIncompleteMultipartUploadOk returns a tuple with the AbortIncompleteMultipartUpload field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Rule) GetAbortIncompleteMultipartUploadOk() (*AbortIncompleteMultipartUpload, bool) {
	if o == nil {
		return nil, false
	}

	return o.AbortIncompleteMultipartUpload, true
}

// SetAbortIncompleteMultipartUpload sets field value
func (o *Rule) SetAbortIncompleteMultipartUpload(v AbortIncompleteMultipartUpload) {

	o.AbortIncompleteMultipartUpload = &v

}

// HasAbortIncompleteMultipartUpload returns a boolean if a field has been set.
func (o *Rule) HasAbortIncompleteMultipartUpload() bool {
	if o != nil && o.AbortIncompleteMultipartUpload != nil {
		return true
	}

	return false
}

func (o Rule) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ID != nil {
		toSerialize["ID"] = o.ID
	}

	if o.Prefix != nil {
		toSerialize["Prefix"] = o.Prefix
	}

	if o.Status != nil {
		toSerialize["Status"] = o.Status
	}

	if o.Expiration != nil {
		toSerialize["Expiration"] = o.Expiration
	}

	if o.NoncurrentVersionExpiration != nil {
		toSerialize["NoncurrentVersionExpiration"] = o.NoncurrentVersionExpiration
	}

	if o.AbortIncompleteMultipartUpload != nil {
		toSerialize["AbortIncompleteMultipartUpload"] = o.AbortIncompleteMultipartUpload
	}

	return json.Marshal(toSerialize)
}

type NullableRule struct {
	value *Rule
	isSet bool
}

func (v NullableRule) Get() *Rule {
	return v.value
}

func (v *NullableRule) Set(val *Rule) {
	v.value = val
	v.isSet = true
}

func (v NullableRule) IsSet() bool {
	return v.isSet
}

func (v *NullableRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRule(val *Rule) *NullableRule {
	return &NullableRule{value: val, isSet: true}
}

func (v NullableRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
