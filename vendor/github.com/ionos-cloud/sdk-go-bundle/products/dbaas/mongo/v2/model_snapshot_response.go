/*
 * IONOS DBaaS MongoDB REST API
 *
 * With IONOS Cloud Database as a Service, you have the ability to quickly set up and manage a MongoDB database. You can also delete clusters, manage backups and users via the API.  MongoDB is an open source, cross-platform, document-oriented database program. Classified as a NoSQL database program, it uses JSON-like documents with optional schemas.  The MongoDB API allows you to create additional database clusters or modify existing ones. Both tools, the Data Center Designer (DCD) and the API use the same concepts consistently and are well suited for smooth and intuitive use.
 *
 * API version: 1.0.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mongo

import (
	"encoding/json"
)

// checks if the SnapshotResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SnapshotResponse{}

// SnapshotResponse A database snapshot.
type SnapshotResponse struct {
	Type *ResourceType `json:"type,omitempty"`
	// The unique ID of the resource.
	Id         *string             `json:"id,omitempty"`
	Properties *SnapshotProperties `json:"properties,omitempty"`
}

// NewSnapshotResponse instantiates a new SnapshotResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSnapshotResponse() *SnapshotResponse {
	this := SnapshotResponse{}

	return &this
}

// NewSnapshotResponseWithDefaults instantiates a new SnapshotResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSnapshotResponseWithDefaults() *SnapshotResponse {
	this := SnapshotResponse{}
	return &this
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *SnapshotResponse) GetType() ResourceType {
	if o == nil || IsNil(o.Type) {
		var ret ResourceType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SnapshotResponse) GetTypeOk() (*ResourceType, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *SnapshotResponse) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given ResourceType and assigns it to the Type field.
func (o *SnapshotResponse) SetType(v ResourceType) {
	o.Type = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *SnapshotResponse) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SnapshotResponse) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *SnapshotResponse) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *SnapshotResponse) SetId(v string) {
	o.Id = &v
}

// GetProperties returns the Properties field value if set, zero value otherwise.
func (o *SnapshotResponse) GetProperties() SnapshotProperties {
	if o == nil || IsNil(o.Properties) {
		var ret SnapshotProperties
		return ret
	}
	return *o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SnapshotResponse) GetPropertiesOk() (*SnapshotProperties, bool) {
	if o == nil || IsNil(o.Properties) {
		return nil, false
	}
	return o.Properties, true
}

// HasProperties returns a boolean if a field has been set.
func (o *SnapshotResponse) HasProperties() bool {
	if o != nil && !IsNil(o.Properties) {
		return true
	}

	return false
}

// SetProperties gets a reference to the given SnapshotProperties and assigns it to the Properties field.
func (o *SnapshotResponse) SetProperties(v SnapshotProperties) {
	o.Properties = &v
}

func (o SnapshotResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SnapshotResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableSnapshotResponse struct {
	value *SnapshotResponse
	isSet bool
}

func (v NullableSnapshotResponse) Get() *SnapshotResponse {
	return v.value
}

func (v *NullableSnapshotResponse) Set(val *SnapshotResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableSnapshotResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableSnapshotResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSnapshotResponse(val *SnapshotResponse) *NullableSnapshotResponse {
	return &NullableSnapshotResponse{value: val, isSet: true}
}

func (v NullableSnapshotResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSnapshotResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
