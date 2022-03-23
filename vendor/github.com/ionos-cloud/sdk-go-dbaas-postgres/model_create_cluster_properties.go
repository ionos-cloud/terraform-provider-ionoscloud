/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CreateClusterProperties Properties with all data needed to create a new PostgreSQL cluster.
type CreateClusterProperties struct {
	// The PostgreSQL version of your cluster.
	PostgresVersion *string `json:"postgresVersion"`
	// The total number of instances in the cluster (one master and n-1 standbys).
	Instances *int32 `json:"instances"`
	// The number of CPU cores per instance.
	Cores *int32 `json:"cores"`
	// The amount of memory per instance in megabytes. Has to be a multiple of 1024.
	Ram *int32 `json:"ram"`
	// The amount of storage per instance in megabytes.
	StorageSize    *int32          `json:"storageSize"`
	StorageType    *StorageType    `json:"storageType"`
	Connections    *[]Connection   `json:"connections"`
	Location       *Location       `json:"location"`
	BackupLocation *BackupLocation `json:"backupLocation,omitempty"`
	// The friendly name of your cluster.
	DisplayName         *string               `json:"displayName"`
	MaintenanceWindow   *MaintenanceWindow    `json:"maintenanceWindow,omitempty"`
	Credentials         *DBUser               `json:"credentials"`
	SynchronizationMode *SynchronizationMode  `json:"synchronizationMode"`
	FromBackup          *CreateRestoreRequest `json:"fromBackup,omitempty"`
}

// NewCreateClusterProperties instantiates a new CreateClusterProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateClusterProperties(postgresVersion string, instances int32, cores int32, ram int32, storageSize int32, storageType StorageType, connections []Connection, location Location, displayName string, credentials DBUser, synchronizationMode SynchronizationMode) *CreateClusterProperties {
	this := CreateClusterProperties{}

	this.PostgresVersion = &postgresVersion
	this.Instances = &instances
	this.Cores = &cores
	this.Ram = &ram
	this.StorageSize = &storageSize
	this.StorageType = &storageType
	this.Connections = &connections
	this.Location = &location
	this.DisplayName = &displayName
	this.Credentials = &credentials
	this.SynchronizationMode = &synchronizationMode

	return &this
}

// NewCreateClusterPropertiesWithDefaults instantiates a new CreateClusterProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateClusterPropertiesWithDefaults() *CreateClusterProperties {
	this := CreateClusterProperties{}
	return &this
}

// GetPostgresVersion returns the PostgresVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateClusterProperties) GetPostgresVersion() *string {
	if o == nil {
		return nil
	}

	return o.PostgresVersion

}

// GetPostgresVersionOk returns a tuple with the PostgresVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetPostgresVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.PostgresVersion, true
}

// SetPostgresVersion sets field value
func (o *CreateClusterProperties) SetPostgresVersion(v string) {

	o.PostgresVersion = &v

}

// HasPostgresVersion returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasPostgresVersion() bool {
	if o != nil && o.PostgresVersion != nil {
		return true
	}

	return false
}

// GetInstances returns the Instances field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateClusterProperties) GetInstances() *int32 {
	if o == nil {
		return nil
	}

	return o.Instances

}

// GetInstancesOk returns a tuple with the Instances field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetInstancesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Instances, true
}

// SetInstances sets field value
func (o *CreateClusterProperties) SetInstances(v int32) {

	o.Instances = &v

}

// HasInstances returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasInstances() bool {
	if o != nil && o.Instances != nil {
		return true
	}

	return false
}

// GetCores returns the Cores field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateClusterProperties) GetCores() *int32 {
	if o == nil {
		return nil
	}

	return o.Cores

}

// GetCoresOk returns a tuple with the Cores field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetCoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Cores, true
}

// SetCores sets field value
func (o *CreateClusterProperties) SetCores(v int32) {

	o.Cores = &v

}

// HasCores returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasCores() bool {
	if o != nil && o.Cores != nil {
		return true
	}

	return false
}

// GetRam returns the Ram field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateClusterProperties) GetRam() *int32 {
	if o == nil {
		return nil
	}

	return o.Ram

}

