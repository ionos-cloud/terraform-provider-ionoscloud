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

// StorageClassEnsure struct for StorageClassEnsure
type StorageClassEnsure struct {
	// The StorageClass of the StorageClass.
	Id *string `json:"id"`
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *StorageClass           `json:"properties"`
}

// NewStorageClassEnsure instantiates a new StorageClassEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStorageClassEnsure(id string, properties StorageClass) *StorageClassEnsure {
	this := StorageClassEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewStorageClassEnsureWithDefaults instantiates a new StorageClassEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStorageClassEnsureWithDefaults() *StorageClassEnsure {
	this := StorageClassEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *StorageClassEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *StorageClassEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *StorageClassEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *StorageClassEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *StorageClassEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *StorageClassEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *StorageClassEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *StorageClassEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for StorageClass will be returned
func (o *StorageClassEnsure) GetProperties() *StorageClass {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *StorageClassEnsure) GetPropertiesOk() (*StorageClass, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *StorageClassEnsure) SetProperties(v StorageClass) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *StorageClassEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o StorageClassEnsure) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableStorageClassEnsure struct {
	value *StorageClassEnsure
	isSet bool
}

func (v NullableStorageClassEnsure) Get() *StorageClassEnsure {
	return v.value
}

func (v *NullableStorageClassEnsure) Set(val *StorageClassEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableStorageClassEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableStorageClassEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStorageClassEnsure(val *StorageClassEnsure) *NullableStorageClassEnsure {
	return &NullableStorageClassEnsure{value: val, isSet: true}
}

func (v NullableStorageClassEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStorageClassEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}