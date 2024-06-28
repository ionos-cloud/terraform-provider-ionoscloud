/*
 * IONOS Logging REST API
 *
 * Logging as a Service (LaaS) is a service that provides a centralized logging system where users are able to push and aggregate their system or application logs. This service also provides a visualization platform where users are able to observe, search and filter the logs and also create dashboards and alerts for their data points. This service can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an API. The API allows you to create logging pipelines or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package logging

import (
	"encoding/json"
)

// checks if the PipelineResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelineResponse{}

// PipelineResponse struct for PipelineResponse
type PipelineResponse struct {
	Public *bool `json:"public,omitempty"`
	// The source parser to be used
	Source *string `json:"source,omitempty"`
	// Tag is to distinguish different pipelines. must be unique amongst the pipeline's array items.
	Tag *string `json:"tag,omitempty"`
	// Protocol to use as intake
	Protocol *string `json:"protocol,omitempty"`
	// Optional custom labels to filter and report logs
	Labels       []string      `json:"labels,omitempty"`
	Destinations []Destination `json:"destinations,omitempty"`
}

// NewPipelineResponse instantiates a new PipelineResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineResponse() *PipelineResponse {
	this := PipelineResponse{}

	return &this
}

// NewPipelineResponseWithDefaults instantiates a new PipelineResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineResponseWithDefaults() *PipelineResponse {
	this := PipelineResponse{}
	return &this
}

// GetPublic returns the Public field value if set, zero value otherwise.
func (o *PipelineResponse) GetPublic() bool {
	if o == nil || IsNil(o.Public) {
		var ret bool
		return ret
	}
	return *o.Public
}

// GetPublicOk returns a tuple with the Public field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetPublicOk() (*bool, bool) {
	if o == nil || IsNil(o.Public) {
		return nil, false
	}
	return o.Public, true
}

// HasPublic returns a boolean if a field has been set.
func (o *PipelineResponse) HasPublic() bool {
	if o != nil && !IsNil(o.Public) {
		return true
	}

	return false
}

// SetPublic gets a reference to the given bool and assigns it to the Public field.
func (o *PipelineResponse) SetPublic(v bool) {
	o.Public = &v
}

// GetSource returns the Source field value if set, zero value otherwise.
func (o *PipelineResponse) GetSource() string {
	if o == nil || IsNil(o.Source) {
		var ret string
		return ret
	}
	return *o.Source
}

// GetSourceOk returns a tuple with the Source field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetSourceOk() (*string, bool) {
	if o == nil || IsNil(o.Source) {
		return nil, false
	}
	return o.Source, true
}

// HasSource returns a boolean if a field has been set.
func (o *PipelineResponse) HasSource() bool {
	if o != nil && !IsNil(o.Source) {
		return true
	}

	return false
}

// SetSource gets a reference to the given string and assigns it to the Source field.
func (o *PipelineResponse) SetSource(v string) {
	o.Source = &v
}

// GetTag returns the Tag field value if set, zero value otherwise.
func (o *PipelineResponse) GetTag() string {
	if o == nil || IsNil(o.Tag) {
		var ret string
		return ret
	}
	return *o.Tag
}

// GetTagOk returns a tuple with the Tag field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetTagOk() (*string, bool) {
	if o == nil || IsNil(o.Tag) {
		return nil, false
	}
	return o.Tag, true
}

// HasTag returns a boolean if a field has been set.
func (o *PipelineResponse) HasTag() bool {
	if o != nil && !IsNil(o.Tag) {
		return true
	}

	return false
}

// SetTag gets a reference to the given string and assigns it to the Tag field.
func (o *PipelineResponse) SetTag(v string) {
	o.Tag = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *PipelineResponse) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *PipelineResponse) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *PipelineResponse) SetProtocol(v string) {
	o.Protocol = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *PipelineResponse) GetLabels() []string {
	if o == nil || IsNil(o.Labels) {
		var ret []string
		return ret
	}
	return o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetLabelsOk() ([]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *PipelineResponse) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given []string and assigns it to the Labels field.
func (o *PipelineResponse) SetLabels(v []string) {
	o.Labels = v
}

// GetDestinations returns the Destinations field value if set, zero value otherwise.
func (o *PipelineResponse) GetDestinations() []Destination {
	if o == nil || IsNil(o.Destinations) {
		var ret []Destination
		return ret
	}
	return o.Destinations
}

// GetDestinationsOk returns a tuple with the Destinations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PipelineResponse) GetDestinationsOk() ([]Destination, bool) {
	if o == nil || IsNil(o.Destinations) {
		return nil, false
	}
	return o.Destinations, true
}

// HasDestinations returns a boolean if a field has been set.
func (o *PipelineResponse) HasDestinations() bool {
	if o != nil && !IsNil(o.Destinations) {
		return true
	}

	return false
}

// SetDestinations gets a reference to the given []Destination and assigns it to the Destinations field.
func (o *PipelineResponse) SetDestinations(v []Destination) {
	o.Destinations = v
}

func (o PipelineResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PipelineResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Public) {
		toSerialize["public"] = o.Public
	}
	if !IsNil(o.Source) {
		toSerialize["source"] = o.Source
	}
	if !IsNil(o.Tag) {
		toSerialize["tag"] = o.Tag
	}
	if !IsNil(o.Protocol) {
		toSerialize["protocol"] = o.Protocol
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.Destinations) {
		toSerialize["destinations"] = o.Destinations
	}
	return toSerialize, nil
}

type NullablePipelineResponse struct {
	value *PipelineResponse
	isSet bool
}

func (v NullablePipelineResponse) Get() *PipelineResponse {
	return v.value
}

func (v *NullablePipelineResponse) Set(val *PipelineResponse) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineResponse) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineResponse(val *PipelineResponse) *NullablePipelineResponse {
	return &NullablePipelineResponse{value: val, isSet: true}
}

func (v NullablePipelineResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
