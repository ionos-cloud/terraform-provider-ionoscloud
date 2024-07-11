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

// checks if the GetObjectLockConfigurationOutputObjectLockConfiguration type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &GetObjectLockConfigurationOutputObjectLockConfiguration{}

// GetObjectLockConfigurationOutputObjectLockConfiguration A container for an object lock configuration.
type GetObjectLockConfigurationOutputObjectLockConfiguration struct {
	XMLName           xml.Name        `xml:"GetObjectLockConfigurationOutputObjectLockConfiguration"`
	ObjectLockEnabled *string         `json:"ObjectLockEnabled,omitempty" xml:"ObjectLockEnabled"`
	Rule              *ObjectLockRule `json:"Rule,omitempty" xml:"Rule"`
}

// NewGetObjectLockConfigurationOutputObjectLockConfiguration instantiates a new GetObjectLockConfigurationOutputObjectLockConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGetObjectLockConfigurationOutputObjectLockConfiguration() *GetObjectLockConfigurationOutputObjectLockConfiguration {
	this := GetObjectLockConfigurationOutputObjectLockConfiguration{}

	return &this
}

// NewGetObjectLockConfigurationOutputObjectLockConfigurationWithDefaults instantiates a new GetObjectLockConfigurationOutputObjectLockConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGetObjectLockConfigurationOutputObjectLockConfigurationWithDefaults() *GetObjectLockConfigurationOutputObjectLockConfiguration {
	this := GetObjectLockConfigurationOutputObjectLockConfiguration{}
	return &this
}

// GetObjectLockEnabled returns the ObjectLockEnabled field value if set, zero value otherwise.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) GetObjectLockEnabled() string {
	if o == nil || IsNil(o.ObjectLockEnabled) {
		var ret string
		return ret
	}
	return *o.ObjectLockEnabled
}

// GetObjectLockEnabledOk returns a tuple with the ObjectLockEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) GetObjectLockEnabledOk() (*string, bool) {
	if o == nil || IsNil(o.ObjectLockEnabled) {
		return nil, false
	}
	return o.ObjectLockEnabled, true
}

// HasObjectLockEnabled returns a boolean if a field has been set.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) HasObjectLockEnabled() bool {
	if o != nil && !IsNil(o.ObjectLockEnabled) {
		return true
	}

	return false
}

// SetObjectLockEnabled gets a reference to the given string and assigns it to the ObjectLockEnabled field.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) SetObjectLockEnabled(v string) {
	o.ObjectLockEnabled = &v
}

// GetRule returns the Rule field value if set, zero value otherwise.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) GetRule() ObjectLockRule {
	if o == nil || IsNil(o.Rule) {
		var ret ObjectLockRule
		return ret
	}
	return *o.Rule
}

// GetRuleOk returns a tuple with the Rule field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) GetRuleOk() (*ObjectLockRule, bool) {
	if o == nil || IsNil(o.Rule) {
		return nil, false
	}
	return o.Rule, true
}

// HasRule returns a boolean if a field has been set.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) HasRule() bool {
	if o != nil && !IsNil(o.Rule) {
		return true
	}

	return false
}

// SetRule gets a reference to the given ObjectLockRule and assigns it to the Rule field.
func (o *GetObjectLockConfigurationOutputObjectLockConfiguration) SetRule(v ObjectLockRule) {
	o.Rule = &v
}

func (o GetObjectLockConfigurationOutputObjectLockConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o GetObjectLockConfigurationOutputObjectLockConfiguration) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ObjectLockEnabled) {
		toSerialize["ObjectLockEnabled"] = o.ObjectLockEnabled
	}
	if !IsNil(o.Rule) {
		toSerialize["Rule"] = o.Rule
	}
	return toSerialize, nil
}

type NullableGetObjectLockConfigurationOutputObjectLockConfiguration struct {
	value *GetObjectLockConfigurationOutputObjectLockConfiguration
	isSet bool
}

func (v NullableGetObjectLockConfigurationOutputObjectLockConfiguration) Get() *GetObjectLockConfigurationOutputObjectLockConfiguration {
	return v.value
}

func (v *NullableGetObjectLockConfigurationOutputObjectLockConfiguration) Set(val *GetObjectLockConfigurationOutputObjectLockConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableGetObjectLockConfigurationOutputObjectLockConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableGetObjectLockConfigurationOutputObjectLockConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGetObjectLockConfigurationOutputObjectLockConfiguration(val *GetObjectLockConfigurationOutputObjectLockConfiguration) *NullableGetObjectLockConfigurationOutputObjectLockConfiguration {
	return &NullableGetObjectLockConfigurationOutputObjectLockConfiguration{value: val, isSet: true}
}

func (v NullableGetObjectLockConfigurationOutputObjectLockConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGetObjectLockConfigurationOutputObjectLockConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
