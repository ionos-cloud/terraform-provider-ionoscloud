/*
 * IONOS Cloud - Managed Stackable Data Platform API
 *
 * *Managed Stackable Data Platform* by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable platform.  The Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API in addition or as an alternative to the *Data Center Designer* (DCD).  ## Getting Started  To get your DataPlatformCluster up and running, the following steps needs to be performed.  ### IONOS Cloud Account  The first step is the creation of a IONOS Cloud account if not already existing.  To register a **new account**, visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).  ### Virtual Data Center (VDC)  The Managed Stackable Data Platform needs a virtual data center (VDC) hosting the cluster. This could either be a VDC that already exists, especially if you want to connect the managed data platform to other services already running within your VDC. Otherwise, if you want to place the Managed Stackable Data Platform in a new VDC or you have not yet created a VDC, you need to do so.  A new VDC can be created via the IONOS Cloud API, the IONOS Cloud CLI (`ionosctl`), or the DCD Web interface. For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/basic-tutorials/data-center-basics).  ### Get a authentication token  To interact with this API a user specific authentication token is needed. This token can be generated using the IONOS Cloud CLI the following way:  ``` ionosctl token generate ```  For more information, [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate).  ### Create a new DataPlatformCluster  Before using the Managed Stackable Data Platform, a new DataPlatformCluster must be created.  To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.  The provisioning of the cluster might take some time. To check the current provisioning status, you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID that was presented to you in the response of the create cluster call.  ### Add a DataPlatformNodePool  To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool that is provisioned along with the cluster is reserved for the Stackable operators. You may create further node pools with resources tailored to your use case.  To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post) endpoint.  ### Receive Kubeconfig  Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API. The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.  To protect the deployment of the Stackable distribution, the kubeconfig does not provide you with administration rights for the cluster. What that means is that your actions and deployments are limited to the **default** namespace.  If you still want to group your deployments, you have the option to create subnamespaces within the default namespace. This is made possible by the concept of *hierarchical namespaces* (HNS). You can find more details [here](https://kubernetes.io/blog/2020/08/14/introducing-hierarchical-namespaces/).  The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint using the cluster ID of the created DataPlatformCluster.  ### Create Stackable Services  You can leverage the `kubeconfig.json` file to access the Managed Stackable Data Platform cluster and manage the deployment of [Stackable data apps](https://stackable.tech/en/platform/).  With the Stackable operators, you can deploy the [data apps](https://docs.stackable.tech/home/stable/getting_started.html#_deploying_stackable_services) you want in your Data Platform cluster.  ## Authorization  All endpoints are secured, so only an authenticated user can access them. As Authentication mechanism the default IONOS Cloud authentication mechanism is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).  ### Basic Auth  The basic auth scheme uses the IONOS Cloud user credentials in form of a *Basic Authentication* header accordingly to [RFC 7617](https://datatracker.ietf.org/doc/html/rfc7617).  ### API Key as Bearer Token  The Bearer auth token used at the API Gateway is a user-related token created with the IONOS Cloud CLI (For details, see the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate)). For every request to be authenticated, the token is passed as *Authorization Bearer* header along with the request.  ### Permissions and Access Roles  Currently, an administrator can see and manipulate all resources in a contract. Furthermore, users with the group privilege `Manage Dataplatform` can access the API.  ## Components  The Managed Stackable Data Platform by IONOS Cloud consists of two components. The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.  ### DataPlatformCluster  A DataPlatformCluster is the virtual instance of the customer services and operations running the managed services like Stackable operators. A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer. Therefore, it's possible to integrate the cluster with other resources as VLANs e.g. to shape the data center in the customer's need and integrate the cluster within the topology the customer wants to build.  In addition to the Kubernetes cluster, a small node pool is provided which is exclusively used to run the Stackable operators.  ### DataPlatformNodePool  A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top. All nodes within a node pool are identical in setup. The nodes of a pool are provisioned into virtual data centers at a location of your choice and you can freely specify the properties of all the nodes at once before creation.  Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are read-only in the customer's VDC and can only be modified or deleted via the API.  ## References
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dataplatform

import (
	"encoding/json"
)

// checks if the NodePool type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &NodePool{}

// NodePool A DataPlatformNodePool of a DataPlatformCluster.
type NodePool struct {
	// The name of your node pool. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.
	Name *string `json:"name,omitempty"`
	// The version of the data platform.
	DataPlatformVersion *string `json:"dataPlatformVersion,omitempty"`
	// The UUID of the virtual data center (VDC) the cluster is provisioned.
	DatacenterId *string `json:"datacenterId,omitempty"`
	// The number of nodes that make up the node pool.
	NodeCount *int32 `json:"nodeCount,omitempty"`
	// A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the data center resource.
	CpuFamily *string `json:"cpuFamily,omitempty"`
	// The number of CPU cores per node.
	CoresCount *int32 `json:"coresCount,omitempty"`
	// The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB.
	RamSize          *int32            `json:"ramSize,omitempty"`
	AvailabilityZone *AvailabilityZone `json:"availabilityZone,omitempty"`
	StorageType      *StorageType      `json:"storageType,omitempty"`
	// The size of the volume in GB. The size must be greater than 10 GB.
	StorageSize       *int32             `json:"storageSize,omitempty"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	// Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
	Labels map[string]interface{} `json:"labels,omitempty"`
	// Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/).
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	AutoScaling *AutoScaling           `json:"autoScaling,omitempty"`
}

// NewNodePool instantiates a new NodePool object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNodePool() *NodePool {
	this := NodePool{}

	var cpuFamily string = "AUTO"
	this.CpuFamily = &cpuFamily
	var coresCount int32 = 4
	this.CoresCount = &coresCount
	var ramSize int32 = 4096
	this.RamSize = &ramSize
	var availabilityZone AvailabilityZone = AVAILABILITYZONE_AUTO
	this.AvailabilityZone = &availabilityZone
	var storageType StorageType = STORAGETYPE_SSD
	this.StorageType = &storageType
	var storageSize int32 = 20
	this.StorageSize = &storageSize

	return &this
}

// NewNodePoolWithDefaults instantiates a new NodePool object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNodePoolWithDefaults() *NodePool {
	this := NodePool{}
	var cpuFamily string = "AUTO"
	this.CpuFamily = &cpuFamily
	var coresCount int32 = 4
	this.CoresCount = &coresCount
	var ramSize int32 = 4096
	this.RamSize = &ramSize
	var availabilityZone AvailabilityZone = AVAILABILITYZONE_AUTO
	this.AvailabilityZone = &availabilityZone
	var storageType StorageType = STORAGETYPE_SSD
	this.StorageType = &storageType
	var storageSize int32 = 20
	this.StorageSize = &storageSize
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *NodePool) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *NodePool) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *NodePool) SetName(v string) {
	o.Name = &v
}

// GetDataPlatformVersion returns the DataPlatformVersion field value if set, zero value otherwise.
func (o *NodePool) GetDataPlatformVersion() string {
	if o == nil || IsNil(o.DataPlatformVersion) {
		var ret string
		return ret
	}
	return *o.DataPlatformVersion
}

// GetDataPlatformVersionOk returns a tuple with the DataPlatformVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetDataPlatformVersionOk() (*string, bool) {
	if o == nil || IsNil(o.DataPlatformVersion) {
		return nil, false
	}
	return o.DataPlatformVersion, true
}

// HasDataPlatformVersion returns a boolean if a field has been set.
func (o *NodePool) HasDataPlatformVersion() bool {
	if o != nil && !IsNil(o.DataPlatformVersion) {
		return true
	}

	return false
}

// SetDataPlatformVersion gets a reference to the given string and assigns it to the DataPlatformVersion field.
func (o *NodePool) SetDataPlatformVersion(v string) {
	o.DataPlatformVersion = &v
}

// GetDatacenterId returns the DatacenterId field value if set, zero value otherwise.
func (o *NodePool) GetDatacenterId() string {
	if o == nil || IsNil(o.DatacenterId) {
		var ret string
		return ret
	}
	return *o.DatacenterId
}

// GetDatacenterIdOk returns a tuple with the DatacenterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetDatacenterIdOk() (*string, bool) {
	if o == nil || IsNil(o.DatacenterId) {
		return nil, false
	}
	return o.DatacenterId, true
}

// HasDatacenterId returns a boolean if a field has been set.
func (o *NodePool) HasDatacenterId() bool {
	if o != nil && !IsNil(o.DatacenterId) {
		return true
	}

	return false
}

// SetDatacenterId gets a reference to the given string and assigns it to the DatacenterId field.
func (o *NodePool) SetDatacenterId(v string) {
	o.DatacenterId = &v
}

// GetNodeCount returns the NodeCount field value if set, zero value otherwise.
func (o *NodePool) GetNodeCount() int32 {
	if o == nil || IsNil(o.NodeCount) {
		var ret int32
		return ret
	}
	return *o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetNodeCountOk() (*int32, bool) {
	if o == nil || IsNil(o.NodeCount) {
		return nil, false
	}
	return o.NodeCount, true
}

// HasNodeCount returns a boolean if a field has been set.
func (o *NodePool) HasNodeCount() bool {
	if o != nil && !IsNil(o.NodeCount) {
		return true
	}

	return false
}

// SetNodeCount gets a reference to the given int32 and assigns it to the NodeCount field.
func (o *NodePool) SetNodeCount(v int32) {
	o.NodeCount = &v
}

// GetCpuFamily returns the CpuFamily field value if set, zero value otherwise.
func (o *NodePool) GetCpuFamily() string {
	if o == nil || IsNil(o.CpuFamily) {
		var ret string
		return ret
	}
	return *o.CpuFamily
}

// GetCpuFamilyOk returns a tuple with the CpuFamily field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetCpuFamilyOk() (*string, bool) {
	if o == nil || IsNil(o.CpuFamily) {
		return nil, false
	}
	return o.CpuFamily, true
}

// HasCpuFamily returns a boolean if a field has been set.
func (o *NodePool) HasCpuFamily() bool {
	if o != nil && !IsNil(o.CpuFamily) {
		return true
	}

	return false
}

// SetCpuFamily gets a reference to the given string and assigns it to the CpuFamily field.
func (o *NodePool) SetCpuFamily(v string) {
	o.CpuFamily = &v
}

// GetCoresCount returns the CoresCount field value if set, zero value otherwise.
func (o *NodePool) GetCoresCount() int32 {
	if o == nil || IsNil(o.CoresCount) {
		var ret int32
		return ret
	}
	return *o.CoresCount
}

// GetCoresCountOk returns a tuple with the CoresCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetCoresCountOk() (*int32, bool) {
	if o == nil || IsNil(o.CoresCount) {
		return nil, false
	}
	return o.CoresCount, true
}

// HasCoresCount returns a boolean if a field has been set.
func (o *NodePool) HasCoresCount() bool {
	if o != nil && !IsNil(o.CoresCount) {
		return true
	}

	return false
}

// SetCoresCount gets a reference to the given int32 and assigns it to the CoresCount field.
func (o *NodePool) SetCoresCount(v int32) {
	o.CoresCount = &v
}

// GetRamSize returns the RamSize field value if set, zero value otherwise.
func (o *NodePool) GetRamSize() int32 {
	if o == nil || IsNil(o.RamSize) {
		var ret int32
		return ret
	}
	return *o.RamSize
}

// GetRamSizeOk returns a tuple with the RamSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetRamSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.RamSize) {
		return nil, false
	}
	return o.RamSize, true
}

// HasRamSize returns a boolean if a field has been set.
func (o *NodePool) HasRamSize() bool {
	if o != nil && !IsNil(o.RamSize) {
		return true
	}

	return false
}

// SetRamSize gets a reference to the given int32 and assigns it to the RamSize field.
func (o *NodePool) SetRamSize(v int32) {
	o.RamSize = &v
}

// GetAvailabilityZone returns the AvailabilityZone field value if set, zero value otherwise.
func (o *NodePool) GetAvailabilityZone() AvailabilityZone {
	if o == nil || IsNil(o.AvailabilityZone) {
		var ret AvailabilityZone
		return ret
	}
	return *o.AvailabilityZone
}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetAvailabilityZoneOk() (*AvailabilityZone, bool) {
	if o == nil || IsNil(o.AvailabilityZone) {
		return nil, false
	}
	return o.AvailabilityZone, true
}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *NodePool) HasAvailabilityZone() bool {
	if o != nil && !IsNil(o.AvailabilityZone) {
		return true
	}

	return false
}

// SetAvailabilityZone gets a reference to the given AvailabilityZone and assigns it to the AvailabilityZone field.
func (o *NodePool) SetAvailabilityZone(v AvailabilityZone) {
	o.AvailabilityZone = &v
}

// GetStorageType returns the StorageType field value if set, zero value otherwise.
func (o *NodePool) GetStorageType() StorageType {
	if o == nil || IsNil(o.StorageType) {
		var ret StorageType
		return ret
	}
	return *o.StorageType
}

// GetStorageTypeOk returns a tuple with the StorageType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetStorageTypeOk() (*StorageType, bool) {
	if o == nil || IsNil(o.StorageType) {
		return nil, false
	}
	return o.StorageType, true
}

// HasStorageType returns a boolean if a field has been set.
func (o *NodePool) HasStorageType() bool {
	if o != nil && !IsNil(o.StorageType) {
		return true
	}

	return false
}

// SetStorageType gets a reference to the given StorageType and assigns it to the StorageType field.
func (o *NodePool) SetStorageType(v StorageType) {
	o.StorageType = &v
}

// GetStorageSize returns the StorageSize field value if set, zero value otherwise.
func (o *NodePool) GetStorageSize() int32 {
	if o == nil || IsNil(o.StorageSize) {
		var ret int32
		return ret
	}
	return *o.StorageSize
}

// GetStorageSizeOk returns a tuple with the StorageSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetStorageSizeOk() (*int32, bool) {
	if o == nil || IsNil(o.StorageSize) {
		return nil, false
	}
	return o.StorageSize, true
}

// HasStorageSize returns a boolean if a field has been set.
func (o *NodePool) HasStorageSize() bool {
	if o != nil && !IsNil(o.StorageSize) {
		return true
	}

	return false
}

// SetStorageSize gets a reference to the given int32 and assigns it to the StorageSize field.
func (o *NodePool) SetStorageSize(v int32) {
	o.StorageSize = &v
}

// GetMaintenanceWindow returns the MaintenanceWindow field value if set, zero value otherwise.
func (o *NodePool) GetMaintenanceWindow() MaintenanceWindow {
	if o == nil || IsNil(o.MaintenanceWindow) {
		var ret MaintenanceWindow
		return ret
	}
	return *o.MaintenanceWindow
}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetMaintenanceWindowOk() (*MaintenanceWindow, bool) {
	if o == nil || IsNil(o.MaintenanceWindow) {
		return nil, false
	}
	return o.MaintenanceWindow, true
}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *NodePool) HasMaintenanceWindow() bool {
	if o != nil && !IsNil(o.MaintenanceWindow) {
		return true
	}

	return false
}

// SetMaintenanceWindow gets a reference to the given MaintenanceWindow and assigns it to the MaintenanceWindow field.
func (o *NodePool) SetMaintenanceWindow(v MaintenanceWindow) {
	o.MaintenanceWindow = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *NodePool) GetLabels() map[string]interface{} {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]interface{}
		return ret
	}
	return o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetLabelsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Labels) {
		return map[string]interface{}{}, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *NodePool) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]interface{} and assigns it to the Labels field.
func (o *NodePool) SetLabels(v map[string]interface{}) {
	o.Labels = v
}

// GetAnnotations returns the Annotations field value if set, zero value otherwise.
func (o *NodePool) GetAnnotations() map[string]interface{} {
	if o == nil || IsNil(o.Annotations) {
		var ret map[string]interface{}
		return ret
	}
	return o.Annotations
}

// GetAnnotationsOk returns a tuple with the Annotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetAnnotationsOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Annotations) {
		return map[string]interface{}{}, false
	}
	return o.Annotations, true
}

// HasAnnotations returns a boolean if a field has been set.
func (o *NodePool) HasAnnotations() bool {
	if o != nil && !IsNil(o.Annotations) {
		return true
	}

	return false
}

// SetAnnotations gets a reference to the given map[string]interface{} and assigns it to the Annotations field.
func (o *NodePool) SetAnnotations(v map[string]interface{}) {
	o.Annotations = v
}

// GetAutoScaling returns the AutoScaling field value if set, zero value otherwise.
func (o *NodePool) GetAutoScaling() AutoScaling {
	if o == nil || IsNil(o.AutoScaling) {
		var ret AutoScaling
		return ret
	}
	return *o.AutoScaling
}

// GetAutoScalingOk returns a tuple with the AutoScaling field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *NodePool) GetAutoScalingOk() (*AutoScaling, bool) {
	if o == nil || IsNil(o.AutoScaling) {
		return nil, false
	}
	return o.AutoScaling, true
}

// HasAutoScaling returns a boolean if a field has been set.
func (o *NodePool) HasAutoScaling() bool {
	if o != nil && !IsNil(o.AutoScaling) {
		return true
	}

	return false
}

// SetAutoScaling gets a reference to the given AutoScaling and assigns it to the AutoScaling field.
func (o *NodePool) SetAutoScaling(v AutoScaling) {
	o.AutoScaling = &v
}

func (o NodePool) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o NodePool) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.DataPlatformVersion) {
		toSerialize["dataPlatformVersion"] = o.DataPlatformVersion
	}
	if !IsNil(o.DatacenterId) {
		toSerialize["datacenterId"] = o.DatacenterId
	}
	if !IsNil(o.NodeCount) {
		toSerialize["nodeCount"] = o.NodeCount
	}
	if !IsNil(o.CpuFamily) {
		toSerialize["cpuFamily"] = o.CpuFamily
	}
	if !IsNil(o.CoresCount) {
		toSerialize["coresCount"] = o.CoresCount
	}
	if !IsNil(o.RamSize) {
		toSerialize["ramSize"] = o.RamSize
	}
	if !IsNil(o.AvailabilityZone) {
		toSerialize["availabilityZone"] = o.AvailabilityZone
	}
	if !IsNil(o.StorageType) {
		toSerialize["storageType"] = o.StorageType
	}
	if !IsNil(o.StorageSize) {
		toSerialize["storageSize"] = o.StorageSize
	}
	if !IsNil(o.MaintenanceWindow) {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.Annotations) {
		toSerialize["annotations"] = o.Annotations
	}
	if !IsNil(o.AutoScaling) {
		toSerialize["autoScaling"] = o.AutoScaling
	}
	return toSerialize, nil
}

type NullableNodePool struct {
	value *NodePool
	isSet bool
}

func (v NullableNodePool) Get() *NodePool {
	return v.value
}

func (v *NullableNodePool) Set(val *NodePool) {
	v.value = val
	v.isSet = true
}

func (v NullableNodePool) IsSet() bool {
	return v.isSet
}

func (v *NullableNodePool) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNodePool(val *NodePool) *NullableNodePool {
	return &NullableNodePool{value: val, isSet: true}
}

func (v NullableNodePool) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNodePool) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
