/*
 * IONOS Cloud - Monitoring REST API
 *
 * The monitoring service offers a centralized platform to collect and store metrics.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package monitoring

import (
	"encoding/json"
)

// checks if the PipelineRead type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PipelineRead{}

// PipelineRead struct for PipelineRead
type PipelineRead struct {
	// The ID (UUID) of the Pipeline.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the Pipeline.
	Href       string               `json:"href"`
	Metadata   MetadataWithEndpoint `json:"metadata"`
	Properties Pipeline             `json:"properties"`
}

// NewPipelineRead instantiates a new PipelineRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPipelineRead(id string, type_ string, href string, metadata MetadataWithEndpoint, properties Pipeline) *PipelineRead {
	this := PipelineRead{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewPipelineReadWithDefaults instantiates a new PipelineRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPipelineReadWithDefaults() *PipelineRead {
	this := PipelineRead{}
	return &this
}

// GetId returns the Id field value
func (o *PipelineRead) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *PipelineRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *PipelineRead) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *PipelineRead) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *PipelineRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *PipelineRead) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *PipelineRead) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *PipelineRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *PipelineRead) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *PipelineRead) GetMetadata() MetadataWithEndpoint {
	if o == nil {
		var ret MetadataWithEndpoint
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *PipelineRead) GetMetadataOk() (*MetadataWithEndpoint, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *PipelineRead) SetMetadata(v MetadataWithEndpoint) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *PipelineRead) GetProperties() Pipeline {
	if o == nil {
		var ret Pipeline
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *PipelineRead) GetPropertiesOk() (*Pipeline, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *PipelineRead) SetProperties(v Pipeline) {
	o.Properties = v
}

func (o PipelineRead) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	toSerialize["metadata"] = o.Metadata
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullablePipelineRead struct {
	value *PipelineRead
	isSet bool
}

func (v NullablePipelineRead) Get() *PipelineRead {
	return v.value
}

func (v *NullablePipelineRead) Set(val *PipelineRead) {
	v.value = val
	v.isSet = true
}

func (v NullablePipelineRead) IsSet() bool {
	return v.isSet
}

func (v *NullablePipelineRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePipelineRead(val *PipelineRead) *NullablePipelineRead {
	return &NullablePipelineRead{value: val, isSet: true}
}

func (v NullablePipelineRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePipelineRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}