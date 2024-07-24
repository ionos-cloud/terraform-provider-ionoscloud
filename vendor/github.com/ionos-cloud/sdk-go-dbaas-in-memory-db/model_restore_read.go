/*
 * In-Memory DB API
 *
 * API description for the IONOS In-Memory DB
 *
 * API version: 1.0.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RestoreRead struct for RestoreRead
type RestoreRead struct {
	// The ID (UUID) of the Restore.
	Id *string `json:"id"`
	// The type of the resource.
	Type *string `json:"type"`
	// The URL of the Restore.
	Href       *string          `json:"href"`
	Metadata   *RestoreMetadata `json:"metadata"`
	Properties *Restore         `json:"properties"`
}

// NewRestoreRead instantiates a new RestoreRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRestoreRead(id string, type_ string, href string, metadata RestoreMetadata, properties Restore) *RestoreRead {
	this := RestoreRead{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href
	this.Metadata = &metadata
	this.Properties = &properties

	return &this
}

// NewRestoreReadWithDefaults instantiates a new RestoreRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRestoreReadWithDefaults() *RestoreRead {
	this := RestoreRead{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RestoreRead) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *RestoreRead) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *RestoreRead) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RestoreRead) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *RestoreRead) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *RestoreRead) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RestoreRead) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *RestoreRead) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *RestoreRead) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for RestoreMetadata will be returned
func (o *RestoreRead) GetMetadata() *RestoreMetadata {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreRead) GetMetadataOk() (*RestoreMetadata, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *RestoreRead) SetMetadata(v RestoreMetadata) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *RestoreRead) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Restore will be returned
func (o *RestoreRead) GetProperties() *Restore {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RestoreRead) GetPropertiesOk() (*Restore, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *RestoreRead) SetProperties(v Restore) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *RestoreRead) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o RestoreRead) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableRestoreRead struct {
	value *RestoreRead
	isSet bool
}

func (v NullableRestoreRead) Get() *RestoreRead {
	return v.value
}

func (v *NullableRestoreRead) Set(val *RestoreRead) {
	v.value = val
	v.isSet = true
}

func (v NullableRestoreRead) IsSet() bool {
	return v.isSet
}

func (v *NullableRestoreRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRestoreRead(val *RestoreRead) *NullableRestoreRead {
	return &NullableRestoreRead{value: val, isSet: true}
}

func (v NullableRestoreRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRestoreRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}