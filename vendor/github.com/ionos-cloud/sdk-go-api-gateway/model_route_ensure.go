/*
 * IONOS Cloud - API Gateway
 *
 * API Gateway is an application that acts as a \"front door\" for backend services and APIs, handling client requests and routing them to the appropriate backend.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RouteEnsure struct for RouteEnsure
type RouteEnsure struct {
	// The ID (UUID) of the Route.
	Id *string `json:"id"`
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *Route                  `json:"properties"`
}

// NewRouteEnsure instantiates a new RouteEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRouteEnsure(id string, properties Route) *RouteEnsure {
	this := RouteEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewRouteEnsureWithDefaults instantiates a new RouteEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRouteEnsureWithDefaults() *RouteEnsure {
	this := RouteEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RouteEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RouteEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *RouteEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *RouteEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *RouteEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RouteEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *RouteEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *RouteEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Route will be returned
func (o *RouteEnsure) GetProperties() *Route {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RouteEnsure) GetPropertiesOk() (*Route, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *RouteEnsure) SetProperties(v Route) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *RouteEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o RouteEnsure) MarshalJSON() ([]byte, error) {
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

type NullableRouteEnsure struct {
	value *RouteEnsure
	isSet bool
}

func (v NullableRouteEnsure) Get() *RouteEnsure {
	return v.value
}

func (v *NullableRouteEnsure) Set(val *RouteEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableRouteEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableRouteEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRouteEnsure(val *RouteEnsure) *NullableRouteEnsure {
	return &NullableRouteEnsure{value: val, isSet: true}
}

func (v NullableRouteEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRouteEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
