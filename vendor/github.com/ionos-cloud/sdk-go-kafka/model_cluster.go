/*
 * Kafka as a Service API
 *
 * An managed Apache Kafka cluster is designed to be highly fault-tolerant and scalable, allowing large volumes of data to be ingested, stored, and processed in real-time. By distributing data across multiple brokers, Kafka achieves high throughput and low latency, making it suitable for applications requiring real-time data processing and analytics.
 *
 * API version: 1.4.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Cluster A Kafka cluster that stores data and serve client requests. Kafka clusters typically have multiple brokers to handle more data and provide high availability. Each broker is identified by a unique ID and manages partitions of different topics.
type Cluster struct {
	// The name of your Kafka cluster. Must be 63 characters or less and must begin and end with an alphanumeric character (`[a-z0-9A-Z]`) with dashes (`-`), underscores (`_`), dots (`.`), and alphanumerics between.
	Name *string `json:"name"`
	// The version of Kafka. Currently only Kafka Version 3.7.0 is supported.
	Version *string `json:"version"`
	// The size of your Kafka cluster. The size of the Kafka cluster is given in T-shirt sizes. Valid values are: \"S\"
	Size        *string                   `json:"size"`
	Connections *[]KafkaClusterConnection `json:"connections"`
}

// NewCluster instantiates a new Cluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCluster(name string, version string, size string, connections []KafkaClusterConnection) *Cluster {
	this := Cluster{}

	this.Name = &name
	this.Version = &version
	this.Size = &size
	this.Connections = &connections

	return &this
}

// NewClusterWithDefaults instantiates a new Cluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewClusterWithDefaults() *Cluster {
	this := Cluster{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *Cluster) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *Cluster) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetVersion returns the Version field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetVersion() *string {
	if o == nil {
		return nil
	}

	return o.Version

}

// GetVersionOk returns a tuple with the Version field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Version, true
}

// SetVersion sets field value
func (o *Cluster) SetVersion(v string) {

	o.Version = &v

}

// HasVersion returns a boolean if a field has been set.
func (o *Cluster) HasVersion() bool {
	if o != nil && o.Version != nil {
		return true
	}

	return false
}

// GetSize returns the Size field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Cluster) GetSize() *string {
	if o == nil {
		return nil
	}

	return o.Size

}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetSizeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Size, true
}

// SetSize sets field value
func (o *Cluster) SetSize(v string) {

	o.Size = &v

}

// HasSize returns a boolean if a field has been set.
func (o *Cluster) HasSize() bool {
	if o != nil && o.Size != nil {
		return true
	}

	return false
}

// GetConnections returns the Connections field value
// If the value is explicit nil, the zero value for []KafkaClusterConnection will be returned
func (o *Cluster) GetConnections() *[]KafkaClusterConnection {
	if o == nil {
		return nil
	}

	return o.Connections

}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Cluster) GetConnectionsOk() (*[]KafkaClusterConnection, bool) {
	if o == nil {
		return nil, false
	}

	return o.Connections, true
}

// SetConnections sets field value
func (o *Cluster) SetConnections(v []KafkaClusterConnection) {

	o.Connections = &v

}

// HasConnections returns a boolean if a field has been set.
func (o *Cluster) HasConnections() bool {
	if o != nil && o.Connections != nil {
		return true
	}

	return false
}

func (o Cluster) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Version != nil {
		toSerialize["version"] = o.Version
	}

	if o.Size != nil {
		toSerialize["size"] = o.Size
	}

	if o.Connections != nil {
		toSerialize["connections"] = o.Connections
	}

	return json.Marshal(toSerialize)
}

type NullableCluster struct {
	value *Cluster
	isSet bool
}

func (v NullableCluster) Get() *Cluster {
	return v.value
}

func (v *NullableCluster) Set(val *Cluster) {
	v.value = val
	v.isSet = true
}

func (v NullableCluster) IsSet() bool {
	return v.isSet
}

func (v *NullableCluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCluster(val *Cluster) *NullableCluster {
	return &NullableCluster{value: val, isSet: true}
}

func (v NullableCluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}