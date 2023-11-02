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
	"time"
)

// Metadata Metadata of the resource
type Metadata struct {
	// The Entity Tag of the resource as defined in http://www.w3.org/Protocols/rfc2616/rfc2616-sec3.html#sec3.11
	ETag *string `json:"ETag,omitempty"`
	// The time the resource was created, ISO 8601 timestamp (UTC).
	CreatedDate *IonosTime `json:"createdDate,omitempty"`
	// The user that created the resource
	CreatedBy *string `json:"createdBy,omitempty"`
	// The ID of the user that created the resource
	CreatedByUserId *string `json:"createdByUserId,omitempty"`
	// The creators contractNumber
	CreatedInContractNumber *string `json:"createdInContractNumber,omitempty"`
	// The last time the resource was modified, ISO 8601 timestamp (UTC).
	LastModifiedDate *IonosTime `json:"lastModifiedDate,omitempty"`
	// The user that last modified the resource
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`
	// The ID of the user that last modified the resource
	LastModifiedByUserId *string `json:"lastModifiedByUserId,omitempty"`
	// The version of the DataPlatform.
	CurrentDataPlatformVersion *string `json:"currentDataPlatformVersion,omitempty"`
	// The current dataplatform revision of a resource. This internal revision is used to rollout non-breaking internal changes. This attribute is read-only.
	CurrentDataPlatformRevision *int64 `json:"currentDataPlatformRevision,omitempty"`
	// List of available upgrades for this cluster
	AvailableUpgradeVersions *[]string `json:"availableUpgradeVersions,omitempty"`
	// State of the resource. *AVAILABLE* There are no pending modification requests for this item; *BUSY* There is at least one modification request pending and all following requests will be queued; *DEPLOYING* Resource state DEPLOYING - the resource is being created; *FAILED* Resource state FAILED - creation of the resource failed; *UPDATING* Resource state UPDATING - the resource is being updated; *FAILED_UPDATING* Resource state FAILED_UPDATING - an update to the resource was not successful; *DESTROYING* Resource state DESTROYING - the resource is being deleted; *FAILED_DESTROYING* Resource state FAILED_DESTROYING - deletion of the resource was not successful; *TERMINATED* Resource state TERMINATED - the resource was deleted.
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

// GetETag returns the ETag field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetETag() *string {
	if o == nil {
		return nil
	}

	return o.ETag

}

// GetETagOk returns a tuple with the ETag field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetETagOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ETag, true
}

// SetETag sets field value
func (o *Metadata) SetETag(v string) {

	o.ETag = &v

}

// HasETag returns a boolean if a field has been set.
func (o *Metadata) HasETag() bool {
	if o != nil && o.ETag != nil {
		return true
	}

	return false
}

// GetCreatedDate returns the CreatedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *Metadata) GetCreatedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.CreatedDate == nil {
		return nil
	}
	return &o.CreatedDate.Time

}

// GetCreatedDateOk returns a tuple with the CreatedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.CreatedDate == nil {
		return nil, false
	}
	return &o.CreatedDate.Time, true

}

// SetCreatedDate sets field value
func (o *Metadata) SetCreatedDate(v time.Time) {

	o.CreatedDate = &IonosTime{v}

}

// HasCreatedDate returns a boolean if a field has been set.
func (o *Metadata) HasCreatedDate() bool {
	if o != nil && o.CreatedDate != nil {
		return true
	}

	return false
}

// GetCreatedBy returns the CreatedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCreatedBy() *string {
	if o == nil {
		return nil
	}

	return o.CreatedBy

}

// GetCreatedByOk returns a tuple with the CreatedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedBy, true
}

// SetCreatedBy sets field value
func (o *Metadata) SetCreatedBy(v string) {

	o.CreatedBy = &v

}

// HasCreatedBy returns a boolean if a field has been set.
func (o *Metadata) HasCreatedBy() bool {
	if o != nil && o.CreatedBy != nil {
		return true
	}

	return false
}

// GetCreatedByUserId returns the CreatedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCreatedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.CreatedByUserId

}

// GetCreatedByUserIdOk returns a tuple with the CreatedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedByUserId, true
}

// SetCreatedByUserId sets field value
func (o *Metadata) SetCreatedByUserId(v string) {

	o.CreatedByUserId = &v

}

// HasCreatedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasCreatedByUserId() bool {
	if o != nil && o.CreatedByUserId != nil {
		return true
	}

	return false
}

// GetCreatedInContractNumber returns the CreatedInContractNumber field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCreatedInContractNumber() *string {
	if o == nil {
		return nil
	}

	return o.CreatedInContractNumber

}

// GetCreatedInContractNumberOk returns a tuple with the CreatedInContractNumber field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCreatedInContractNumberOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CreatedInContractNumber, true
}

// SetCreatedInContractNumber sets field value
func (o *Metadata) SetCreatedInContractNumber(v string) {

	o.CreatedInContractNumber = &v

}

// HasCreatedInContractNumber returns a boolean if a field has been set.
func (o *Metadata) HasCreatedInContractNumber() bool {
	if o != nil && o.CreatedInContractNumber != nil {
		return true
	}

	return false
}

// GetLastModifiedDate returns the LastModifiedDate field value
// If the value is explicit nil, the zero value for time.Time will be returned
func (o *Metadata) GetLastModifiedDate() *time.Time {
	if o == nil {
		return nil
	}

	if o.LastModifiedDate == nil {
		return nil
	}
	return &o.LastModifiedDate.Time

}

// GetLastModifiedDateOk returns a tuple with the LastModifiedDate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedDateOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}

	if o.LastModifiedDate == nil {
		return nil, false
	}
	return &o.LastModifiedDate.Time, true

}

// SetLastModifiedDate sets field value
func (o *Metadata) SetLastModifiedDate(v time.Time) {

	o.LastModifiedDate = &IonosTime{v}

}

// HasLastModifiedDate returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedDate() bool {
	if o != nil && o.LastModifiedDate != nil {
		return true
	}

	return false
}

// GetLastModifiedBy returns the LastModifiedBy field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetLastModifiedBy() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedBy

}

// GetLastModifiedByOk returns a tuple with the LastModifiedBy field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedByOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedBy, true
}

// SetLastModifiedBy sets field value
func (o *Metadata) SetLastModifiedBy(v string) {

	o.LastModifiedBy = &v

}

// HasLastModifiedBy returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedBy() bool {
	if o != nil && o.LastModifiedBy != nil {
		return true
	}

	return false
}

// GetLastModifiedByUserId returns the LastModifiedByUserId field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetLastModifiedByUserId() *string {
	if o == nil {
		return nil
	}

	return o.LastModifiedByUserId

}

// GetLastModifiedByUserIdOk returns a tuple with the LastModifiedByUserId field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetLastModifiedByUserIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.LastModifiedByUserId, true
}

// SetLastModifiedByUserId sets field value
func (o *Metadata) SetLastModifiedByUserId(v string) {

	o.LastModifiedByUserId = &v

}

// HasLastModifiedByUserId returns a boolean if a field has been set.
func (o *Metadata) HasLastModifiedByUserId() bool {
	if o != nil && o.LastModifiedByUserId != nil {
		return true
	}

	return false
}

// GetCurrentDataPlatformVersion returns the CurrentDataPlatformVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetCurrentDataPlatformVersion() *string {
	if o == nil {
		return nil
	}

	return o.CurrentDataPlatformVersion

}

// GetCurrentDataPlatformVersionOk returns a tuple with the CurrentDataPlatformVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCurrentDataPlatformVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CurrentDataPlatformVersion, true
}

// SetCurrentDataPlatformVersion sets field value
func (o *Metadata) SetCurrentDataPlatformVersion(v string) {

	o.CurrentDataPlatformVersion = &v

}

// HasCurrentDataPlatformVersion returns a boolean if a field has been set.
func (o *Metadata) HasCurrentDataPlatformVersion() bool {
	if o != nil && o.CurrentDataPlatformVersion != nil {
		return true
	}

	return false
}

// GetCurrentDataPlatformRevision returns the CurrentDataPlatformRevision field value
// If the value is explicit nil, the zero value for int64 will be returned
func (o *Metadata) GetCurrentDataPlatformRevision() *int64 {
	if o == nil {
		return nil
	}

	return o.CurrentDataPlatformRevision

}

// GetCurrentDataPlatformRevisionOk returns a tuple with the CurrentDataPlatformRevision field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetCurrentDataPlatformRevisionOk() (*int64, bool) {
	if o == nil {
		return nil, false
	}

	return o.CurrentDataPlatformRevision, true
}

// SetCurrentDataPlatformRevision sets field value
func (o *Metadata) SetCurrentDataPlatformRevision(v int64) {

	o.CurrentDataPlatformRevision = &v

}

// HasCurrentDataPlatformRevision returns a boolean if a field has been set.
func (o *Metadata) HasCurrentDataPlatformRevision() bool {
	if o != nil && o.CurrentDataPlatformRevision != nil {
		return true
	}

	return false
}

// GetAvailableUpgradeVersions returns the AvailableUpgradeVersions field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *Metadata) GetAvailableUpgradeVersions() *[]string {
	if o == nil {
		return nil
	}

	return o.AvailableUpgradeVersions

}

// GetAvailableUpgradeVersionsOk returns a tuple with the AvailableUpgradeVersions field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetAvailableUpgradeVersionsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.AvailableUpgradeVersions, true
}

// SetAvailableUpgradeVersions sets field value
func (o *Metadata) SetAvailableUpgradeVersions(v []string) {

	o.AvailableUpgradeVersions = &v

}

// HasAvailableUpgradeVersions returns a boolean if a field has been set.
func (o *Metadata) HasAvailableUpgradeVersions() bool {
	if o != nil && o.AvailableUpgradeVersions != nil {
		return true
	}

	return false
}

// GetState returns the State field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Metadata) GetState() *string {
	if o == nil {
		return nil
	}

	return o.State

}

// GetStateOk returns a tuple with the State field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Metadata) GetStateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.State, true
}

// SetState sets field value
func (o *Metadata) SetState(v string) {

	o.State = &v

}

// HasState returns a boolean if a field has been set.
func (o *Metadata) HasState() bool {
	if o != nil && o.State != nil {
		return true
	}

	return false
}

func (o Metadata) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.ETag != nil {
		toSerialize["ETag"] = o.ETag
	}

	if o.CreatedDate != nil {
		toSerialize["createdDate"] = o.CreatedDate
	}

	if o.CreatedBy != nil {
		toSerialize["createdBy"] = o.CreatedBy
	}

	if o.CreatedByUserId != nil {
		toSerialize["createdByUserId"] = o.CreatedByUserId
	}

	if o.CreatedInContractNumber != nil {
		toSerialize["createdInContractNumber"] = o.CreatedInContractNumber
	}

	if o.LastModifiedDate != nil {
		toSerialize["lastModifiedDate"] = o.LastModifiedDate
	}

	if o.LastModifiedBy != nil {
		toSerialize["lastModifiedBy"] = o.LastModifiedBy
	}

	if o.LastModifiedByUserId != nil {
		toSerialize["lastModifiedByUserId"] = o.LastModifiedByUserId
	}

	if o.CurrentDataPlatformVersion != nil {
		toSerialize["currentDataPlatformVersion"] = o.CurrentDataPlatformVersion
	}

	if o.CurrentDataPlatformRevision != nil {
		toSerialize["currentDataPlatformRevision"] = o.CurrentDataPlatformRevision
	}

	if o.AvailableUpgradeVersions != nil {
		toSerialize["availableUpgradeVersions"] = o.AvailableUpgradeVersions
	}

	if o.State != nil {
		toSerialize["state"] = o.State
	}

	return json.Marshal(toSerialize)
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