// GetRamOk returns a tuple with the Ram field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetRamOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.Ram, true
}

// SetRam sets field value
func (o *CreateClusterProperties) SetRam(v int32) {

	o.Ram = &v

}

// HasRam returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasRam() bool {
	if o != nil && o.Ram != nil {
		return true
	}

	return false
}

// GetStorageSize returns the StorageSize field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateClusterProperties) GetStorageSize() *int32 {
	if o == nil {
		return nil
	}

	return o.StorageSize

}

// GetStorageSizeOk returns a tuple with the StorageSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetStorageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageSize, true
}

// SetStorageSize sets field value
func (o *CreateClusterProperties) SetStorageSize(v int32) {

	o.StorageSize = &v

}

// HasStorageSize returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasStorageSize() bool {
	if o != nil && o.StorageSize != nil {
		return true
	}

	return false
}

// GetStorageType returns the StorageType field value
// If the value is explicit nil, the zero value for StorageType will be returned
func (o *CreateClusterProperties) GetStorageType() *StorageType {
	if o == nil {
		return nil
	}

	return o.StorageType

}

// GetStorageTypeOk returns a tuple with the StorageType field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageType, true
}

// SetStorageType sets field value
func (o *CreateClusterProperties) SetStorageType(v StorageType) {

	o.StorageType = &v

}

// HasStorageType returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasStorageType() bool {
	if o != nil && o.StorageType != nil {
		return true
	}

	return false
}

// GetConnections returns the Connections field value
// If the value is explicit nil, the zero value for []Connection will be returned
func (o *CreateClusterProperties) GetConnections() *[]Connection {
	if o == nil {
		return nil
	}

	return o.Connections

}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetConnectionsOk() (*[]Connection, bool) {
	if o == nil {
		return nil, false
	}

	return o.Connections, true
}

// SetConnections sets field value
func (o *CreateClusterProperties) SetConnections(v []Connection) {

	o.Connections = &v

}

// HasConnections returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasConnections() bool {
	if o != nil && o.Connections != nil {
		return true
	}

	return false
}

// GetLocation returns the Location field value
// If the value is explicit nil, the zero value for Location will be returned
func (o *CreateClusterProperties) GetLocation() *Location {
	if o == nil {
		return nil
	}

	return o.Location

}

// GetLocationOk returns a tuple with the Location field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetLocationOk() (*Location, bool) {
	if o == nil {
		return nil, false
	}

	return o.Location, true
}

// SetLocation sets field value
func (o *CreateClusterProperties) SetLocation(v Location) {

	o.Location = &v

}

// HasLocation returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasLocation() bool {
	if o != nil && o.Location != nil {
		return true
	}

	return false
}

// GetBackupLocation returns the BackupLocation field value
// If the value is explicit nil, the zero value for BackupLocation will be returned
func (o *CreateClusterProperties) GetBackupLocation() *BackupLocation {
	if o == nil {
		return nil
	}

	return o.BackupLocation

}

// GetBackupLocationOk returns a tuple with the BackupLocation field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetBackupLocationOk() (*BackupLocation, bool) {
	if o == nil {
		return nil, false
	}

	return o.BackupLocation, true
}

// SetBackupLocation sets field value
func (o *CreateClusterProperties) SetBackupLocation(v BackupLocation) {

	o.BackupLocation = &v

}

// HasBackupLocation returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasBackupLocation() bool {
	if o != nil && o.BackupLocation != nil {
		return true
	}

	return false
}

// GetDisplayName returns the DisplayName field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateClusterProperties) GetDisplayName() *string {
	if o == nil {
		return nil
	}

	return o.DisplayName

}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.DisplayName, true
}

// SetDisplayName sets field value
func (o *CreateClusterProperties) SetDisplayName(v string) {

	o.DisplayName = &v

}

// HasDisplayName returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasDisplayName() bool {
	if o != nil && o.DisplayName != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for MaintenanceWindow will be returned
func (o *CreateClusterProperties) GetMaintenanceWindow() *MaintenanceWindow {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *CreateClusterProperties) SetMaintenanceWindow(v MaintenanceWindow) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetCredentials returns the Credentials field value
// If the value is explicit nil, the zero value for DBUser will be returned
func (o *CreateClusterProperties) GetCredentials() *DBUser {
	if o == nil {
		return nil
	}

	return o.Credentials

}

