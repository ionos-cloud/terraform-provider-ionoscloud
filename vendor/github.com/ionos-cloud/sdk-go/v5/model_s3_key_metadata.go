/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"time"
)

// S3KeyMetadata struct for S3KeyMetadata
type S3KeyMetadata struct {
	// Resource's Entity Tag as defined in http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.11 . Entity Tag is also added as an 'ETag response header to requests which don't use 'depth' parameter.
	Etag *string `json:"etag,omitempty"`
	// The time the S3 key was created
	CreatedDate *IonosTime
}

// GetEtag returns the Etag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *S3KeyMetadata) GetEtag() *string {
	if o == nil {
		return nil
	}

	return o.Etag

}

// GetEtagOk returns a tuple with the Etag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *S3KeyMetadata) GetEtagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Etag, true
}

// SetEtag sets field value
func (o *S3KeyMetadata) SetEtag(v string) {

	o.Etag = &v

}

// HasEtag returns a boolean if a field has been set.
func (o *S3KeyMetadata) HasEtag() bool {
	if o != nil && o.Etag != nil {
		return true
	}

	return false
}

// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *S3KeyMetadata) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time

}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *S3KeyMetadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *S3KeyMetadata) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}

}

// HasCreatedDate returns a boolean if a field has been set.
func (o *S3KeyMetadata) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}

func (o S3KeyMetadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Etag != nil {
		toSerialize["etag"] = o.Etag
	}
	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}
	return json.Marshal(toSerialize)
}

type NullableS3KeyMetadata struct {
	value *S3KeyMetadata
	isSet bool
}

func (v NullableS3KeyMetadata) Get() *S3KeyMetadata {
	return v.value
}

func (v *NullableS3KeyMetadata) Set(val *S3KeyMetadata) {
	v.value = val
	v.isSet = true
}

func (v NullableS3KeyMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableS3KeyMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableS3KeyMetadata(val *S3KeyMetadata) *NullableS3KeyMetadata {
	return &NullableS3KeyMetadata{value: val, isSet: true}
}

func (v NullableS3KeyMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableS3KeyMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
