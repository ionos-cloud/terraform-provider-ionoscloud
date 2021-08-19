/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Cdroms struct for Cdroms
type Cdroms struct {
	// The resource's unique identifier
	Id *string `json:"id,omitempty"`
	// The type of object that has been created
	Type *Type `json:"type,omitempty"`
	// URL to the object representation (absolute path)
	Href *string `json:"href,omitempty"`
	// Array of items in that collection
	Items *[]Image `json:"items,omitempty"`
	// the offset (if specified in the request)
	Offset *float32 `json:"offset,omitempty"`
	// the limit (if specified in the request)
	Limit *float32 `json:"limit,omitempty"`
	Links *PaginationLinks `json:"_links,omitempty"`
}



// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cdroms) GetId() *string {
	if o == nil {
		return nil
	}


	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Id, true
}

// SetId sets field value
func (o *Cdroms) SetId(v string) {


	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *Cdroms) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}



// GetType returns the Type field value
// If the value is explicit nil, the zero value for Type will be returned
func (o *Cdroms) GetType() *Type {
	if o == nil {
		return nil
	}


	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetTypeOk() (*Type, bool) {
	if o == nil {
		return nil, false
	}


	return o.Type, true
}

// SetType sets field value
func (o *Cdroms) SetType(v Type) {


	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Cdroms) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}



// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cdroms) GetHref() *string {
	if o == nil {
		return nil
	}


	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Href, true
}

// SetHref sets field value
func (o *Cdroms) SetHref(v string) {


	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *Cdroms) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}



// GetItems returns the Items field value
// If the value is explicit nil, the zero value for []Image will be returned
func (o *Cdroms) GetItems() *[]Image {
	if o == nil {
		return nil
	}


	return o.Items

}

// GetItemsOk returns a tuple with the Items field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetItemsOk() (*[]Image, bool) {
	if o == nil {
		return nil, false
	}


	return o.Items, true
}

// SetItems sets field value
func (o *Cdroms) SetItems(v []Image) {


	o.Items = &v

}

// HasItems returns a boolean if a field has been set.
func (o *Cdroms) HasItems() bool {
	if o != nil && o.Items != nil {
		return true
	}

	return false
}



// GetOffset returns the Offset field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Cdroms) GetOffset() *float32 {
	if o == nil {
		return nil
	}


	return o.Offset

}

// GetOffsetOk returns a tuple with the Offset field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetOffsetOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Offset, true
}

// SetOffset sets field value
func (o *Cdroms) SetOffset(v float32) {


	o.Offset = &v

}

// HasOffset returns a boolean if a field has been set.
func (o *Cdroms) HasOffset() bool {
	if o != nil && o.Offset != nil {
		return true
	}

	return false
}



// GetLimit returns the Limit field value
// If the value is explicit nil, the zero value for float32 will be returned
func (o *Cdroms) GetLimit() *float32 {
	if o == nil {
		return nil
	}


	return o.Limit

}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetLimitOk() (*float32, bool) {
	if o == nil {
		return nil, false
	}


	return o.Limit, true
}

// SetLimit sets field value
func (o *Cdroms) SetLimit(v float32) {


	o.Limit = &v

}

// HasLimit returns a boolean if a field has been set.
func (o *Cdroms) HasLimit() bool {
	if o != nil && o.Limit != nil {
		return true
	}

	return false
}



// GetLinks returns the Links field value
// If the value is explicit nil, the zero value for PaginationLinks will be returned
func (o *Cdroms) GetLinks() *PaginationLinks {
	if o == nil {
		return nil
	}


	return o.Links

}

// GetLinksOk returns a tuple with the Links field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cdroms) GetLinksOk() (*PaginationLinks, bool) {
	if o == nil {
		return nil, false
	}


	return o.Links, true
}

// SetLinks sets field value
func (o *Cdroms) SetLinks(v PaginationLinks) {


	o.Links = &v

}

// HasLinks returns a boolean if a field has been set.
func (o *Cdroms) HasLinks() bool {
	if o != nil && o.Links != nil {
		return true
	}

	return false
}


func (o Cdroms) MarshalJSON() ([]byte, error) {
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

type NullableCdroms struct {
	value *Cdroms
	isSet bool
}

func (v NullableCdroms) Get() *Cdroms {
	return v.value
}

func (v *NullableCdroms) Set(val *Cdroms) {
	v.value = val
	v.isSet = true
}

func (v NullableCdroms) IsSet() bool {
	return v.isSet
}

func (v *NullableCdroms) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCdroms(val *Cdroms) *NullableCdroms {
	return &NullableCdroms{value: val, isSet: true}
}

func (v NullableCdroms) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCdroms) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


