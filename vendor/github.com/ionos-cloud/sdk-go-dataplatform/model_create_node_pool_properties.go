/*
 * IONOS Cloud - Managed Stackable Data Platform API
 *
 * Managed Stackable Data Platform by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable Platform.  Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API in addition or as an alternative to the \"Data Center Designer\" (DCD).  ## Getting Started  To get your DataPlatformCluster up and running, the following steps needs to be performed.  ### IONOS Cloud Account  The first step is the creation of a IONOS Cloud account if not already existing.  To register a **new account** visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).  ### Virtual Datacenter (VDC)  The Managed Data Stack needs a virtual datacenter (VDC) hosting the cluster. This could either be a VDC that already exists, especially if you want to connect the managed DataPlatform to other services already running within your VDC. Otherwise, if you want to place the Managed Data Stack in a new VDC or you have not yet created a VDC, you need to do so.  A new VDC can be created via the IONOS Cloud API, the IONOS-CLI or the DCD Web interface. For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/tutorials/data-center-basics)  ### Get a authentication token  To interact with this API a user specific authentication token is needed. This token can be generated using the IONOS-CLI the following way:  ``` ionosctl token generate ```  For more information [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate)  ### Create a new DataPlatformCluster  Before using Managed Stackable Data Platform, a new DataPlatformCluster must be created.  To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.  The provisioning of the cluster might take some time. To check the current provisioning status, you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID that was presented to you in the response of the create cluster call.  ### Add a DataPlatformNodePool  To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool that is provisioned along with the cluster is reserved for the Stackable operators. You may create further node pools with resources tailored to your use-case.  To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post) endpoint.  ### Receive Kubeconfig  Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API. The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.  The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint using the cluster ID of the created DataPlatformCluster.  ### Create Stackable Service  To create the desired application, the Stackable service needs to be provided, using the received kubeconfig and [deploy a Stackable service](https://docs.stackable.tech/home/getting_started.html#_deploying_stackable_services)  ## Authorization  All endpoints are secured, so only an authenticated user can access them. As Authentication mechanism the default IONOS Cloud authentication mechanism is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).  ### Basic-Auth  The basic auth scheme uses the IONOS Cloud user credentials in form of a Basic Authentication Header accordingly to [RFC7617](https://datatracker.ietf.org/doc/html/rfc7617)  ### API-Key as Bearer Token  The Bearer auth token used at the API-Gateway is a user related token created with the IONOS-CLI. (See the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate) for details) For every request to be authenticated, the token is passed as 'Authorization Bearer' header along with the request.  ### Permissions and access roles  Currently, an admin can see and manipulate all resources in a contract. A normal authenticated user can only see and manipulate resources he created.   ## Components  The Managed Stackable Data Platform by IONOS Cloud consists of two components. The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.  ### DataPlatformCluster  A DataPlatformCluster is the virtual instance of the customer services and operations running the managed Services like Stackable operators. A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer. Therefore, it's possible to integrate the cluster with other resources as vLANs e.G. to shape the datacenter in the customer's need and integrate the cluster within the topology the customer wants to build.  In addition to the Kubernetes cluster a small node pool is provided which is exclusively used to run the Stackable operators.  ### DataPlatformNodePool  A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top. All nodes within a node pool are identical in setup. The nodes of a pool are provisioned into virtual data centers at a location of your choice and you can freely specify the properties of all the nodes at once before creation.  Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are readonly in the customer's VDC and can only be modified or deleted via the API.  ### References
 *
 * API version: 0.0.7
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CreateNodePoolProperties struct for CreateNodePoolProperties
type CreateNodePoolProperties struct {
	// The name of your node pool. Must be 63 characters or less and must begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
	Name *string `json:"name"`
	// The number of nodes that make up the node pool.
	NodeCount *int32 `json:"nodeCount"`
	// A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource.
	CpuFamily *string `json:"cpuFamily,omitempty"`
	// The number of CPU cores per node.
	CoresCount *int32 `json:"coresCount,omitempty"`
	// The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.
	RamSize          *int32            `json:"ramSize,omitempty"`
	AvailabilityZone *AvailabilityZone `json:"availabilityZone,omitempty"`
	StorageType      *StorageType      `json:"storageType,omitempty"`
	// The size of the volume in GB. The size must be greater than 10GB.
	StorageSize       *int32             `json:"storageSize,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	// Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)
	Labels *map[string]interface{} `json:"labels,omitempty"`
	// Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/)
	Annotations *map[string]interface{} `json:"annotations,omitempty"`
}

// NewCreateNodePoolProperties instantiates a new CreateNodePoolProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateNodePoolProperties(name string, nodeCount int32) *CreateNodePoolProperties {
	this := CreateNodePoolProperties{}

	this.Name = &name
	this.NodeCount = &nodeCount
	var cpuFamily string = "AUTO"
	this.CpuFamily = &cpuFamily
	var coresCount int32 = 4
	this.CoresCount = &coresCount
	var ramSize int32 = 4096
	this.RamSize = &ramSize
	var availabilityZone AvailabilityZone = AUTO
	this.AvailabilityZone = &availabilityZone
	var storageType StorageType = SSD
	this.StorageType = &storageType
	var storageSize int32 = 20
	this.StorageSize = &storageSize

	return &this
}

// NewCreateNodePoolPropertiesWithDefaults instantiates a new CreateNodePoolProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateNodePoolPropertiesWithDefaults() *CreateNodePoolProperties {
	this := CreateNodePoolProperties{}
	var cpuFamily string = "AUTO"
	this.CpuFamily = &cpuFamily
	var coresCount int32 = 4
	this.CoresCount = &coresCount
	var ramSize int32 = 4096
	this.RamSize = &ramSize
	var availabilityZone AvailabilityZone = AUTO
	this.AvailabilityZone = &availabilityZone
	var storageType StorageType = SSD
	this.StorageType = &storageType
	var storageSize int32 = 20
	this.StorageSize = &storageSize
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateNodePoolProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *CreateNodePoolProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetNodeCount returns the NodeCount field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateNodePoolProperties) GetNodeCount() *int32 {
	if o == nil {
		return nil
	}

	return o.NodeCount

}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.NodeCount, true
}

// SetNodeCount sets field value
func (o *CreateNodePoolProperties) SetNodeCount(v int32) {

	o.NodeCount = &v

}

// HasNodeCount returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasNodeCount() bool {
	if o != nil && o.NodeCount != nil {
		return true
	}

	return false
}

// GetCpuFamily returns the CpuFamily field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateNodePoolProperties) GetCpuFamily() *string {
	if o == nil {
		return nil
	}

	return o.CpuFamily

}

// GetCpuFamilyOk returns a tuple with the CpuFamily field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetCpuFamilyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CpuFamily, true
}

// SetCpuFamily sets field value
func (o *CreateNodePoolProperties) SetCpuFamily(v string) {

	o.CpuFamily = &v

}

// HasCpuFamily returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasCpuFamily() bool {
	if o != nil && o.CpuFamily != nil {
		return true
	}

	return false
}

// GetCoresCount returns the CoresCount field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateNodePoolProperties) GetCoresCount() *int32 {
	if o == nil {
		return nil
	}

	return o.CoresCount

}

// GetCoresCountOk returns a tuple with the CoresCount field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetCoresCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.CoresCount, true
}

// SetCoresCount sets field value
func (o *CreateNodePoolProperties) SetCoresCount(v int32) {

	o.CoresCount = &v

}

// HasCoresCount returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasCoresCount() bool {
	if o != nil && o.CoresCount != nil {
		return true
	}

	return false
}

// GetRamSize returns the RamSize field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateNodePoolProperties) GetRamSize() *int32 {
	if o == nil {
		return nil
	}

	return o.RamSize

}

// GetRamSizeOk returns a tuple with the RamSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetRamSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.RamSize, true
}

// SetRamSize sets field value
func (o *CreateNodePoolProperties) SetRamSize(v int32) {

	o.RamSize = &v

}

// HasRamSize returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasRamSize() bool {
	if o != nil && o.RamSize != nil {
		return true
	}

	return false
}

// GetAvailabilityZone returns the AvailabilityZone field value
// If the value is explicit nil, the zero value for AvailabilityZone will be returned
func (o *CreateNodePoolProperties) GetAvailabilityZone() *AvailabilityZone {
	if o == nil {
		return nil
	}

	return o.AvailabilityZone

}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetAvailabilityZoneOk() (*AvailabilityZone, bool) {
	if o == nil {
		return nil, false
	}

	return o.AvailabilityZone, true
}

// SetAvailabilityZone sets field value
func (o *CreateNodePoolProperties) SetAvailabilityZone(v AvailabilityZone) {

	o.AvailabilityZone = &v

}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasAvailabilityZone() bool {
	if o != nil && o.AvailabilityZone != nil {
		return true
	}

	return false
}

// GetStorageType returns the StorageType field value
// If the value is explicit nil, the zero value for StorageType will be returned
func (o *CreateNodePoolProperties) GetStorageType() *StorageType {
	if o == nil {
		return nil
	}

	return o.StorageType

}

// GetStorageTypeOk returns a tuple with the StorageType field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageType, true
}

// SetStorageType sets field value
func (o *CreateNodePoolProperties) SetStorageType(v StorageType) {

	o.StorageType = &v

}

// HasStorageType returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasStorageType() bool {
	if o != nil && o.StorageType != nil {
		return true
	}

	return false
}

// GetStorageSize returns the StorageSize field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *CreateNodePoolProperties) GetStorageSize() *int32 {
	if o == nil {
		return nil
	}

	return o.StorageSize

}

// GetStorageSizeOk returns a tuple with the StorageSize field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetStorageSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.StorageSize, true
}

// SetStorageSize sets field value
func (o *CreateNodePoolProperties) SetStorageSize(v int32) {

	o.StorageSize = &v

}

// HasStorageSize returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasStorageSize() bool {
	if o != nil && o.StorageSize != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for MaintenanceWindow will be returned
func (o *CreateNodePoolProperties) GetMaintenanceWindow() *MaintenanceWindow {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *CreateNodePoolProperties) SetMaintenanceWindow(v MaintenanceWindow) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetLabels returns the Labels field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *CreateNodePoolProperties) GetLabels() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Labels

}

// GetLabelsOk returns a tuple with the Labels field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetLabelsOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Labels, true
}

// SetLabels sets field value
func (o *CreateNodePoolProperties) SetLabels(v map[string]interface{}) {

	o.Labels = &v

}

// HasLabels returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasLabels() bool {
	if o != nil && o.Labels != nil {
		return true
	}

	return false
}

// GetAnnotations returns the Annotations field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *CreateNodePoolProperties) GetAnnotations() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Annotations

}

// GetAnnotationsOk returns a tuple with the Annotations field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateNodePoolProperties) GetAnnotationsOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Annotations, true
}

// SetAnnotations sets field value
func (o *CreateNodePoolProperties) SetAnnotations(v map[string]interface{}) {

	o.Annotations = &v

}

// HasAnnotations returns a boolean if a field has been set.
func (o *CreateNodePoolProperties) HasAnnotations() bool {
	if o != nil && o.Annotations != nil {
		return true
	}

	return false
}

func (o CreateNodePoolProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.NodeCount != nil {
		toSerialize["nodeCount"] = o.NodeCount
	}

	if o.CpuFamily != nil {
		toSerialize["cpuFamily"] = o.CpuFamily
	}

	if o.CoresCount != nil {
		toSerialize["coresCount"] = o.CoresCount
	}

	if o.RamSize != nil {
		toSerialize["ramSize"] = o.RamSize
	}

	if o.AvailabilityZone != nil {
		toSerialize["availabilityZone"] = o.AvailabilityZone
	}

	if o.StorageType != nil {
		toSerialize["storageType"] = o.StorageType
	}

	if o.StorageSize != nil {
		toSerialize["storageSize"] = o.StorageSize
	}

	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}

	if o.Labels != nil {
		toSerialize["labels"] = o.Labels
	}

	if o.Annotations != nil {
		toSerialize["annotations"] = o.Annotations
	}

	return json.Marshal(toSerialize)
}

type NullableCreateNodePoolProperties struct {
	value *CreateNodePoolProperties
	isSet bool
}

func (v NullableCreateNodePoolProperties) Get() *CreateNodePoolProperties {
	return v.value
}

func (v *NullableCreateNodePoolProperties) Set(val *CreateNodePoolProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableCreateNodePoolProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableCreateNodePoolProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCreateNodePoolProperties(val *CreateNodePoolProperties) *NullableCreateNodePoolProperties {
	return &NullableCreateNodePoolProperties{value: val, isSet: true}
}

func (v NullableCreateNodePoolProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCreateNodePoolProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
