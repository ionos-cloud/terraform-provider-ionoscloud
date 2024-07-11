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

// checks if the Rule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Rule{}

// Rule Specifies lifecycle rules for an IONOS S3 Object Storage bucket.
type Rule struct {
	XMLName xml.Name `xml:"Rule"`
	// Unique identifier for the rule. The value can't be longer than 255 characters.
	ID *string `json:"ID,omitempty" xml:"ID"`
	// Object key prefix that identifies one or more objects to which this rule applies. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.
	Prefix                         string                          `json:"Prefix" xml:"Prefix"`
	Status                         ExpirationStatus                `json:"Status" xml:"Status"`
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

	this.Prefix = prefix
	this.Status = status

	return &this
}

// NewRuleWithDefaults instantiates a new Rule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRuleWithDefaults() *Rule {
	this := Rule{}
	return &this
}

// GetID returns the ID field value if set, zero value otherwise.
func (o *Rule) GetID() string {
	if o == nil || IsNil(o.ID) {
		var ret string
		return ret
	}
	return *o.ID
}

// GetIDOk returns a tuple with the ID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Rule) GetIDOk() (*string, bool) {
	if o == nil || IsNil(o.ID) {
		return nil, false
	}
	return o.ID, true
}

// HasID returns a boolean if a field has been set.
func (o *Rule) HasID() bool {
	if o != nil && !IsNil(o.ID) {
		return true
	}

	return false
}

// SetID gets a reference to the given string and assigns it to the ID field.
func (o *Rule) SetID(v string) {
	o.ID = &v
}

// GetPrefix returns the Prefix field value
func (o *Rule) GetPrefix() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value
// and a boolean to check if the value has been set.
func (o *Rule) GetPrefixOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Prefix, true
}

// SetPrefix sets field value
func (o *Rule) SetPrefix(v string) {
	o.Prefix = v
}

// GetStatus returns the Status field value
func (o *Rule) GetStatus() ExpirationStatus {
	if o == nil {
		var ret ExpirationStatus
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *Rule) GetStatusOk() (*ExpirationStatus, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *Rule) SetStatus(v ExpirationStatus) {
	o.Status = v
}

// GetExpiration returns the Expiration field value if set, zero value otherwise.
func (o *Rule) GetExpiration() LifecycleExpiration {
	if o == nil || IsNil(o.Expiration) {
		var ret LifecycleExpiration
		return ret
	}
	return *o.Expiration
}

// GetExpirationOk returns a tuple with the Expiration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Rule) GetExpirationOk() (*LifecycleExpiration, bool) {
	if o == nil || IsNil(o.Expiration) {
		return nil, false
	}
	return o.Expiration, true
}

// HasExpiration returns a boolean if a field has been set.
func (o *Rule) HasExpiration() bool {
	if o != nil && !IsNil(o.Expiration) {
		return true
	}

	return false
}

// SetExpiration gets a reference to the given LifecycleExpiration and assigns it to the Expiration field.
func (o *Rule) SetExpiration(v LifecycleExpiration) {
	o.Expiration = &v
}

// GetNoncurrentVersionExpiration returns the NoncurrentVersionExpiration field value if set, zero value otherwise.
func (o *Rule) GetNoncurrentVersionExpiration() NoncurrentVersionExpiration {
	if o == nil || IsNil(o.NoncurrentVersionExpiration) {
		var ret NoncurrentVersionExpiration
		return ret
	}
	return *o.NoncurrentVersionExpiration
}

// GetNoncurrentVersionExpirationOk returns a tuple with the NoncurrentVersionExpiration field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Rule) GetNoncurrentVersionExpirationOk() (*NoncurrentVersionExpiration, bool) {
	if o == nil || IsNil(o.NoncurrentVersionExpiration) {
		return nil, false
	}
	return o.NoncurrentVersionExpiration, true
}

// HasNoncurrentVersionExpiration returns a boolean if a field has been set.
func (o *Rule) HasNoncurrentVersionExpiration() bool {
	if o != nil && !IsNil(o.NoncurrentVersionExpiration) {
		return true
	}

	return false
}

// SetNoncurrentVersionExpiration gets a reference to the given NoncurrentVersionExpiration and assigns it to the NoncurrentVersionExpiration field.
func (o *Rule) SetNoncurrentVersionExpiration(v NoncurrentVersionExpiration) {
	o.NoncurrentVersionExpiration = &v
}

// GetAbortIncompleteMultipartUpload returns the AbortIncompleteMultipartUpload field value if set, zero value otherwise.
func (o *Rule) GetAbortIncompleteMultipartUpload() AbortIncompleteMultipartUpload {
	if o == nil || IsNil(o.AbortIncompleteMultipartUpload) {
		var ret AbortIncompleteMultipartUpload
		return ret
	}
	return *o.AbortIncompleteMultipartUpload
}

// GetAbortIncompleteMultipartUploadOk returns a tuple with the AbortIncompleteMultipartUpload field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Rule) GetAbortIncompleteMultipartUploadOk() (*AbortIncompleteMultipartUpload, bool) {
	if o == nil || IsNil(o.AbortIncompleteMultipartUpload) {
		return nil, false
	}
	return o.AbortIncompleteMultipartUpload, true
}

// HasAbortIncompleteMultipartUpload returns a boolean if a field has been set.
func (o *Rule) HasAbortIncompleteMultipartUpload() bool {
	if o != nil && !IsNil(o.AbortIncompleteMultipartUpload) {
		return true
	}

	return false
}

// SetAbortIncompleteMultipartUpload gets a reference to the given AbortIncompleteMultipartUpload and assigns it to the AbortIncompleteMultipartUpload field.
func (o *Rule) SetAbortIncompleteMultipartUpload(v AbortIncompleteMultipartUpload) {
	o.AbortIncompleteMultipartUpload = &v
}

func (o Rule) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Rule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ID) {
		toSerialize["ID"] = o.ID
	}
	if !IsZero(o.Prefix) {
		toSerialize["Prefix"] = o.Prefix
	}
	if !IsZero(o.Status) {
		toSerialize["Status"] = o.Status
	}
	if !IsNil(o.Expiration) {
		toSerialize["Expiration"] = o.Expiration
	}
	if !IsNil(o.NoncurrentVersionExpiration) {
		toSerialize["NoncurrentVersionExpiration"] = o.NoncurrentVersionExpiration
	}
	if !IsNil(o.AbortIncompleteMultipartUpload) {
		toSerialize["AbortIncompleteMultipartUpload"] = o.AbortIncompleteMultipartUpload
	}
	return toSerialize, nil
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
