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

// AccessKeyRead struct for AccessKeyRead
type AccessKeyRead struct {
	// The ID (UUID) of the AccessKey.
	Id *string `json:"id"`
	// The type of the resource.
	Type *string `json:"type"`
	// The URL of the AccessKey.
	Href       *string                       `json:"href"`
	Metadata   *MetadataWithSupportedRegions `json:"metadata"`
	Properties *AccessKey                    `json:"properties"`
}

// NewAccessKeyRead instantiates a new AccessKeyRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAccessKeyRead(id string, type_ string, href string, metadata MetadataWithSupportedRegions, properties AccessKey) *AccessKeyRead {
	this := AccessKeyRead{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href
	this.Metadata = &metadata
	this.Properties = &properties

	return &this
}

// NewAccessKeyReadWithDefaults instantiates a new AccessKeyRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAccessKeyReadWithDefaults() *AccessKeyRead {
	this := AccessKeyRead{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKeyRead) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKeyRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *AccessKeyRead) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *AccessKeyRead) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKeyRead) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKeyRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *AccessKeyRead) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *AccessKeyRead) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *AccessKeyRead) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKeyRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *AccessKeyRead) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *AccessKeyRead) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for MetadataWithSupportedRegions will be returned
func (o *AccessKeyRead) GetMetadata() *MetadataWithSupportedRegions {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKeyRead) GetMetadataOk() (*MetadataWithSupportedRegions, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *AccessKeyRead) SetMetadata(v MetadataWithSupportedRegions) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *AccessKeyRead) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for AccessKey will be returned
func (o *AccessKeyRead) GetProperties() *AccessKey {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *AccessKeyRead) GetPropertiesOk() (*AccessKey, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *AccessKeyRead) SetProperties(v AccessKey) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *AccessKeyRead) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o AccessKeyRead) MarshalJSON() ([]byte, error) {
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

type NullableAccessKeyRead struct {
	value *AccessKeyRead
	isSet bool
}

func (v NullableAccessKeyRead) Get() *AccessKeyRead {
	return v.value
}

func (v *NullableAccessKeyRead) Set(val *AccessKeyRead) {
	v.value = val
	v.isSet = true
}

func (v NullableAccessKeyRead) IsSet() bool {
	return v.isSet
}

func (v *NullableAccessKeyRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAccessKeyRead(val *AccessKeyRead) *NullableAccessKeyRead {
	return &NullableAccessKeyRead{value: val, isSet: true}
}

func (v NullableAccessKeyRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAccessKeyRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}