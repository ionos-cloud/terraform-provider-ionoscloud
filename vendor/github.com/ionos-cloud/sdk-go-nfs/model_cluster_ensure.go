/*
 * IONOS Cloud - Network File Storage API
 *
 * The RESTful API for managing Network File Storage.
 *
 * API version: 0.1.3
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ClusterEnsure struct for ClusterEnsure
type ClusterEnsure struct {
	// The identifier (UUID) of the cluster.
	Id *string `json:"id"`
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *Cluster                `json:"properties"`
}

// NewClusterEnsure instantiates a new ClusterEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterEnsure(id string, properties Cluster) *ClusterEnsure {
	this := ClusterEnsure{}

	this.Id = &id
	this.Properties = &properties

	return &this
}

// NewClusterEnsureWithDefaults instantiates a new ClusterEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterEnsureWithDefaults() *ClusterEnsure {
	this := ClusterEnsure{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterEnsure) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *ClusterEnsure) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *ClusterEnsure) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *ClusterEnsure) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterEnsure) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *ClusterEnsure) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *ClusterEnsure) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Cluster will be returned
func (o *ClusterEnsure) GetProperties() *Cluster {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterEnsure) GetPropertiesOk() (*Cluster, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *ClusterEnsure) SetProperties(v Cluster) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *ClusterEnsure) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o ClusterEnsure) MarshalJSON() ([]byte, error) {
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

type NullableClusterEnsure struct {
	value *ClusterEnsure
	isSet bool
}

func (v NullableClusterEnsure) Get() *ClusterEnsure {
	return v.value
}

func (v *NullableClusterEnsure) Set(val *ClusterEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterEnsure(val *ClusterEnsure) *NullableClusterEnsure {
	return &NullableClusterEnsure{value: val, isSet: true}
}

func (v NullableClusterEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
