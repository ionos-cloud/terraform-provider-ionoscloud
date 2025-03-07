/*
 * IONOS Cloud - Managed Stackable Data Platform API
 *
 * *Managed Stackable Data Platform* by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable platform.  The Managed Stackable Data Platform by IONOS Cloud can be configured through the IONOS Cloud API in addition or as an alternative to the *Data Center Designer* (DCD).  ## Getting Started  To get your DataPlatformCluster up and running, the following steps needs to be performed.  ### IONOS Cloud Account  The first step is the creation of a IONOS Cloud account if not already existing.  To register a **new account**, visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).  ### Virtual Data Center (VDC)  The Managed Stackable Data Platform needs a virtual data center (VDC) hosting the cluster. This could either be a VDC that already exists, especially if you want to connect the managed data platform to other services already running within your VDC. Otherwise, if you want to place the Managed Stackable Data Platform in a new VDC or you have not yet created a VDC, you need to do so.  A new VDC can be created via the IONOS Cloud API, the IONOS Cloud CLI (`ionosctl`), or the DCD Web interface. For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/basic-tutorials/data-center-basics).  ### Get a authentication token  To interact with this API a user specific authentication token is needed. This token can be generated using the IONOS Cloud CLI the following way:  ``` ionosctl token generate ```  For more information, [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate).  ### Create a new DataPlatformCluster  Before using the Managed Stackable Data Platform, a new DataPlatformCluster must be created.  To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.  The provisioning of the cluster might take some time. To check the current provisioning status, you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID that was presented to you in the response of the create cluster call.  ### Add a DataPlatformNodePool  To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool that is provisioned along with the cluster is reserved for the Stackable operators. You may create further node pools with resources tailored to your use case.  To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post) endpoint.  ### Receive Kubeconfig  Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API. The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.  To protect the deployment of the Stackable distribution, the kubeconfig does not provide you with administration rights for the cluster. What that means is that your actions and deployments are limited to the **default** namespace.  If you still want to group your deployments, you have the option to create subnamespaces within the default namespace. This is made possible by the concept of *hierarchical namespaces* (HNS). You can find more details [here](https://kubernetes.io/blog/2020/08/14/introducing-hierarchical-namespaces/).  The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint using the cluster ID of the created DataPlatformCluster.  ### Create Stackable Services  You can leverage the `kubeconfig.json` file to access the Managed Stackable Data Platform cluster and manage the deployment of [Stackable data apps](https://stackable.tech/en/platform/).  With the Stackable operators, you can deploy the [data apps](https://docs.stackable.tech/home/stable/getting_started.html#_deploying_stackable_services) you want in your Data Platform cluster.  ## Authorization  All endpoints are secured, so only an authenticated user can access them. As Authentication mechanism the default IONOS Cloud authentication mechanism is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).  ### Basic Auth  The basic auth scheme uses the IONOS Cloud user credentials in form of a *Basic Authentication* header accordingly to [RFC 7617](https://datatracker.ietf.org/doc/html/rfc7617).  ### API Key as Bearer Token  The Bearer auth token used at the API Gateway is a user-related token created with the IONOS Cloud CLI (For details, see the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token/generate)). For every request to be authenticated, the token is passed as *Authorization Bearer* header along with the request.  ### Permissions and Access Roles  Currently, an administrator can see and manipulate all resources in a contract. Furthermore, users with the group privilege `Manage Dataplatform` can access the API.  ## Components  The Managed Stackable Data Platform by IONOS Cloud consists of two components. The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.  ### DataPlatformCluster  A DataPlatformCluster is the virtual instance of the customer services and operations running the managed services like Stackable operators. A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer. Therefore, it's possible to integrate the cluster with other resources as VLANs e.g. to shape the data center in the customer's need and integrate the cluster within the topology the customer wants to build.  In addition to the Kubernetes cluster, a small node pool is provided which is exclusively used to run the Stackable operators.  ### DataPlatformNodePool  A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top. All nodes within a node pool are identical in setup. The nodes of a pool are provisioned into virtual data centers at a location of your choice and you can freely specify the properties of all the nodes at once before creation.  Nodes in node pools provisioned by the Managed Stackable Data Platform Cloud API are read-only in the customer's VDC and can only be modified or deleted via the API.  ## References
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CreateClusterProperties struct for CreateClusterProperties
type CreateClusterProperties struct {
	// The name of your cluster. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.
	Name *string `json:"name"`
	// The version of the data platform.
	DataPlatformVersion *string `json:"dataPlatformVersion,omitempty"`
	// The UUID of the virtual data center (VDC) the cluster is provisioned.
	DatacenterId      *string            `json:"datacenterId"`
	MaintenanceWindow *MaintenanceWindow `json:"maintenanceWindow,omitempty"`
	// A list of LANs you want this node pool to be part of.
	Lans *[]Lan `json:"lans,omitempty"`
}

// NewCreateClusterProperties instantiates a new CreateClusterProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCreateClusterProperties(name string, datacenterId string) *CreateClusterProperties {
	this := CreateClusterProperties{}

	this.Name = &name
	this.DatacenterId = &datacenterId

	return &this
}

// NewCreateClusterPropertiesWithDefaults instantiates a new CreateClusterProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCreateClusterPropertiesWithDefaults() *CreateClusterProperties {
	this := CreateClusterProperties{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateClusterProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *CreateClusterProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetDataPlatformVersion returns the DataPlatformVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateClusterProperties) GetDataPlatformVersion() *string {
	if o == nil {
		return nil
	}

	return o.DataPlatformVersion

}

// GetDataPlatformVersionOk returns a tuple with the DataPlatformVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetDataPlatformVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.DataPlatformVersion, true
}

// SetDataPlatformVersion sets field value
func (o *CreateClusterProperties) SetDataPlatformVersion(v string) {

	o.DataPlatformVersion = &v

}

// HasDataPlatformVersion returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasDataPlatformVersion() bool {
	if o != nil && o.DataPlatformVersion != nil {
		return true
	}

	return false
}

// GetDatacenterId returns the DatacenterId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CreateClusterProperties) GetDatacenterId() *string {
	if o == nil {
		return nil
	}

	return o.DatacenterId

}

// GetDatacenterIdOk returns a tuple with the DatacenterId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetDatacenterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.DatacenterId, true
}

// SetDatacenterId sets field value
func (o *CreateClusterProperties) SetDatacenterId(v string) {

	o.DatacenterId = &v

}

// HasDatacenterId returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasDatacenterId() bool {
	if o != nil && o.DatacenterId != nil {
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

// GetLans returns the Lans field value
// If the value is explicit nil, the zero value for []Lan will be returned
func (o *CreateClusterProperties) GetLans() *[]Lan {
	if o == nil {
		return nil
	}

	return o.Lans

}

// GetLansOk returns a tuple with the Lans field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CreateClusterProperties) GetLansOk() (*[]Lan, bool) {
	if o == nil {
		return nil, false
	}

	return o.Lans, true
}

// SetLans sets field value
func (o *CreateClusterProperties) SetLans(v []Lan) {

	o.Lans = &v

}

// HasLans returns a boolean if a field has been set.
func (o *CreateClusterProperties) HasLans() bool {
	if o != nil && o.Lans != nil {
		return true
	}

	return false
}

func (o CreateClusterProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.DataPlatformVersion != nil {
		toSerialize["dataPlatformVersion"] = o.DataPlatformVersion
	}

	if o.DatacenterId != nil {
		toSerialize["datacenterId"] = o.DatacenterId
	}

	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}

	if o.Lans != nil {
		toSerialize["lans"] = o.Lans
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
