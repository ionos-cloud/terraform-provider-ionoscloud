/*
 * IONOS DBaaS PostgreSQL REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional PostgreSQL database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ClusterList List of clusters.
type ClusterList struct {
	// The offset specified in the request (if none was specified, the default offset is 0).
	Offset *int32 `json:"offset,omitempty"`
	// The limit specified in the request (if none was specified, the default limit is 100).
	Limit *int32           `json:"limit,omitempty"`
	Links *PaginationLinks `json:"links,omitempty"`
	Type  *ResourceType    `json:"type,omitempty"`
	// The unique ID of the resource.
	Id    *string            `json:"id,omitempty"`
	Items *[]ClusterResponse `json:"items,omitempty"`
}

// NewClusterList instantiates a new ClusterList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterList() *ClusterList {
	this := ClusterList{}

	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit

	return &this
}

// NewClusterListWithDefaults instantiates a new ClusterList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterListWithDefaults() *ClusterList {
	this := ClusterList{}
	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit
	return &this
}

// GetOffset returns the Offset field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterList) GetOffset() *int32 {
	if o == nil {
		return nil
	}

	return o.Offset

}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Offset, true
}

// SetOffset sets field value
func (o *ClusterList) SetOffset(v int32) {

	o.Offset = &v

}

// HasOffset returns a boolean if a field has been set.
func (o *ClusterList) HasOffset() bool {
	if o != nil && o.Offset != nil {
		return true
	}

	return false
}

// GetLimit returns the Limit field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *ClusterList) GetLimit() *int32 {
	if o == nil {
		return nil
	}

	return o.Limit

}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetLimitOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Limit, true
}

// SetLimit sets field value
func (o *ClusterList) SetLimit(v int32) {

	o.Limit = &v

}

// HasLimit returns a boolean if a field has been set.
func (o *ClusterList) HasLimit() bool {
	if o != nil && o.Limit != nil {
		return true
	}

	return false
}

// GetLinks returns the Links field value
// If the value is explicit nil, the zero value for PaginationLinks will be returned
func (o *ClusterList) GetLinks() *PaginationLinks {
	if o == nil {
		return nil
	}

	return o.Links

}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetLinksOk() (*PaginationLinks, bool) {
	if o == nil {
		return nil, false
	}

	return o.Links, true
}

// SetLinks sets field value
func (o *ClusterList) SetLinks(v PaginationLinks) {

	o.Links = &v

}

// HasLinks returns a boolean if a field has been set.
func (o *ClusterList) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for ResourceType will be returned
func (o *ClusterList) GetType() *ResourceType {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetTypeOk() (*ResourceType, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *ClusterList) SetType(v ResourceType) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *ClusterList) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterList) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *ClusterList) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ClusterList) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []ClusterResponse will be returned
func (o *ClusterList) GetItems() *[]ClusterResponse {
	if o == nil {
		return nil
	}

	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterList) GetItemsOk() (*[]ClusterResponse, bool) {
	if o == nil {
		return nil, false
	}

	return o.Items, true
}

// SetItems sets field value
func (o *ClusterList) SetItems(v []ClusterResponse) {

	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *ClusterList) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}

func (o ClusterList) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Offset != nil {
		toSerialize["offset"] = o.Offset
	}

	if o.Limit != nil {
		toSerialize["limit"] = o.Limit
	}

	if o.Links != nil {
		toSerialize["links"] = o.Links
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Items != nil {
		toSerialize["items"] = o.Items
	}

	return json.Marshal(toSerialize)
}

type NullableClusterList struct {
	value *ClusterList
	isSet bool
}

func (v NullableClusterList) Get() *ClusterList {
	return v.value
}

func (v *NullableClusterList) Set(val *ClusterList) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterList) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterList(val *ClusterList) *NullableClusterList {
	return &NullableClusterList{value: val, isSet: true}
}

func (v NullableClusterList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
