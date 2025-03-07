/*
 * IONOS Cloud - Object Storage Management API
 *
 * Object Storage Management API is a RESTful API that manages the object storage service configuration for IONOS Cloud.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// StorageClassCreate struct for StorageClassCreate
type StorageClassCreate struct {
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *StorageClass           `json:"properties"`
}

// NewStorageClassCreate instantiates a new StorageClassCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStorageClassCreate(properties StorageClass) *StorageClassCreate {
	this := StorageClassCreate{}

	this.Properties = &properties

	return &this
}

// NewStorageClassCreateWithDefaults instantiates a new StorageClassCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStorageClassCreateWithDefaults() *StorageClassCreate {
	this := StorageClassCreate{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *StorageClassCreate) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *StorageClassCreate) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *StorageClassCreate) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *StorageClassCreate) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for StorageClass will be returned
func (o *StorageClassCreate) GetProperties() *StorageClass {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *StorageClassCreate) GetPropertiesOk() (*StorageClass, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *StorageClassCreate) SetProperties(v StorageClass) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *StorageClassCreate) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o StorageClassCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableStorageClassCreate struct {
	value *StorageClassCreate
	isSet bool
}

func (v NullableStorageClassCreate) Get() *StorageClassCreate {
	return v.value
}

func (v *NullableStorageClassCreate) Set(val *StorageClassCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableStorageClassCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableStorageClassCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStorageClassCreate(val *StorageClassCreate) *NullableStorageClassCreate {
	return &NullableStorageClassCreate{value: val, isSet: true}
}

func (v NullableStorageClassCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStorageClassCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
