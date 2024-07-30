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

// DeletionError Container for all error elements.
type DeletionError struct {
	XMLName xml.Name `xml:"DeletionError"`
	// The object key.
	Key *string `json:"Key,omitempty" xml:"Key"`
	// The version ID of the object.
	VersionId *string `json:"VersionId,omitempty" xml:"VersionId"`
	Code      *string `json:"Code,omitempty" xml:"Code"`
	Message   *string `json:"Message,omitempty" xml:"Message"`
}

// NewDeletionError instantiates a new DeletionError object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDeletionError() *DeletionError {
	this := DeletionError{}

	return &this
}

// NewDeletionErrorWithDefaults instantiates a new DeletionError object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDeletionErrorWithDefaults() *DeletionError {
	this := DeletionError{}
	return &this
}

// GetKey returns the Key field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DeletionError) GetKey() *string {
	if o == nil {
		return nil
	}

	return o.Key

}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DeletionError) GetKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Key, true
}

// SetKey sets field value
func (o *DeletionError) SetKey(v string) {

	o.Key = &v

}

// HasKey returns a boolean if a field has been set.
func (o *DeletionError) HasKey() bool {
	if o != nil && o.Key != nil {
		return true
	}

	return false
}

// GetVersionId returns the VersionId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DeletionError) GetVersionId() *string {
	if o == nil {
		return nil
	}

	return o.VersionId

}

// GetVersionIdOk returns a tuple with the VersionId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DeletionError) GetVersionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.VersionId, true
}

// SetVersionId sets field value
func (o *DeletionError) SetVersionId(v string) {

	o.VersionId = &v

}

// HasVersionId returns a boolean if a field has been set.
func (o *DeletionError) HasVersionId() bool {
	if o != nil && o.VersionId != nil {
		return true
	}

	return false
}

// GetCode returns the Code field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DeletionError) GetCode() *string {
	if o == nil {
		return nil
	}

	return o.Code

}

// GetCodeOk returns a tuple with the Code field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DeletionError) GetCodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Code, true
}

// SetCode sets field value
func (o *DeletionError) SetCode(v string) {

	o.Code = &v

}

// HasCode returns a boolean if a field has been set.
func (o *DeletionError) HasCode() bool {
	if o != nil && o.Code != nil {
		return true
	}

	return false
}

// GetMessage returns the Message field value
// If the value is explicit nil, the zero value for string will be returned
func (o *DeletionError) GetMessage() *string {
	if o == nil {
		return nil
	}

	return o.Message

}

// GetMessageOk returns a tuple with the Message field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DeletionError) GetMessageOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Message, true
}

// SetMessage sets field value
func (o *DeletionError) SetMessage(v string) {

	o.Message = &v

}

// HasMessage returns a boolean if a field has been set.
func (o *DeletionError) HasMessage() bool {
	if o != nil && o.Message != nil {
		return true
	}

	return false
}

func (o DeletionError) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Key != nil {
		toSerialize["Key"] = o.Key
	}

	if o.VersionId != nil {
		toSerialize["VersionId"] = o.VersionId
	}

	if o.Code != nil {
		toSerialize["Code"] = o.Code
	}

	if o.Message != nil {
		toSerialize["Message"] = o.Message
	}

	return json.Marshal(toSerialize)
}

type NullableDeletionError struct {
	value *DeletionError
	isSet bool
}

func (v NullableDeletionError) Get() *DeletionError {
	return v.value
}

func (v *NullableDeletionError) Set(val *DeletionError) {
	v.value = val
	v.isSet = true
}

func (v NullableDeletionError) IsSet() bool {
	return v.isSet
}

func (v *NullableDeletionError) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDeletionError(val *DeletionError) *NullableDeletionError {
	return &NullableDeletionError{value: val, isSet: true}
}

func (v NullableDeletionError) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDeletionError) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
