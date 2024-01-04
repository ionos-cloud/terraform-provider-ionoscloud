/*
 * Container Registry service
 *
 * ## Overview Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls. ## Changelog ### 1.1.0  - Added new endpoints for Repositories  - Added new endpoints for Artifacts  - Added new endpoints for Vulnerabilities  - Added registry vulnerabilityScanning feature
 *
 * API version: 1.1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RegistryArtifactsReadList struct for RegistryArtifactsReadList
type RegistryArtifactsReadList struct {
	Id    *string         `json:"id"`
	Type  *string         `json:"type"`
	Href  *string         `json:"href"`
	Items *[]ArtifactRead `json:"items,omitempty"`
	// The offset specified in the request (if none was specified, the default offset is 0) (not implemented yet).
	Offset *int32 `json:"offset"`
	// The limit specified in the request (if none was specified, use the endpoint's default pagination limit) (not implemented yet, always return number of items).
	Limit *int32 `json:"limit"`
	Links *Links `json:"_links"`
}

// NewRegistryArtifactsReadList instantiates a new RegistryArtifactsReadList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegistryArtifactsReadList(id string, type_ string, href string, offset int32, limit int32, links Links) *RegistryArtifactsReadList {
	this := RegistryArtifactsReadList{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href
	this.Offset = &offset
	this.Limit = &limit
	this.Links = &links

	return &this
}

// NewRegistryArtifactsReadListWithDefaults instantiates a new RegistryArtifactsReadList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegistryArtifactsReadListWithDefaults() *RegistryArtifactsReadList {
	this := RegistryArtifactsReadList{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryArtifactsReadList) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *RegistryArtifactsReadList) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryArtifactsReadList) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *RegistryArtifactsReadList) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryArtifactsReadList) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *RegistryArtifactsReadList) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []ArtifactRead will be returned
func (o *RegistryArtifactsReadList) GetItems() *[]ArtifactRead {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetItemsOk() (*[]ArtifactRead, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *RegistryArtifactsReadList) SetItems(v []ArtifactRead) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

// GetOffset returns the Offset field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *RegistryArtifactsReadList) GetOffset() *int32 {
	if o == nil {
		return nil
	}

	return o.Offset

}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Offset, true
}

// SetOffset sets field value
func (o *RegistryArtifactsReadList) SetOffset(v int32) {

	o.Offset = &v

}

// HasOffset returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasOffset() bool {
	if o != nil && o.Offset != nil {
		return true
	}

	return false
}

// GetLimit returns the Limit field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *RegistryArtifactsReadList) GetLimit() *int32 {
	if o == nil {
		return nil
	}

	return o.Limit

}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetLimitOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Limit, true
}

// SetLimit sets field value
func (o *RegistryArtifactsReadList) SetLimit(v int32) {

	o.Limit = &v

}

// HasLimit returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasLimit() bool {
	if o != nil && o.Limit != nil {
		return true
	}

	return false
}

// GetLinks returns the Links field value
// If the value is explicit nil, the zero value for Links will be returned
func (o *RegistryArtifactsReadList) GetLinks() *Links {
	if o == nil {
		return nil
	}

	return o.Links

}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryArtifactsReadList) GetLinksOk() (*Links, bool) {
	if o == nil {
		return nil, false
	}

	return o.Links, true
}

// SetLinks sets field value
func (o *RegistryArtifactsReadList) SetLinks(v Links) {

	o.Links = &v

}

// HasLinks returns a boolean if a field has been set.
func (o *RegistryArtifactsReadList) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

func (o RegistryArtifactsReadList) MarshalJSON() ([]byte, error) {
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

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	if o.Offset != nil {
		toSerialize["offset"] = o.Offset
	}

	if o.Limit != nil {
		toSerialize["limit"] = o.Limit
	}

	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}

	return json.Marshal(toSerialize)
}

type NullableRegistryArtifactsReadList struct {
	value *RegistryArtifactsReadList
	isSet bool
}

func (v NullableRegistryArtifactsReadList) Get() *RegistryArtifactsReadList {
	return v.value
}

func (v *NullableRegistryArtifactsReadList) Set(val *RegistryArtifactsReadList) {
	v.value = val
	v.isSet = true
}

func (v NullableRegistryArtifactsReadList) IsSet() bool {
	return v.isSet
}

func (v *NullableRegistryArtifactsReadList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegistryArtifactsReadList(val *RegistryArtifactsReadList) *NullableRegistryArtifactsReadList {
	return &NullableRegistryArtifactsReadList{value: val, isSet: true}
}

func (v NullableRegistryArtifactsReadList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegistryArtifactsReadList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
