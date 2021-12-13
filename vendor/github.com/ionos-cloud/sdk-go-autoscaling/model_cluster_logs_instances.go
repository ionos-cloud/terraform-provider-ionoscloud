/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ClusterLogsInstances struct for ClusterLogsInstances
type ClusterLogsInstances struct {
	// The name of the PostgreSQL instance.
	Name     *string                   `json:"name,omitempty"`
	Messages *[]map[string]interface{} `json:"messages,omitempty"`
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ClusterLogsInstances) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterLogsInstances) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *ClusterLogsInstances) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *ClusterLogsInstances) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetMessages returns the Messages field value
// If the value is explicit nil, the zero value for []map[string]interface{} will be returned
func (o *ClusterLogsInstances) GetMessages() *[]map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Messages

}

// GetMessagesOk returns a tuple with the Messages field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ClusterLogsInstances) GetMessagesOk() (*[]map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Messages, true
}

// SetMessages sets field value
func (o *ClusterLogsInstances) SetMessages(v []map[string]interface{}) {

	o.Messages = &v

}

// HasMessages returns a boolean if a field has been set.
func (o *ClusterLogsInstances) HasMessages() bool {
	if o != nil && o.Messages != nil {
		return true
	}

	return false
}

func (o ClusterLogsInstances) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Messages != nil {
		toSerialize["messages"] = o.Messages
	}

	return json.Marshal(toSerialize)
}

type NullableClusterLogsInstances struct {
	value *ClusterLogsInstances
	isSet bool
}

func (v NullableClusterLogsInstances) Get() *ClusterLogsInstances {
	return v.value
}

func (v *NullableClusterLogsInstances) Set(val *ClusterLogsInstances) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterLogsInstances) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterLogsInstances) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterLogsInstances(val *ClusterLogsInstances) *NullableClusterLogsInstances {
	return &NullableClusterLogsInstances{value: val, isSet: true}
}

func (v NullableClusterLogsInstances) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterLogsInstances) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
