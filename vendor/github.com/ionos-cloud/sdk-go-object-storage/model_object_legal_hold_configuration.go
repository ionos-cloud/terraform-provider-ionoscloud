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

// ObjectLegalHoldConfiguration A Legal Hold configuration for an object.
type ObjectLegalHoldConfiguration struct {
	XMLName xml.Name `xml:"LegalHold"`
	// Object Legal Hold status
	Status *string `json:"Status,omitempty" xml:"Status"`
}

// NewObjectLegalHoldConfiguration instantiates a new ObjectLegalHoldConfiguration object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewObjectLegalHoldConfiguration() *ObjectLegalHoldConfiguration {
	this := ObjectLegalHoldConfiguration{}

	return &this
}

// NewObjectLegalHoldConfigurationWithDefaults instantiates a new ObjectLegalHoldConfiguration object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewObjectLegalHoldConfigurationWithDefaults() *ObjectLegalHoldConfiguration {
	this := ObjectLegalHoldConfiguration{}
	return &this
}

// GetStatus returns the Status field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ObjectLegalHoldConfiguration) GetStatus() *string {
	if o == nil {
		return nil
	}

	return o.Status

}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ObjectLegalHoldConfiguration) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Status, true
}

// SetStatus sets field value
func (o *ObjectLegalHoldConfiguration) SetStatus(v string) {

	o.Status = &v

}

// HasStatus returns a boolean if a field has been set.
func (o *ObjectLegalHoldConfiguration) HasStatus() bool {
	if o != nil && o.Status != nil {
		return true
	}

	return false
}

func (o ObjectLegalHoldConfiguration) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Status != nil {
		toSerialize["Status"] = o.Status
	}

	return json.Marshal(toSerialize)
}

type NullableObjectLegalHoldConfiguration struct {
	value *ObjectLegalHoldConfiguration
	isSet bool
}

func (v NullableObjectLegalHoldConfiguration) Get() *ObjectLegalHoldConfiguration {
	return v.value
}

func (v *NullableObjectLegalHoldConfiguration) Set(val *ObjectLegalHoldConfiguration) {
	v.value = val
	v.isSet = true
}

func (v NullableObjectLegalHoldConfiguration) IsSet() bool {
	return v.isSet
}

func (v *NullableObjectLegalHoldConfiguration) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableObjectLegalHoldConfiguration(val *ObjectLegalHoldConfiguration) *NullableObjectLegalHoldConfiguration {
	return &NullableObjectLegalHoldConfiguration{value: val, isSet: true}
}

func (v NullableObjectLegalHoldConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableObjectLegalHoldConfiguration) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}