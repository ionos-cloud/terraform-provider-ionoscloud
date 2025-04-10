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

// checks if the MaintenanceWindow type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MaintenanceWindow{}

// MaintenanceWindow Starting time of a weekly 4-hour-long window, during which maintenance might occur in the `HH:MM:SS` format.
type MaintenanceWindow struct {
	// Time at which the maintenance should start. Must conform to the `HH:MM:SS` 24-hour format.
	Time         string `json:"time"`
	DayOfTheWeek string `json:"dayOfTheWeek"`
}

// NewMaintenanceWindow instantiates a new MaintenanceWindow object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMaintenanceWindow(time string, dayOfTheWeek string) *MaintenanceWindow {
	this := MaintenanceWindow{}

	this.Time = time
	this.DayOfTheWeek = dayOfTheWeek

	return &this
}

// NewMaintenanceWindowWithDefaults instantiates a new MaintenanceWindow object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMaintenanceWindowWithDefaults() *MaintenanceWindow {
	this := MaintenanceWindow{}
	return &this
}

// GetTime returns the Time field value
func (o *MaintenanceWindow) GetTime() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Time
}

// GetTimeOk returns a tuple with the Time field value
// and a boolean to check if the value has been set.
func (o *MaintenanceWindow) GetTimeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Time, true
}

// SetTime sets field value
func (o *MaintenanceWindow) SetTime(v string) {
	o.Time = v
}

// GetDayOfTheWeek returns the DayOfTheWeek field value
func (o *MaintenanceWindow) GetDayOfTheWeek() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DayOfTheWeek
}

// GetDayOfTheWeekOk returns a tuple with the DayOfTheWeek field value
// and a boolean to check if the value has been set.
func (o *MaintenanceWindow) GetDayOfTheWeekOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DayOfTheWeek, true
}

// SetDayOfTheWeek sets field value
func (o *MaintenanceWindow) SetDayOfTheWeek(v string) {
	o.DayOfTheWeek = v
}

func (o MaintenanceWindow) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MaintenanceWindow) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["time"] = o.Time
	toSerialize["dayOfTheWeek"] = o.DayOfTheWeek
	return toSerialize, nil
}

type NullableMaintenanceWindow struct {
	value *MaintenanceWindow
	isSet bool
}

func (v NullableMaintenanceWindow) Get() *MaintenanceWindow {
	return v.value
}

func (v *NullableMaintenanceWindow) Set(val *MaintenanceWindow) {
	v.value = val
	v.isSet = true
}

func (v NullableMaintenanceWindow) IsSet() bool {
	return v.isSet
}

func (v *NullableMaintenanceWindow) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMaintenanceWindow(val *MaintenanceWindow) *NullableMaintenanceWindow {
	return &NullableMaintenanceWindow{value: val, isSet: true}
}

func (v NullableMaintenanceWindow) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMaintenanceWindow) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
