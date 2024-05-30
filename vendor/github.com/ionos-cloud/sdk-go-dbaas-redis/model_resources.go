/*
 * Redis DB API
 *
 * Redis Database API
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Resources The resources of the individual replicas.
type Resources struct {
	// The number of CPU cores per instance.
	Cores *int32 `json:"cores"`
	// The amount of memory per instance in gigabytes (GB).
	Ram *int32 `json:"ram"`
	// The size of the storage in GB. The size is derived from the amount of RAM and the persistence mode and is not configurable.
	Storage *int32 `json:"storage"`
}

// NewResources instantiates a new Resources object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewResources(cores int32, ram int32, storage int32) *Resources {
	this := Resources{}

	this.Cores = &cores
	this.Ram = &ram
	this.Storage = &storage

	return &this
}

// NewResourcesWithDefaults instantiates a new Resources object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewResourcesWithDefaults() *Resources {
	this := Resources{}
	return &this
}

// GetCores returns the Cores field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Resources) GetCores() *int32 {
	if o == nil {
		return nil
	}

	return o.Cores

}

// GetCoresOk returns a tuple with the Cores field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Resources) GetCoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Cores, true
}

// SetCores sets field value
func (o *Resources) SetCores(v int32) {

	o.Cores = &v

}

// HasCores returns a boolean if a field has been set.
func (o *Resources) HasCores() bool {
	if o != nil && o.Cores != nil {
		return true
	}

	return false
}

// GetRam returns the Ram field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Resources) GetRam() *int32 {
	if o == nil {
		return nil
	}

	return o.Ram

}

// GetRamOk returns a tuple with the Ram field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Resources) GetRamOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Ram, true
}

// SetRam sets field value
func (o *Resources) SetRam(v int32) {

	o.Ram = &v

}

// HasRam returns a boolean if a field has been set.
func (o *Resources) HasRam() bool {
	if o != nil && o.Ram != nil {
		return true
	}

	return false
}

// GetStorage returns the Storage field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Resources) GetStorage() *int32 {
	if o == nil {
		return nil
	}

	return o.Storage

}

// GetStorageOk returns a tuple with the Storage field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Resources) GetStorageOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Storage, true
}

// SetStorage sets field value
func (o *Resources) SetStorage(v int32) {

	o.Storage = &v

}

// HasStorage returns a boolean if a field has been set.
func (o *Resources) HasStorage() bool {
	if o != nil && o.Storage != nil {
		return true
	}

	return false
}

func (o Resources) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Cores != nil {
		toSerialize["cores"] = o.Cores
	}

	if o.Ram != nil {
		toSerialize["ram"] = o.Ram
	}

	if o.Storage != nil {
		toSerialize["storage"] = o.Storage
	}

	return json.Marshal(toSerialize)
}

type NullableResources struct {
	value *Resources
	isSet bool
}

func (v NullableResources) Get() *Resources {
	return v.value
}

func (v *NullableResources) Set(val *Resources) {
	v.value = val
	v.isSet = true
}

func (v NullableResources) IsSet() bool {
	return v.isSet
}

func (v *NullableResources) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResources(val *Resources) *NullableResources {
	return &NullableResources{value: val, isSet: true}
}

func (v NullableResources) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResources) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
