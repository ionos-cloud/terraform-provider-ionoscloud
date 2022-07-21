/*
 * IONOS Cloud - Managed Data Stack Solution API
 *
 * The Managed Data Stack Solution by IONOS Cloud provides a preconfigured Kubernetes cluster with pre-installed and managed Stackable operators. After the provision of these Stackable operators, the customer can interact with them directly and build his desired application on top of the Stackable Platform.  The Managed Data Stack Solution by IONOS Cloud can be configured through the IONOS Cloud API in addition or as an alternative to the \"Data Center Designer\" (DCD).  ## Getting Started  To get your DataPlatformCluster up and running, the following steps needs to be performed.  ### IONOS Cloud Account  The first step is the creation of a IONOS Cloud account if not already existing.  To register a **new account** visit [cloud.ionos.com](https://cloud.ionos.com/compute/signup).  ### Virtual Datacenter (VDC)  The Managed Data Stack needs a virtual datacenter (VDC) hosting the cluster. This could either be a VDC that already exists, especially if you want to connect the managed DataPlatform to other services already running within your VDC. Otherwise, if you want to place the Managed Data Stack in a new VDC or you have not yet created a VDC, you need to do so.  A new VDC can be created via the IONOS Cloud API, the IONOS-CLI or the DCD Web interface. For more information, see the [official documentation](https://docs.ionos.com/cloud/getting-started/tutorials/data-center-basics)  ### Get a authentication token  To interact with this API a user specific authentication token is needed. This token can be generated using the IONOS-CLI the following way:  ``` ionosctl token generate ```  For more information [see](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate)  ### Create a new DataPlatformCluster  Before using the managed Stackable solution, a new DataPlatformCluster must be created.  To create a cluster, use the [Create DataPlatformCluster](paths./clusters.post) API endpoint.  The provisioning of the cluster might take some time. To check the current provisioning status, you can query the cluster by calling the [Get Endpoint](#/DataPlatformCluster/getCluster) with the cluster ID that was presented to you in the response of the create cluster call.  ### Add a DataPlatformNodePool  To deploy and run a Stackable service, the cluster must have enough computational resources. The node pool that is provisioned along with the cluster is reserved for the Stackable operators. You may create further node pools with resources tailored to your use-case.  To create a new node pool use the [Create DataPlatformNodepool](paths./clusters/{clusterId}/nodepools.post) endpoint.  ### Receive Kubeconfig  Once the DataPlatformCluster is created, the kubeconfig can be accessed by the API. The kubeconfig allows the interaction with the provided cluster as with any regular Kubernetes cluster.  The kubeconfig can be downloaded with the [Get Kubeconfig](paths./clusters/{clusterId}/kubeconfig.get) endpoint using the cluster ID of the created DataPlatformCluster.  ### Create Stackable Service  To create the desired application, the Stackable service needs to be provided, using the received kubeconfig and [deploy a Stackable service](https://docs.stackable.tech/home/getting_started.html#_deploying_stackable_services)  ## Authorization  All endpoints are secured, so only an authenticated user can access them. As Authentication mechanism the default IONOS Cloud authentication mechanism is used. A detailed description can be found [here](https://api.ionos.com/docs/authentication/).  ### Basic-Auth  The basic auth scheme uses the IONOS Cloud user credentials in form of a Basic Authentication Header accordingly to [RFC7617](https://datatracker.ietf.org/doc/html/rfc7617)  ### API-Key as Bearer Token  The Bearer auth token used at the API-Gateway is a user related token created with the IONOS-CLI. (See the [documentation](https://docs.ionos.com/cli-ionosctl/subcommands/authentication/token-generate) for details) For every request to be authenticated, the token is passed as 'Authorization Bearer' header along with the request.  ### Permissions and access roles  Currently, an admin can see and manipulate all resources in a contract. A normal authenticated user can only see and manipulate resources he created.   ## Components  The Managed Data Stack Solution by IONOS Cloud consists of two components. The concept of a DataPlatformClusters and the backing DataPlatformNodePools the cluster is build on.  ### DataPlatformCluster  A DataPlatformCluster is the virtual instance of the customer services and operations running the managed Services like Stackable operators. A DataPlatformCluster is a Kubernetes Cluster in the VDC of the customer. Therefore, it's possible to integrate the cluster with other resources as vLANs e.G. to shape the datacenter in the customer's need and integrate the cluster within the topology the customer wants to build.  In addition to the Kubernetes cluster a small node pool is provided which is exclusively used to run the Stackable operators.  ### DataPlatformNodePool  A DataPlatformNodePool represents the physical machines a DataPlatformCluster is build on top. All nodes within a node pool are identical in setup. The nodes of a pool are provisioned into virtual data centers at a location of your choice and you can freely specify the properties of all the nodes at once before creation.  Nodes in node pools provisioned by the Managed Data Stack Solution Cloud API are readonly in the customer's VDC and can only be modified or deleted via the API.  ### References
 *
 * API version: 0.0.6
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// PatchClusterRequest Data to update selected properties of a DataPlatformCluster
type PatchClusterRequest struct {
	Properties *PatchClusterProperties `json:"properties,omitempty"`
}

// NewPatchClusterRequest instantiates a new PatchClusterRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPatchClusterRequest() *PatchClusterRequest {
	this := PatchClusterRequest{}

	return &this
}

// NewPatchClusterRequestWithDefaults instantiates a new PatchClusterRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPatchClusterRequestWithDefaults() *PatchClusterRequest {
	this := PatchClusterRequest{}
	return &this
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for PatchClusterProperties will be returned
func (o *PatchClusterRequest) GetProperties() *PatchClusterProperties {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PatchClusterRequest) GetPropertiesOk() (*PatchClusterProperties, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *PatchClusterRequest) SetProperties(v PatchClusterProperties) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *PatchClusterRequest) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o PatchClusterRequest) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullablePatchClusterRequest struct {
	value *PatchClusterRequest
	isSet bool
}

func (v NullablePatchClusterRequest) Get() *PatchClusterRequest {
	return v.value
}

func (v *NullablePatchClusterRequest) Set(val *PatchClusterRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePatchClusterRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePatchClusterRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePatchClusterRequest(val *PatchClusterRequest) *NullablePatchClusterRequest {
	return &NullablePatchClusterRequest{value: val, isSet: true}
}

func (v NullablePatchClusterRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePatchClusterRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}