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

	"time"
)

// checks if the Metadata type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Metadata{}

// Metadata Metadata of the resource.
type Metadata struct {
	// The *entity tag* of the resource as defined in [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html#name-etag).
	ETag *string `json:"ETag,omitempty"`
	// The time the resource was created, ISO 8601 timestamp (UTC).
	CreatedDate *IonosTime `json:"createdDate,omitempty"`
	// The user that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`
	// The ID of the user that created the resource.
	CreatedByUserId *string `json:"createdByUserId,omitempty"`
	// The creators' contract number.
	CreatedInContractNumber *string `json:"createdInContractNumber,omitempty"`
	// The last time the resource was modified, ISO 8601 timestamp (UTC).
	LastModifiedDate *IonosTime `json:"lastModifiedDate,omitempty"`
	// The user that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// The ID of the user that last modified the resource.
	LastModifiedByUserId *string `json:"lastModifiedByUserId,omitempty"`
	// The version of the data platform.
	CurrentDataPlatformVersion *string `json:"currentDataPlatformVersion,omitempty"`
	// The current data platform revision of a resource. This internal revision is used to roll out non-breaking internal changes. This attribute is read-only.
	CurrentDataPlatformRevision *int64 `json:"currentDataPlatformRevision,omitempty"`
	// List of available upgrades for this cluster.
	AvailableUpgradeVersions []string `json:"availableUpgradeVersions,omitempty"`
	// State of the resource. Resource states: `AVAILABLE`: There are no pending modification requests for this item. `BUSY`: There is at least one modification request pending and all following requests will be queued. `DEPLOYING`: The resource is being created. `FAILED`: The creation of the resource failed. `UPDATING`: The resource is being updated. `FAILED_UPDATING`: An update to the resource was not successful. `DESTROYING`: The resource is being deleted. `FAILED_DESTROYING`: The deletion of the resource was not successful. `TERMINATED`: The resource has been deleted.
	State *string `json:"state,omitempty"`
}

// NewMetadata instantiates a new Metadata object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadata() *Metadata {
	this := Metadata{}

	return &this
}

// NewMetadataWithDefaults instantiates a new Metadata object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithDefaults() *Metadata {
	this := Metadata{}
	return &this
}

// GetETag returns the ETag field value if set, zero value otherwise.
func (o *Metadata) GetETag() string {
	if o == nil || IsNil(o.ETag) {
		var ret string
		return ret
	}
	return *o.ETag
}

// GetETagOk returns a tuple with the ETag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetETagOk() (*string, bool) {
	if o == nil || IsNil(o.ETag) {
		return nil, false
	}
	return o.ETag, true
}

// HasETag returns a boolean if a field has been set.
func (o *Metadata) HasETag() bool {
	if o != nil && !IsNil(o.ETag) {
		return true
	}

	return false
}

// SetETag gets a reference to the given string and assigns it to the ETag field.
func (o *Metadata) SetETag(v string) {
	o.ETag = &v
}

// GetCreatedDate returns the CreatedDate field value if set, zero value otherwise.
func (o *Metadata) GetCreatedDate() time.Time {
	if o == nil || IsNil(o.CreatedDate) {
		var ret time.Time
		return ret
	}
	return o.CreatedDate.Time
}

// GetCreatedDateOk returns a tuple with the CreatedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreatedDate) {
		return nil, false
	}
	return &o.CreatedDate.Time, true
}

// HasCreatedDate returns a boolean if a field has been set.
func (o *Metadata) HasCreatedDate() bool {
	if o != nil && !IsNil(o.CreatedDate) {
		return true
	}

	return false
}

// SetCreatedDate gets a reference to the given time.Time and assigns it to the CreatedDate field.
func (o *Metadata) SetCreatedDate(v time.Time) {
	o.CreatedDate = &IonosTime{v}
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *Metadata) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *Metadata) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *Metadata) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetCreatedByUserId returns the CreatedByUserId field value if set, zero value otherwise.
func (o *Metadata) GetCreatedByUserId() string {
	if o == nil || IsNil(o.CreatedByUserId) {
		var ret string
		return ret
	}
	return *o.CreatedByUserId
}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedByUserId) {
		return nil, false
	}
	return o.CreatedByUserId, true
}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasCreatedByUserId() bool {
	if o != nil && !IsNil(o.CreatedByUserId) {
		return true
	}

	return false
}