// GetCredentialsOk returns a tuple with the Credentials field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetCredentialsOk() (*DBUser, bool) {
	if o == nil {
		return nil, false
	}

	return o.Credentials, true
}

// SetCredentials sets field value
func (o *CreateClusterProperties) SetCredentials(v DBUser) {

	o.Credentials = &v

}

// HasCredentials returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasCredentials() bool {
	if o != nil && o.Credentials != nil {
		return true
	}

	return false
}

// GetSynchronizationMode returns the SynchronizationMode field value
// If the value is explicit nil, the zero value for SynchronizationMode will be returned
func (o *CreateClusterProperties) GetSynchronizationMode() *SynchronizationMode {
	if o == nil {
		return nil
	}

	return o.SynchronizationMode

}

// GetSynchronizationModeOk returns a tuple with the SynchronizationMode field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetSynchronizationModeOk() (*SynchronizationMode, bool) {
	if o == nil {
		return nil, false
	}

	return o.SynchronizationMode, true
}

// SetSynchronizationMode sets field value
func (o *CreateClusterProperties) SetSynchronizationMode(v SynchronizationMode) {

	o.SynchronizationMode = &v

}

// HasSynchronizationMode returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasSynchronizationMode() bool {
	if o != nil && o.SynchronizationMode != nil {
		return true
	}

	return false
}

// GetFromBackup returns the FromBackup field value
// If the value is explicit nil, the zero value for CreateRestoreRequest will be returned
func (o *CreateClusterProperties) GetFromBackup() *CreateRestoreRequest {
	if o == nil {
		return nil
	}

	return o.FromBackup

}

// GetFromBackupOk returns a tuple with the FromBackup field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetFromBackupOk() (*CreateRestoreRequest, bool) {
	if o == nil {
		return nil, false
	}

	return o.FromBackup, true
}

// SetFromBackup sets field value
func (o *CreateClusterProperties) SetFromBackup(v CreateRestoreRequest) {

	o.FromBackup = &v

}

// HasFromBackup returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasFromBackup() bool {
	if o != nil && o.FromBackup != nil {
		return true
	}

	return false
}

func (o CreateClusterProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.PostgresVersion != nil {
		toSerialize["postgresVersion"] = o.PostgresVersion
	}

	if o.Instances != nil {
		toSerialize["instances"] = o.Instances
	}

	if o.Cores != nil {
		toSerialize["cores"] = o.Cores
	}

	if o.Ram != nil {
		toSerialize["ram"] = o.Ram
	}

	if o.StorageSize != nil {
		toSerialize["storageSize"] = o.StorageSize
	}

	if o.StorageType != nil {
		toSerialize["storageType"] = o.StorageType
	}

	if o.Connections != nil {
		toSerialize["connections"] = o.Connections
	}

	if o.Location != nil {
		toSerialize["location"] = o.Location
	}

	if o.BackupLocation != nil {
		toSerialize["backupLocation"] = o.BackupLocation
	}

	if o.DisplayName != nil {
		toSerialize["displayName"] = o.DisplayName
	}

	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}

	if o.Credentials != nil {
		toSerialize["credentials"] = o.Credentials
	}

	if o.SynchronizationMode != nil {
		toSerialize["synchronizationMode"] = o.SynchronizationMode
	}

	if o.FromBackup != nil {
		toSerialize["fromBackup"] = o.FromBackup
	}

	return json.Marshal(toSerialize)
}

type NullableCreateClusterProperties struct {
	value *CreateClusterProperties
	isSet bool
}

func (v NullableCreateClusterProperties) Get() *CreateClusterProperties {
	return v.value
}

func (v *NullableCreateClusterProperties) Set(val *CreateClusterProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateClusterProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateClusterProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateClusterProperties(val *CreateClusterProperties) *NullableCreateClusterProperties {
	return &NullableCreateClusterProperties{value: val, isSet: true}
}

func (v NullableCreateClusterProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateClusterProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
