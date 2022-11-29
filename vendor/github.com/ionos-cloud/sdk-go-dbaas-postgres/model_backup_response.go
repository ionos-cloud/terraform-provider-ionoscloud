/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// BackupResponse A database backup.
type BackupResponse struct {
	Type *ResourceType `json:"type,omitempty"`
	// The unique ID of the resource.
	Id         *string         `json:"id,omitempty"`
	Metadata   *BackupMetadata `json:"metadata,omitempty"`
	Properties *ClusterBackup  `json:"properties,omitempty"`
}

// NewBackupResponse instantiates a new BackupResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupResponse() *BackupResponse {
	this := BackupResponse{}

	return &this
}

// NewBackupResponseWithDefaults instantiates a new BackupResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupResponseWithDefaults() *BackupResponse {
	this := BackupResponse{}
	return &this
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for ResourceType will be returned
func (o *BackupResponse) GetType() *ResourceType {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BackupResponse) GetTypeOk() (*ResourceType, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *BackupResponse) SetType(v ResourceType) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *BackupResponse) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *BackupResponse) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BackupResponse) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *BackupResponse) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *BackupResponse) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for BackupMetadata will be returned
func (o *BackupResponse) GetMetadata() *BackupMetadata {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BackupResponse) GetMetadataOk() (*BackupMetadata, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *BackupResponse) SetMetadata(v BackupMetadata) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *BackupResponse) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for ClusterBackup will be returned
func (o *BackupResponse) GetProperties() *ClusterBackup {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *BackupResponse) GetPropertiesOk() (*ClusterBackup, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *BackupResponse) SetProperties(v ClusterBackup) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *BackupResponse) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o BackupResponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

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

type NullableBackupResponse struct {
	value *BackupResponse
	isSet bool
}

func (v NullableBackupResponse) Get() *BackupResponse {
	return v.value
}

func (v *NullableBackupResponse) Set(val *BackupResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableBackupResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableBackupResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBackupResponse(val *BackupResponse) *NullableBackupResponse {
	return &NullableBackupResponse{value: val, isSet: true}
}

func (v NullableBackupResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBackupResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
