/*
 * IONOS DBaaS MariaDB REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional MariaDB database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package mariadb

import (
	"encoding/json"
)

// checks if the CreateClusterProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CreateClusterProperties{}

// CreateClusterProperties Properties with all data needed to create a new MariaDB cluster.
type CreateClusterProperties struct {
	MariadbVersion MariadbVersion `json:"mariadbVersion"`
	// The total number of instances in the cluster (one primary and n-1 secondary).
	Instances int32 `json:"instances"`
	// The number of CPU cores per instance.
	Cores int32 `json:"cores"`
	// The amount of memory per instance in gigabytes (GB).
	Ram int32 `json:"ram"`
	// The amount of storage per instance in gigabytes (GB).
	StorageSize int32 `json:"storageSize"`
	// The network connection for your cluster. Only one connection is allowed.
	Connections []Connection `json:"connections"`
	// The friendly name of your cluster.
	DisplayName       string             `json:"displayName"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	Backup            *BackupProperties  `json:"backup,omitempty"`
	Credentials       DBUser             `json:"credentials"`
	FromBackup        *RestoreRequest    `json:"fromBackup,omitempty"`
}

// NewCreateClusterProperties instantiates a new CreateClusterProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateClusterProperties(mariadbVersion MariadbVersion, instances int32, cores int32, ram int32, storageSize int32, connections []Connection, displayName string, credentials DBUser) *CreateClusterProperties {
	this := CreateClusterProperties{}

	this.MariadbVersion = mariadbVersion
	this.Instances = instances
	this.Cores = cores
	this.Ram = ram
	this.StorageSize = storageSize
	this.Connections = connections
	this.DisplayName = displayName
	this.Credentials = credentials

	return &this
}

// NewCreateClusterPropertiesWithDefaults instantiates a new CreateClusterProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateClusterPropertiesWithDefaults() *CreateClusterProperties {
	this := CreateClusterProperties{}
	return &this
}

// GetMariadbVersion returns the MariadbVersion field value
func (o *CreateClusterProperties) GetMariadbVersion() MariadbVersion {
	if o == nil {
		var ret MariadbVersion
		return ret
	}

	return o.MariadbVersion
}

// GetMariadbVersionOk returns a tuple with the MariadbVersion field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetMariadbVersionOk() (*MariadbVersion, bool) {
	if o == nil {
		return nil, false
	}
	return &o.MariadbVersion, true
}

// SetMariadbVersion sets field value
func (o *CreateClusterProperties) SetMariadbVersion(v MariadbVersion) {
	o.MariadbVersion = v
}

// GetInstances returns the Instances field value
func (o *CreateClusterProperties) GetInstances() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Instances
}

// GetInstancesOk returns a tuple with the Instances field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetInstancesOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Instances, true
}

// SetInstances sets field value
func (o *CreateClusterProperties) SetInstances(v int32) {
	o.Instances = v
}

// GetCores returns the Cores field value
func (o *CreateClusterProperties) GetCores() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Cores
}

// GetCoresOk returns a tuple with the Cores field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetCoresOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Cores, true
}

// SetCores sets field value
func (o *CreateClusterProperties) SetCores(v int32) {
	o.Cores = v
}

// GetRam returns the Ram field value
func (o *CreateClusterProperties) GetRam() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Ram
}

// GetRamOk returns a tuple with the Ram field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetRamOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ram, true
}

// SetRam sets field value
func (o *CreateClusterProperties) SetRam(v int32) {
	o.Ram = v
}

// GetStorageSize returns the StorageSize field value
func (o *CreateClusterProperties) GetStorageSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.StorageSize
}

// GetStorageSizeOk returns a tuple with the StorageSize field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetStorageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StorageSize, true
}

// SetStorageSize sets field value
func (o *CreateClusterProperties) SetStorageSize(v int32) {
	o.StorageSize = v
}

// GetConnections returns the Connections field value
func (o *CreateClusterProperties) GetConnections() []Connection {
	if o == nil {
		var ret []Connection
		return ret
	}

	return o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetConnectionsOk() ([]Connection, bool) {
	if o == nil {
		return nil, false
	}
	return o.Connections, true
}

// SetConnections sets field value
func (o *CreateClusterProperties) SetConnections(v []Connection) {
	o.Connections = v
}

// GetDisplayName returns the DisplayName field value
func (o *CreateClusterProperties) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *CreateClusterProperties) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetMaintenanceWindow returns the MaintenanceWindow field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetMaintenanceWindow() MaintenanceWindow {
	if o == nil || IsNil(o.MaintenanceWindow) {
		var ret MaintenanceWindow
		return ret
	}
	return *o.MaintenanceWindow
}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil || IsNil(o.MaintenanceWindow) {
		return nil, false
	}
	return o.MaintenanceWindow, true
}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasMaintenanceWindow() bool {
	if o != nil && !IsNil(o.MaintenanceWindow) {
		return true
	}

	return false
}

// SetMaintenanceWindow gets a reference to the given MaintenanceWindow and assigns it to the MaintenanceWindow field.
func (o *CreateClusterProperties) SetMaintenanceWindow(v MaintenanceWindow) {
	o.MaintenanceWindow = &v
}

// GetBackup returns the Backup field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetBackup() BackupProperties {
	if o == nil || IsNil(o.Backup) {
		var ret BackupProperties
		return ret
	}
	return *o.Backup
}

// GetBackupOk returns a tuple with the Backup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetBackupOk() (*BackupProperties, bool) {
	if o == nil || IsNil(o.Backup) {
		return nil, false
	}
	return o.Backup, true
}

// HasBackup returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasBackup() bool {
	if o != nil && !IsNil(o.Backup) {
		return true
	}

	return false
}

// SetBackup gets a reference to the given BackupProperties and assigns it to the Backup field.
func (o *CreateClusterProperties) SetBackup(v BackupProperties) {
	o.Backup = &v
}

// GetCredentials returns the Credentials field value
func (o *CreateClusterProperties) GetCredentials() DBUser {
	if o == nil {
		var ret DBUser
		return ret
	}

	return o.Credentials
}

// GetCredentialsOk returns a tuple with the Credentials field value
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetCredentialsOk() (*DBUser, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Credentials, true
}

// SetCredentials sets field value
func (o *CreateClusterProperties) SetCredentials(v DBUser) {
	o.Credentials = v
}

// GetFromBackup returns the FromBackup field value if set, zero value otherwise.
func (o *CreateClusterProperties) GetFromBackup() RestoreRequest {
	if o == nil || IsNil(o.FromBackup) {
		var ret RestoreRequest
		return ret
	}
	return *o.FromBackup
}

// GetFromBackupOk returns a tuple with the FromBackup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CreateClusterProperties) GetFromBackupOk() (*RestoreRequest, bool) {
	if o == nil || IsNil(o.FromBackup) {
		return nil, false
	}
	return o.FromBackup, true
}

// HasFromBackup returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasFromBackup() bool {
	if o != nil && !IsNil(o.FromBackup) {
		return true
	}

	return false
}

// SetFromBackup gets a reference to the given RestoreRequest and assigns it to the FromBackup field.
func (o *CreateClusterProperties) SetFromBackup(v RestoreRequest) {
	o.FromBackup = &v
}

func (o CreateClusterProperties) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CreateClusterProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["mariadbVersion"] = o.MariadbVersion
	toSerialize["instances"] = o.Instances
	toSerialize["cores"] = o.Cores
	toSerialize["ram"] = o.Ram
	toSerialize["storageSize"] = o.StorageSize
	toSerialize["connections"] = o.Connections
	toSerialize["displayName"] = o.DisplayName
	if !IsNil(o.MaintenanceWindow) {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}
	if !IsNil(o.Backup) {
		toSerialize["backup"] = o.Backup
	}
	toSerialize["credentials"] = o.Credentials
	if !IsNil(o.FromBackup) {
		toSerialize["fromBackup"] = o.FromBackup
	}
	return toSerialize, nil
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
