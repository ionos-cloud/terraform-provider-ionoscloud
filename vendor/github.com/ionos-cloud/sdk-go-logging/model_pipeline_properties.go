/*
 * IONOS Logging REST API
 *
 * The logging service offers a centralized platform to collect and store logs from various systems and applications. It includes tools to search, filter, visualize, and create alerts based on your log data.  This API provides programmatic control over logging pipelines, enabling you to create new pipelines or modify existing ones. It mirrors the functionality of the DCD visual tool, ensuring a consistent experience regardless of your chosen interface.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// PipelineProperties A pipeline properties
type PipelineProperties struct {
	// The friendly name of your pipeline.
	Name *string `json:"name,omitempty"`
	// The information of the log aggregator
	Logs *[]PipelineResponse `json:"logs,omitempty"`
	// The address to connect fluentBit compatible logging agents to
	TcpAddress *string `json:"tcpAddress,omitempty"`
	// The address to post logs to using JSON with basic auth
	HttpAddress *string `json:"httpAddress,omitempty"`
	// The address of the client's grafana instance
	GrafanaAddress *string `json:"grafanaAddress,omitempty"`
}

// NewPipelineProperties instantiates a new PipelineProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineProperties() *PipelineProperties {
	this := PipelineProperties{}

	return &this
}

// NewPipelinePropertiesWithDefaults instantiates a new PipelineProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelinePropertiesWithDefaults() *PipelineProperties {
	this := PipelineProperties{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PipelineProperties) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *PipelineProperties) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *PipelineProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetLogs returns the Logs field value
// If the value is explicit nil, the zero value for []PipelineResponse will be returned
func (o *PipelineProperties) GetLogs() *[]PipelineResponse {
	if o == nil {
		return nil
	}

	return o.Logs

}

// GetLogsOk returns a tuple with the Logs field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineProperties) GetLogsOk() (*[]PipelineResponse, bool) {
	if o == nil {
		return nil, false
	}

	return o.Logs, true
}

// SetLogs sets field value
func (o *PipelineProperties) SetLogs(v []PipelineResponse) {

	o.Logs = &v

}

// HasLogs returns a boolean if a field has been set.
func (o *PipelineProperties) HasLogs() bool {
	if o != nil && o.Logs != nil {
		return true
	}

	return false
}

// GetTcpAddress returns the TcpAddress field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PipelineProperties) GetTcpAddress() *string {
	if o == nil {
		return nil
	}

	return o.TcpAddress

}

// GetTcpAddressOk returns a tuple with the TcpAddress field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineProperties) GetTcpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.TcpAddress, true
}

// SetTcpAddress sets field value
func (o *PipelineProperties) SetTcpAddress(v string) {

	o.TcpAddress = &v

}

// HasTcpAddress returns a boolean if a field has been set.
func (o *PipelineProperties) HasTcpAddress() bool {
	if o != nil && o.TcpAddress != nil {
		return true
	}

	return false
}

// GetHttpAddress returns the HttpAddress field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PipelineProperties) GetHttpAddress() *string {
	if o == nil {
		return nil
	}

	return o.HttpAddress

}

// GetHttpAddressOk returns a tuple with the HttpAddress field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineProperties) GetHttpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.HttpAddress, true
}

// SetHttpAddress sets field value
func (o *PipelineProperties) SetHttpAddress(v string) {

	o.HttpAddress = &v

}

// HasHttpAddress returns a boolean if a field has been set.
func (o *PipelineProperties) HasHttpAddress() bool {
	if o != nil && o.HttpAddress != nil {
		return true
	}

	return false
}

// GetGrafanaAddress returns the GrafanaAddress field value
// If the value is explicit nil, the zero value for string will be returned
func (o *PipelineProperties) GetGrafanaAddress() *string {
	if o == nil {
		return nil
	}

	return o.GrafanaAddress

}

// GetGrafanaAddressOk returns a tuple with the GrafanaAddress field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PipelineProperties) GetGrafanaAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.GrafanaAddress, true
}

// SetGrafanaAddress sets field value
func (o *PipelineProperties) SetGrafanaAddress(v string) {

	o.GrafanaAddress = &v

}

// HasGrafanaAddress returns a boolean if a field has been set.
func (o *PipelineProperties) HasGrafanaAddress() bool {
	if o != nil && o.GrafanaAddress != nil {
		return true
	}

	return false
}

func (o PipelineProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Logs != nil {
		toSerialize["logs"] = o.Logs
	}

	if o.TcpAddress != nil {
		toSerialize["tcpAddress"] = o.TcpAddress
	}

	if o.HttpAddress != nil {
		toSerialize["httpAddress"] = o.HttpAddress
	}

	if o.GrafanaAddress != nil {
		toSerialize["grafanaAddress"] = o.GrafanaAddress
	}

	return json.Marshal(toSerialize)
}

type NullablePipelineProperties struct {
	value *PipelineProperties
	isSet bool
}

func (v NullablePipelineProperties) Get() *PipelineProperties {
	return v.value
}

func (v *NullablePipelineProperties) Set(val *PipelineProperties) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineProperties) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineProperties(val *PipelineProperties) *NullablePipelineProperties {
	return &NullablePipelineProperties{value: val, isSet: true}
}

func (v NullablePipelineProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