// SetCreatedByUserId gets a reference to the given string and assigns it to the CreatedByUserId field.
func (o *Metadata) SetCreatedByUserId(v string) {
	o.CreatedByUserId = &v
}

// GetCreatedInContractNumber returns the CreatedInContractNumber field value if set, zero value otherwise.
func (o *Metadata) GetCreatedInContractNumber() string {
	if o == nil || IsNil(o.CreatedInContractNumber) {
		var ret string
		return ret
	}
	return *o.CreatedInContractNumber
}

// GetCreatedInContractNumberOk returns a tuple with the CreatedInContractNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCreatedInContractNumberOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedInContractNumber) {
		return nil, false
	}
	return o.CreatedInContractNumber, true
}

// HasCreatedInContractNumber returns a boolean if a field has been set.
func (o *Metadata) HasCreatedInContractNumber() bool {
	if o != nil && !IsNil(o.CreatedInContractNumber) {
		return true
	}

	return false
}

// SetCreatedInContractNumber gets a reference to the given string and assigns it to the CreatedInContractNumber field.
func (o *Metadata) SetCreatedInContractNumber(v string) {
	o.CreatedInContractNumber = &v
}

// GetLastModifiedDate returns the LastModifiedDate field value if set, zero value otherwise.
func (o *Metadata) GetLastModifiedDate() time.Time {
	if o == nil || IsNil(o.LastModifiedDate) {
		var ret time.Time
		return ret
	}
	return o.LastModifiedDate.Time
}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil || IsNil(o.LastModifiedDate) {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true
}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedDate() bool {
	if o != nil && !IsNil(o.LastModifiedDate) {
		return true
	}

	return false
}

// SetLastModifiedDate gets a reference to the given time.Time and assigns it to the LastModifiedDate field.
func (o *Metadata) SetLastModifiedDate(v time.Time) {
	o.LastModifiedDate = &IonosTime{v}
}

// GetLastModifiedBy returns the LastModifiedBy field value if set, zero value otherwise.
func (o *Metadata) GetLastModifiedBy() string {
	if o == nil || IsNil(o.LastModifiedBy) {
		var ret string
		return ret
	}
	return *o.LastModifiedBy
}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetLastModifiedByOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedBy) {
		return nil, false
	}
	return o.LastModifiedBy, true
}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedBy() bool {
	if o != nil && !IsNil(o.LastModifiedBy) {
		return true
	}

	return false
}

// SetLastModifiedBy gets a reference to the given string and assigns it to the LastModifiedBy field.
func (o *Metadata) SetLastModifiedBy(v string) {
	o.LastModifiedBy = &v
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value if set, zero value otherwise.
func (o *Metadata) GetLastModifiedByUserId() string {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		var ret string
		return ret
	}
	return *o.LastModifiedByUserId
}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil || IsNil(o.LastModifiedByUserId) {
		return nil, false
	}
	return o.LastModifiedByUserId, true
}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedByUserId() bool {
	if o != nil && !IsNil(o.LastModifiedByUserId) {
		return true
	}

	return false
}

// SetLastModifiedByUserId gets a reference to the given string and assigns it to the LastModifiedByUserId field.
func (o *Metadata) SetLastModifiedByUserId(v string) {
	o.LastModifiedByUserId = &v
}

// GetCurrentDataPlatformVersion returns the CurrentDataPlatformVersion field value if set, zero value otherwise.
func (o *Metadata) GetCurrentDataPlatformVersion() string {
	if o == nil || IsNil(o.CurrentDataPlatformVersion) {
		var ret string
		return ret
	}
	return *o.CurrentDataPlatformVersion
}

// GetCurrentDataPlatformVersionOk returns a tuple with the CurrentDataPlatformVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCurrentDataPlatformVersionOk() (*string, bool) {
	if o == nil || IsNil(o.CurrentDataPlatformVersion) {
		return nil, false
	}
	return o.CurrentDataPlatformVersion, true
}

// HasCurrentDataPlatformVersion returns a boolean if a field has been set.
func (o *Metadata) HasCurrentDataPlatformVersion() bool {
	if o != nil && !IsNil(o.CurrentDataPlatformVersion) {
		return true
	}

	return false
}

