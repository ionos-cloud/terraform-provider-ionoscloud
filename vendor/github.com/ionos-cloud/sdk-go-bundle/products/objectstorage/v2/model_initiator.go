/*
 * IONOS Object Storage API for contract-owned buckets
 *
 * ## Overview The IONOS Object Storage API for contract-owned buckets is a REST-based API that allows developers and applications to interact directly with IONOS' scalable storage solution, leveraging the S3 protocol for object storage operations. Its design ensures seamless compatibility with existing tools and libraries tailored for S3 systems.  ### API References - [S3 API Reference for contract-owned buckets](https://api.ionos.com/docs/s3-contract-owned-buckets/v2/) ### User documentation [IONOS Object Storage User Guide](https://docs.ionos.com/cloud/managed-services/s3-object-storage) * [Documentation on user-owned and contract-owned buckets](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/buckets) * [Documentation on S3 API Compatibility](https://docs.ionos.com/cloud/managed-services/s3-object-storage/concepts/s3-api-compatibility) * [S3 Tools](https://docs.ionos.com/cloud/managed-services/s3-object-storage/s3-tools)  ## Endpoints for contract-owned buckets | Location | Region Name | Bucket Type | Endpoint | | --- | --- | --- | --- | | **Berlin, Germany** | **eu-central-3** | Contract-owned | `https://s3.eu-central-3.ionoscloud.com` |  ## Changelog - 30.05.2024 Initial version
 *
 * API version: 2.0.2
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package objectstorage

import (
	"encoding/json"
)

import "encoding/xml"

// checks if the Initiator type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Initiator{}

// Initiator Container element that identifies who initiated the multipart upload.
type Initiator struct {
	XMLName xml.Name `xml:"Initiator"`
	// Container for the Contract Number of the owner.
	ID *int32 `json:"ID,omitempty" xml:"ID"`
	// Container for the display name of the owner.
	DisplayName *string `json:"DisplayName,omitempty" xml:"DisplayName"`
}

// NewInitiator instantiates a new Initiator object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInitiator() *Initiator {
	this := Initiator{}

	return &this
}

// NewInitiatorWithDefaults instantiates a new Initiator object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInitiatorWithDefaults() *Initiator {
	this := Initiator{}
	return &this
}

// GetID returns the ID field value if set, zero value otherwise.
func (o *Initiator) GetID() int32 {
	if o == nil || IsNil(o.ID) {
		var ret int32
		return ret
	}
	return *o.ID
}

// GetIDOk returns a tuple with the ID field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Initiator) GetIDOk() (*int32, bool) {
	if o == nil || IsNil(o.ID) {
		return nil, false
	}
	return o.ID, true
}

// HasID returns a boolean if a field has been set.
func (o *Initiator) HasID() bool {
	if o != nil && !IsNil(o.ID) {
		return true
	}

	return false
}

// SetID gets a reference to the given int32 and assigns it to the ID field.
func (o *Initiator) SetID(v int32) {
	o.ID = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *Initiator) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Initiator) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *Initiator) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *Initiator) SetDisplayName(v string) {
	o.DisplayName = &v
}

func (o Initiator) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Initiator) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ID) {
		toSerialize["ID"] = o.ID
	}
	if !IsNil(o.DisplayName) {
		toSerialize["DisplayName"] = o.DisplayName
	}
	return toSerialize, nil
}

type NullableInitiator struct {
	value *Initiator
	isSet bool
}

func (v NullableInitiator) Get() *Initiator {
	return v.value
}

func (v *NullableInitiator) Set(val *Initiator) {
	v.value = val
	v.isSet = true
}

func (v NullableInitiator) IsSet() bool {
	return v.isSet
}

func (v *NullableInitiator) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInitiator(val *Initiator) *NullableInitiator {
	return &NullableInitiator{value: val, isSet: true}
}

func (v NullableInitiator) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInitiator) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
