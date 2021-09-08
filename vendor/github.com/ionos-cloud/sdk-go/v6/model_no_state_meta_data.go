/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"time"
)

// NoStateMetaData struct for NoStateMetaData
type NoStateMetaData struct {
	// Resource's Entity Tag as defined in http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.11 . Entity Tag is also added as an 'ETag response header to requests which don't use 'depth' parameter. 
	Etag *string `json:"etag,omitempty"`
	// The time the Resource was created
	CreatedDate *IonosTime
	// The user who has created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`
	// The user id of the user who has created the resource.
	CreatedByUserId *string `json:"createdByUserId,omitempty"`
	// The last time the resource has been modified
	LastModifiedDate *IonosTime
	// The user who last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// The user id of the user who has last modified the resource.
	LastModifiedByUserId *string `json:"lastModifiedByUserId,omitempty"`
}



// GetEtag returns the Etag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *NoStateMetaData) GetEtag() *string {
	if o == nil {
		return nil
	}


	return o.Etag

}

// GetEtagOk returns a tuple with the Etag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetEtagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Etag, true
}

// SetEtag sets field value
func (o *NoStateMetaData) SetEtag(v string) {


	o.Etag = &v

}

// HasEtag returns a boolean if a field has been set.
func (o *NoStateMetaData) HasEtag() bool {
	if o != nil && o.Etag != nil {
		return true
	}

	return false
}



// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *NoStateMetaData) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time


}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *NoStateMetaData) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}


}

// HasCreatedDate returns a boolean if a field has been set.
func (o *NoStateMetaData) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}



// GetCreatedBy returns the CreatedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *NoStateMetaData) GetCreatedBy() *string {
	if o == nil {
		return nil
	}


	return o.CreatedBy

}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetCreatedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *NoStateMetaData) SetCreatedBy(v string) {


	o.CreatedBy = &v

}

// HasCreatedBy returns a boolean if a field has been set.
func (o *NoStateMetaData) HasCreatedBy() bool {
	if o != nil && o.CreatedBy != nil {
		return true
	}

	return false
}



// GetCreatedByUserId returns the CreatedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *NoStateMetaData) GetCreatedByUserId() *string {
	if o == nil {
		return nil
	}


	return o.CreatedByUserId

}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.CreatedByUserId, true
}

// SetCreatedByUserId sets field value
func (o *NoStateMetaData) SetCreatedByUserId(v string) {


	o.CreatedByUserId = &v

}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *NoStateMetaData) HasCreatedByUserId() bool {
	if o != nil && o.CreatedByUserId != nil {
		return true
	}

	return false
}



// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *NoStateMetaData) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time


}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *NoStateMetaData) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}


}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *NoStateMetaData) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}



// GetLastModifiedBy returns the LastModifiedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *NoStateMetaData) GetLastModifiedBy() *string {
	if o == nil {
		return nil
	}


	return o.LastModifiedBy

}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetLastModifiedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LastModifiedBy, true
}

// SetLastModifiedBy sets field value
func (o *NoStateMetaData) SetLastModifiedBy(v string) {


	o.LastModifiedBy = &v

}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *NoStateMetaData) HasLastModifiedBy() bool {
	if o != nil && o.LastModifiedBy != nil {
		return true
	}

	return false
}



// GetLastModifiedByUserId returns the LastModifiedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *NoStateMetaData) GetLastModifiedByUserId() *string {
	if o == nil {
		return nil
	}


	return o.LastModifiedByUserId

}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *NoStateMetaData) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.LastModifiedByUserId, true
}

// SetLastModifiedByUserId sets field value
func (o *NoStateMetaData) SetLastModifiedByUserId(v string) {


	o.LastModifiedByUserId = &v

}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *NoStateMetaData) HasLastModifiedByUserId() bool {
	if o != nil && o.LastModifiedByUserId != nil {
		return true
	}

	return false
}


func (o NoStateMetaData) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Etag != nil {
		toSerialize["etag"] = o.Etag
	}
	

	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}
	

	if o.CreatedBy != nil {
		toSerialize["createdBy"] = o.CreatedBy
	}
	

	if o.CreatedByUserId != nil {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}
	

	if o.LastModifiedDate != nil {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}
	

	if o.LastModifiedBy != nil {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}
	

	if o.LastModifiedByUserId != nil {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}
	
	return json.Marshal(toSerialize)
}

type NullableNoStateMetaData struct {
	value *NoStateMetaData
	isSet bool
}

func (v NullableNoStateMetaData) Get() *NoStateMetaData {
	return v.value
}

func (v *NullableNoStateMetaData) Set(val *NoStateMetaData) {
	v.value = val
	v.isSet = true
}

func (v NullableNoStateMetaData) IsSet() bool {
	return v.isSet
}

func (v *NullableNoStateMetaData) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNoStateMetaData(val *NoStateMetaData) *NullableNoStateMetaData {
	return &NullableNoStateMetaData{value: val, isSet: true}
}

func (v NullableNoStateMetaData) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNoStateMetaData) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