// SetCurrentDataPlatformVersion gets a reference to the given string and assigns it to the CurrentDataPlatformVersion field.
func (o *Metadata) SetCurrentDataPlatformVersion(v string) {
	o.CurrentDataPlatformVersion = &v
}

// GetCurrentDataPlatformRevision returns the CurrentDataPlatformRevision field value if set, zero value otherwise.
func (o *Metadata) GetCurrentDataPlatformRevision() int64 {
	if o == nil || IsNil(o.CurrentDataPlatformRevision) {
		var ret int64
		return ret
	}
	return *o.CurrentDataPlatformRevision
}

// GetCurrentDataPlatformRevisionOk returns a tuple with the CurrentDataPlatformRevision field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetCurrentDataPlatformRevisionOk() (*int64, bool) {
	if o == nil || IsNil(o.CurrentDataPlatformRevision) {
		return nil, false
	}
	return o.CurrentDataPlatformRevision, true
}

// HasCurrentDataPlatformRevision returns a boolean if a field has been set.
func (o *Metadata) HasCurrentDataPlatformRevision() bool {
	if o != nil && !IsNil(o.CurrentDataPlatformRevision) {
		return true
	}

	return false
}

// SetCurrentDataPlatformRevision gets a reference to the given int64 and assigns it to the CurrentDataPlatformRevision field.
func (o *Metadata) SetCurrentDataPlatformRevision(v int64) {
	o.CurrentDataPlatformRevision = &v
}

// GetAvailableUpgradeVersions returns the AvailableUpgradeVersions field value if set, zero value otherwise.
func (o *Metadata) GetAvailableUpgradeVersions() []string {
	if o == nil || IsNil(o.AvailableUpgradeVersions) {
		var ret []string
		return ret
	}
	return o.AvailableUpgradeVersions
}

// GetAvailableUpgradeVersionsOk returns a tuple with the AvailableUpgradeVersions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetAvailableUpgradeVersionsOk() ([]string, bool) {
	if o == nil || IsNil(o.AvailableUpgradeVersions) {
		return nil, false
	}
	return o.AvailableUpgradeVersions, true
}

// HasAvailableUpgradeVersions returns a boolean if a field has been set.
func (o *Metadata) HasAvailableUpgradeVersions() bool {
	if o != nil && !IsNil(o.AvailableUpgradeVersions) {
		return true
	}

	return false
}

// SetAvailableUpgradeVersions gets a reference to the given []string and assigns it to the AvailableUpgradeVersions field.
func (o *Metadata) SetAvailableUpgradeVersions(v []string) {
	o.AvailableUpgradeVersions = v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Metadata) GetState() string {
	if o == nil || IsNil(o.State) {
		var ret string
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Metadata) GetStateOk() (*string, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Metadata) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given string and assigns it to the State field.
func (o *Metadata) SetState(v string) {
	o.State = &v
}

func (o Metadata) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Metadata) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ETag) {
		toSerialize["ETag"] = o.ETag
	}
	if !IsNil(o.CreatedDate) {
		toSerialize["createdDate"] = o.CreatedDate
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["createdBy"] = o.CreatedBy
	}
	if !IsNil(o.CreatedByUserId) {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}
	if !IsNil(o.CreatedInContractNumber) {
		toSerialize["createdInContractNumber"] = o.CreatedInContractNumber
	}
	if !IsNil(o.LastModifiedDate) {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}
	if !IsNil(o.LastModifiedBy) {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}
	if !IsNil(o.LastModifiedByUserId) {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}
	if !IsNil(o.CurrentDataPlatformVersion) {
		toSerialize["currentDataPlatformVersion"] = o.CurrentDataPlatformVersion
	}
	if !IsNil(o.CurrentDataPlatformRevision) {
		toSerialize["currentDataPlatformRevision"] = o.CurrentDataPlatformRevision
	}
	if !IsNil(o.AvailableUpgradeVersions) {
		toSerialize["availableUpgradeVersions"] = o.AvailableUpgradeVersions
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	return toSerialize, nil
}

type NullableMetadata struct {
	value *Metadata
	isSet bool
}

func (v NullableMetadata) Get() *Metadata {
	return v.value
}

func (v *NullableMetadata) Set(val *Metadata) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadata) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadata) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadata(val *Metadata) *NullableMetadata {
	return &NullableMetadata{value: val, isSet: true}
}

func (v NullableMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadata) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
