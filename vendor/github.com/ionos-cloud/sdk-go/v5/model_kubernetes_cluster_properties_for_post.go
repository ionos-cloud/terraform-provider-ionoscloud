/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesClusterPropertiesForPost struct for KubernetesClusterPropertiesForPost
type KubernetesClusterPropertiesForPost struct {
	// A Kubernetes Cluster Name. Valid Kubernetes Cluster name must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
	Name *string `json:"name"`
	// The kubernetes version in which a cluster is running. This imposes restrictions on what kubernetes versions can be run in a cluster's nodepools. Additionally, not all kubernetes versions are viable upgrade targets for all prior versions.
	K8sVersion        *string                      `json:"k8sVersion,omitempty"`
	MaintenanceWindow *KubernetesMaintenanceWindow `json:"maintenanceWindow,omitempty"`
	// Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6.
	ApiSubnetAllowList *[]string `json:"apiSubnetAllowList,omitempty"`
	// List of S3 bucket configured for K8s usage. For now it contains only one S3 bucket used to store K8s API audit logs
	S3Buckets *[]S3Bucket `json:"s3Buckets,omitempty"`
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesClusterPropertiesForPost) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesClusterPropertiesForPost) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *KubernetesClusterPropertiesForPost) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *KubernetesClusterPropertiesForPost) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetK8sVersion returns the K8sVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesClusterPropertiesForPost) GetK8sVersion() *string {
	if o == nil {
		return nil
	}

	return o.K8sVersion

}

// GetK8sVersionOk returns a tuple with the K8sVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesClusterPropertiesForPost) GetK8sVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.K8sVersion, true
}

// SetK8sVersion sets field value
func (o *KubernetesClusterPropertiesForPost) SetK8sVersion(v string) {

	o.K8sVersion = &v

}

// HasK8sVersion returns a boolean if a field has been set.
func (o *KubernetesClusterPropertiesForPost) HasK8sVersion() bool {
	if o != nil && o.K8sVersion != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for KubernetesMaintenanceWindow will be returned
func (o *KubernetesClusterPropertiesForPost) GetMaintenanceWindow() *KubernetesMaintenanceWindow {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesClusterPropertiesForPost) GetMaintenanceWindowOk() (*KubernetesMaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *KubernetesClusterPropertiesForPost) SetMaintenanceWindow(v KubernetesMaintenanceWindow) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *KubernetesClusterPropertiesForPost) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetApiSubnetAllowList returns the ApiSubnetAllowList field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *KubernetesClusterPropertiesForPost) GetApiSubnetAllowList() *[]string {
	if o == nil {
		return nil
	}

	return o.ApiSubnetAllowList

}

// GetApiSubnetAllowListOk returns a tuple with the ApiSubnetAllowList field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesClusterPropertiesForPost) GetApiSubnetAllowListOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.ApiSubnetAllowList, true
}

// SetApiSubnetAllowList sets field value
func (o *KubernetesClusterPropertiesForPost) SetApiSubnetAllowList(v []string) {

	o.ApiSubnetAllowList = &v

}

// HasApiSubnetAllowList returns a boolean if a field has been set.
func (o *KubernetesClusterPropertiesForPost) HasApiSubnetAllowList() bool {
	if o != nil && o.ApiSubnetAllowList != nil {
		return true
	}

	return false
}

// GetS3Buckets returns the S3Buckets field value
// If the value is explicit nil, the zero value for []S3Bucket will be returned
func (o *KubernetesClusterPropertiesForPost) GetS3Buckets() *[]S3Bucket {
	if o == nil {
		return nil
	}

	return o.S3Buckets

}

// GetS3BucketsOk returns a tuple with the S3Buckets field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesClusterPropertiesForPost) GetS3BucketsOk() (*[]S3Bucket, bool) {
	if o == nil {
		return nil, false
	}

	return o.S3Buckets, true
}

// SetS3Buckets sets field value
func (o *KubernetesClusterPropertiesForPost) SetS3Buckets(v []S3Bucket) {

	o.S3Buckets = &v

}

// HasS3Buckets returns a boolean if a field has been set.
func (o *KubernetesClusterPropertiesForPost) HasS3Buckets() bool {
	if o != nil && o.S3Buckets != nil {
		return true
	}

	return false
}

func (o KubernetesClusterPropertiesForPost) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.K8sVersion != nil {
		toSerialize["k8sVersion"] = o.K8sVersion
	}
	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}
	if o.ApiSubnetAllowList != nil {
		toSerialize["apiSubnetAllowList"] = o.ApiSubnetAllowList
	}
	if o.S3Buckets != nil {
		toSerialize["s3Buckets"] = o.S3Buckets
	}
	return json.Marshal(toSerialize)
}

type NullableKubernetesClusterPropertiesForPost struct {
	value *KubernetesClusterPropertiesForPost
	isSet bool
}

func (v NullableKubernetesClusterPropertiesForPost) Get() *KubernetesClusterPropertiesForPost {
	return v.value
}

func (v *NullableKubernetesClusterPropertiesForPost) Set(val *KubernetesClusterPropertiesForPost) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesClusterPropertiesForPost) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesClusterPropertiesForPost) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesClusterPropertiesForPost(val *KubernetesClusterPropertiesForPost) *NullableKubernetesClusterPropertiesForPost {
	return &NullableKubernetesClusterPropertiesForPost{value: val, isSet: true}
}

func (v NullableKubernetesClusterPropertiesForPost) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesClusterPropertiesForPost) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
