/*
 * IONOS Cloud - Network File Storage API
 *
 * The RESTful API for managing Network File Storage.
 *
 * API version: 0.1.3
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nfs

import (
	"encoding/json"
)

// checks if the ClusterNfs type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ClusterNfs{}

// ClusterNfs struct for ClusterNfs
type ClusterNfs struct {
	// The version of the NFS cluster, that is supported at minimum.  Currently supported version: * `4.2` - NFSv4.2
	MinVersion *string `json:"minVersion,omitempty"`
}

// NewClusterNfs instantiates a new ClusterNfs object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewClusterNfs() *ClusterNfs {
	this := ClusterNfs{}

	var minVersion string = "4.2"
	this.MinVersion = &minVersion

	return &this
}

// NewClusterNfsWithDefaults instantiates a new ClusterNfs object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterNfsWithDefaults() *ClusterNfs {
	this := ClusterNfs{}
	var minVersion string = "4.2"
	this.MinVersion = &minVersion
	return &this
}

// GetMinVersion returns the MinVersion field value if set, zero value otherwise.
func (o *ClusterNfs) GetMinVersion() string {
	if o == nil || IsNil(o.MinVersion) {
		var ret string
		return ret
	}
	return *o.MinVersion
}

// GetMinVersionOk returns a tuple with the MinVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ClusterNfs) GetMinVersionOk() (*string, bool) {
	if o == nil || IsNil(o.MinVersion) {
		return nil, false
	}
	return o.MinVersion, true
}

// HasMinVersion returns a boolean if a field has been set.
func (o *ClusterNfs) HasMinVersion() bool {
	if o != nil && !IsNil(o.MinVersion) {
		return true
	}

	return false
}

// SetMinVersion gets a reference to the given string and assigns it to the MinVersion field.
func (o *ClusterNfs) SetMinVersion(v string) {
	o.MinVersion = &v
}

func (o ClusterNfs) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.MinVersion) {
		toSerialize["minVersion"] = o.MinVersion
	}
	return toSerialize, nil
}

type NullableClusterNfs struct {
	value *ClusterNfs
	isSet bool
}

func (v NullableClusterNfs) Get() *ClusterNfs {
	return v.value
}

func (v *NullableClusterNfs) Set(val *ClusterNfs) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterNfs) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterNfs) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterNfs(val *ClusterNfs) *NullableClusterNfs {
	return &NullableClusterNfs{value: val, isSet: true}
}

func (v NullableClusterNfs) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterNfs) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
