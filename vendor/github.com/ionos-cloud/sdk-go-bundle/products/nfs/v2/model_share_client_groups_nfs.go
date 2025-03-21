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

// checks if the ShareClientGroupsNfs type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ShareClientGroupsNfs{}

// ShareClientGroupsNfs NFS specific configurations.
type ShareClientGroupsNfs struct {
	// The squash mode for the export. The squash mode can be: * `none` - No squash mode. no mapping (no_all_squash,no_root_squash). * `root-anonymous` - Map root user to anonymous uid (root_squash,anonuid=<uid>,anongid=<gid>). * `all-anonymous` - Map all users to anonymous uid (all_squash,anonuid=<uid>,anongid=<gid>).
	Squash *string `json:"squash,omitempty"`
}

// NewShareClientGroupsNfs instantiates a new ShareClientGroupsNfs object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewShareClientGroupsNfs() *ShareClientGroupsNfs {
	this := ShareClientGroupsNfs{}

	var squash string = "none"
	this.Squash = &squash

	return &this
}

// NewShareClientGroupsNfsWithDefaults instantiates a new ShareClientGroupsNfs object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewShareClientGroupsNfsWithDefaults() *ShareClientGroupsNfs {
	this := ShareClientGroupsNfs{}
	var squash string = "none"
	this.Squash = &squash
	return &this
}

// GetSquash returns the Squash field value if set, zero value otherwise.
func (o *ShareClientGroupsNfs) GetSquash() string {
	if o == nil || IsNil(o.Squash) {
		var ret string
		return ret
	}
	return *o.Squash
}

// GetSquashOk returns a tuple with the Squash field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ShareClientGroupsNfs) GetSquashOk() (*string, bool) {
	if o == nil || IsNil(o.Squash) {
		return nil, false
	}
	return o.Squash, true
}

// HasSquash returns a boolean if a field has been set.
func (o *ShareClientGroupsNfs) HasSquash() bool {
	if o != nil && !IsNil(o.Squash) {
		return true
	}

	return false
}

// SetSquash gets a reference to the given string and assigns it to the Squash field.
func (o *ShareClientGroupsNfs) SetSquash(v string) {
	o.Squash = &v
}

func (o ShareClientGroupsNfs) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Squash) {
		toSerialize["squash"] = o.Squash
	}
	return toSerialize, nil
}

type NullableShareClientGroupsNfs struct {
	value *ShareClientGroupsNfs
	isSet bool
}

func (v NullableShareClientGroupsNfs) Get() *ShareClientGroupsNfs {
	return v.value
}

func (v *NullableShareClientGroupsNfs) Set(val *ShareClientGroupsNfs) {
	v.value = val
	v.isSet = true
}

func (v NullableShareClientGroupsNfs) IsSet() bool {
	return v.isSet
}

func (v *NullableShareClientGroupsNfs) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableShareClientGroupsNfs(val *ShareClientGroupsNfs) *NullableShareClientGroupsNfs {
	return &NullableShareClientGroupsNfs{value: val, isSet: true}
}

func (v NullableShareClientGroupsNfs) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableShareClientGroupsNfs) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
