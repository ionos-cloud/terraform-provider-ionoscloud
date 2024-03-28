/*
 * IONOS DBaaS MariaDB REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional MariaDB database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Pagination Pagination information in list responses.
type Pagination struct {
	// The offset specified in the request (if none was specified, the default offset is 0).
	Offset *int32 `json:"offset,omitempty"`
	// The limit specified in the request (if none was specified, the default limit is 100).
	Limit *int32 `json:"limit,omitempty"`
	// The total number of elements matching the request (without pagination).
	Total *int32           `json:"total,omitempty"`
	Links *PaginationLinks `json:"_links,omitempty"`
}

// NewPagination instantiates a new Pagination object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPagination() *Pagination {
	this := Pagination{}

	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit
	var total int32 = 0
	this.Total = &total

	return &this
}

// NewPaginationWithDefaults instantiates a new Pagination object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPaginationWithDefaults() *Pagination {
	this := Pagination{}
	var offset int32 = 0
	this.Offset = &offset
	var limit int32 = 100
	this.Limit = &limit
	var total int32 = 0
	this.Total = &total
	return &this
}

// GetOffset returns the Offset field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Pagination) GetOffset() *int32 {
	if o == nil {
		return nil
	}

	return o.Offset

}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Pagination) GetOffsetOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Offset, true
}

// SetOffset sets field value
func (o *Pagination) SetOffset(v int32) {

	o.Offset = &v

}

// HasOffset returns a boolean if a field has been set.
func (o *Pagination) HasOffset() bool {
	if o != nil && o.Offset != nil {
		return true
	}

	return false
}

// GetLimit returns the Limit field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Pagination) GetLimit() *int32 {
	if o == nil {
		return nil
	}

	return o.Limit

}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Pagination) GetLimitOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Limit, true
}

// SetLimit sets field value
func (o *Pagination) SetLimit(v int32) {

	o.Limit = &v

}

// HasLimit returns a boolean if a field has been set.
func (o *Pagination) HasLimit() bool {
	if o != nil && o.Limit != nil {
		return true
	}

	return false
}

// GetTotal returns the Total field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Pagination) GetTotal() *int32 {
	if o == nil {
		return nil
	}

	return o.Total

}

// GetTotalOk returns a tuple with the Total field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Pagination) GetTotalOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Total, true
}

// SetTotal sets field value
func (o *Pagination) SetTotal(v int32) {

	o.Total = &v

}

// HasTotal returns a boolean if a field has been set.
func (o *Pagination) HasTotal() bool {
	if o != nil && o.Total != nil {
		return true
	}

	return false
}

// GetLinks returns the Links field value
// If the value is explicit nil, the zero value for PaginationLinks will be returned
func (o *Pagination) GetLinks() *PaginationLinks {
	if o == nil {
		return nil
	}

	return o.Links

}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Pagination) GetLinksOk() (*PaginationLinks, bool) {
	if o == nil {
		return nil, false
	}

	return o.Links, true
}

// SetLinks sets field value
func (o *Pagination) SetLinks(v PaginationLinks) {

	o.Links = &v

}

// HasLinks returns a boolean if a field has been set.
func (o *Pagination) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}

func (o Pagination) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Offset != nil {
		toSerialize["offset"] = o.Offset
	}

	if o.Limit != nil {
		toSerialize["limit"] = o.Limit
	}

	if o.Total != nil {
		toSerialize["total"] = o.Total
	}

	if o.Links != nil {
		toSerialize["_links"] = o.Links
	}

	return json.Marshal(toSerialize)
}

type NullablePagination struct {
	value *Pagination
	isSet bool
}

func (v NullablePagination) Get() *Pagination {
	return v.value
}

func (v *NullablePagination) Set(val *Pagination) {
	v.value = val
	v.isSet = true
}

func (v NullablePagination) IsSet() bool {
	return v.isSet
}

func (v *NullablePagination) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePagination(val *Pagination) *NullablePagination {
	return &NullablePagination{value: val, isSet: true}
}

func (v NullablePagination) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePagination) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
