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

// checks if the ObjectIdentifier type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ObjectIdentifier{}

// ObjectIdentifier Object Identifier is unique value to identify objects.
type ObjectIdentifier struct {
	XMLName xml.Name `xml:"Object"`
	// The object key.
	Key string `json:"Key" xml:"Key"`
	// VersionId for the specific version of the object to delete.
	VersionId *string `json:"VersionId,omitempty" xml:"VersionId"`
}

// NewObjectIdentifier instantiates a new ObjectIdentifier object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewObjectIdentifier(key string) *ObjectIdentifier {
	this := ObjectIdentifier{}

	this.Key = key

	return &this
}

// NewObjectIdentifierWithDefaults instantiates a new ObjectIdentifier object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewObjectIdentifierWithDefaults() *ObjectIdentifier {
	this := ObjectIdentifier{}
	return &this
}

// GetKey returns the Key field value
func (o *ObjectIdentifier) GetKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Key
}

// GetKeyOk returns a tuple with the Key field value
// and a boolean to check if the value has been set.
func (o *ObjectIdentifier) GetKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Key, true
}

// SetKey sets field value
func (o *ObjectIdentifier) SetKey(v string) {
	o.Key = v
}

// GetVersionId returns the VersionId field value if set, zero value otherwise.
func (o *ObjectIdentifier) GetVersionId() string {
	if o == nil || IsNil(o.VersionId) {
		var ret string
		return ret
	}
	return *o.VersionId
}

// GetVersionIdOk returns a tuple with the VersionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ObjectIdentifier) GetVersionIdOk() (*string, bool) {
	if o == nil || IsNil(o.VersionId) {
		return nil, false
	}
	return o.VersionId, true
}

// HasVersionId returns a boolean if a field has been set.
func (o *ObjectIdentifier) HasVersionId() bool {
	if o != nil && !IsNil(o.VersionId) {
		return true
	}

	return false
}

// SetVersionId gets a reference to the given string and assigns it to the VersionId field.
func (o *ObjectIdentifier) SetVersionId(v string) {
	o.VersionId = &v
}

func (o ObjectIdentifier) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ObjectIdentifier) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Key) {
		toSerialize["Key"] = o.Key
	}
	if !IsNil(o.VersionId) {
		toSerialize["VersionId"] = o.VersionId
	}
	return toSerialize, nil
}

type NullableObjectIdentifier struct {
	value *ObjectIdentifier
	isSet bool
}

func (v NullableObjectIdentifier) Get() *ObjectIdentifier {
	return v.value
}

func (v *NullableObjectIdentifier) Set(val *ObjectIdentifier) {
	v.value = val
	v.isSet = true
}

func (v NullableObjectIdentifier) IsSet() bool {
	return v.isSet
}

func (v *NullableObjectIdentifier) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableObjectIdentifier(val *ObjectIdentifier) *NullableObjectIdentifier {
	return &NullableObjectIdentifier{value: val, isSet: true}
}

func (v NullableObjectIdentifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableObjectIdentifier) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
