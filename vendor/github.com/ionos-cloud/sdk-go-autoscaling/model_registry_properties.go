/*
 * Container Registry service (CloudAPI)
 *
 * Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their manage Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls.
 *
 * API version: 1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// RegistryProperties struct for RegistryProperties
type RegistryProperties struct {
	GarbageCollectionSchedule *WeeklySchedule `json:"garbageCollectionSchedule,omitempty"`
	Hostname                  *string         `json:"hostname,omitempty"`
	Location                  *string         `json:"location"`
	MaintenanceWindow         *WeeklySchedule `json:"maintenanceWindow,omitempty"`
	Name                      *string         `json:"name"`
	StorageUsage              *StorageUsage   `json:"storageUsage,omitempty"`
}

// NewRegistryProperties instantiates a new RegistryProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegistryProperties(location string, name string) *RegistryProperties {
	this := RegistryProperties{}

	this.Location = &location
	this.Name = &name

	return &this
}

// NewRegistryPropertiesWithDefaults instantiates a new RegistryProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegistryPropertiesWithDefaults() *RegistryProperties {
	this := RegistryProperties{}
	return &this
}

// GetGarbageCollectionSchedule returns the GarbageCollectionSchedule field value
// If the value is explicit nil, the zero value for WeeklySchedule will be returned
func (o *RegistryProperties) GetGarbageCollectionSchedule() *WeeklySchedule {
	if o == nil {
		return nil
	}

	return o.GarbageCollectionSchedule

}

// GetGarbageCollectionScheduleOk returns a tuple with the GarbageCollectionSchedule field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetGarbageCollectionScheduleOk() (*WeeklySchedule, bool) {
	if o == nil {
		return nil, false
	}

	return o.GarbageCollectionSchedule, true
}

// SetGarbageCollectionSchedule sets field value
func (o *RegistryProperties) SetGarbageCollectionSchedule(v WeeklySchedule) {

	o.GarbageCollectionSchedule = &v

}

// HasGarbageCollectionSchedule returns a boolean if a field has been set.
func (o *RegistryProperties) HasGarbageCollectionSchedule() bool {
	if o != nil && o.GarbageCollectionSchedule != nil {
		return true
	}

	return false
}

// GetHostname returns the Hostname field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryProperties) GetHostname() *string {
	if o == nil {
		return nil
	}

	return o.Hostname

}

// GetHostnameOk returns a tuple with the Hostname field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetHostnameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Hostname, true
}

// SetHostname sets field value
func (o *RegistryProperties) SetHostname(v string) {

	o.Hostname = &v

}

// HasHostname returns a boolean if a field has been set.
func (o *RegistryProperties) HasHostname() bool {
	if o != nil && o.Hostname != nil {
		return true
	}

	return false
}

// GetLocation returns the Location field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryProperties) GetLocation() *string {
	if o == nil {
		return nil
	}

	return o.Location

}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetLocationOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Location, true
}

// SetLocation sets field value
func (o *RegistryProperties) SetLocation(v string) {

	o.Location = &v

}

// HasLocation returns a boolean if a field has been set.
func (o *RegistryProperties) HasLocation() bool {
	if o != nil && o.Location != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for WeeklySchedule will be returned
func (o *RegistryProperties) GetMaintenanceWindow() *WeeklySchedule {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetMaintenanceWindowOk() (*WeeklySchedule, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *RegistryProperties) SetMaintenanceWindow(v WeeklySchedule) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *RegistryProperties) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *RegistryProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *RegistryProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *RegistryProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetStorageUsage returns the StorageUsage field value
// If the value is explicit nil, the zero value for StorageUsage will be returned
func (o *RegistryProperties) GetStorageUsage() *StorageUsage {
	if o == nil {
		return nil
	}

	return o.StorageUsage

}

// GetStorageUsageOk returns a tuple with the StorageUsage field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *RegistryProperties) GetStorageUsageOk() (*StorageUsage, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageUsage, true
}

// SetStorageUsage sets field value
func (o *RegistryProperties) SetStorageUsage(v StorageUsage) {

	o.StorageUsage = &v

}

// HasStorageUsage returns a boolean if a field has been set.
func (o *RegistryProperties) HasStorageUsage() bool {
	if o != nil && o.StorageUsage != nil {
		return true
	}

	return false
}

func (o RegistryProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["garbageCollectionSchedule"] = o.GarbageCollectionSchedule

	if o.Hostname != nil {
		toSerialize["hostname"] = o.Hostname
	}

	if o.Location != nil {
		toSerialize["location"] = o.Location
	}

	toSerialize["maintenanceWindow"] = o.MaintenanceWindow

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	toSerialize["storageUsage"] = o.StorageUsage

	return json.Marshal(toSerialize)
}

type NullableRegistryProperties struct {
	value *RegistryProperties
	isSet bool
}

func (v NullableRegistryProperties) Get() *RegistryProperties {
	return v.value
}

func (v *NullableRegistryProperties) Set(val *RegistryProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableRegistryProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableRegistryProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegistryProperties(val *RegistryProperties) *NullableRegistryProperties {
	return &NullableRegistryProperties{value: val, isSet: true}
}

func (v NullableRegistryProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegistryProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
