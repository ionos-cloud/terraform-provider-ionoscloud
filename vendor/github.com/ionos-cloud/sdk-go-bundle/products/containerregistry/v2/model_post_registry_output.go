/*
 * Container Registry service
 *
 * ## Overview Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls. ## Changelog ### 1.1.0  - Added new endpoints for Repositories  - Added new endpoints for Artifacts  - Added new endpoints for Vulnerabilities  - Added registry vulnerabilityScanning feature ### 1.2.0  - Added registry `apiSubnetAllowList` ### 1.2.1  - Amended `apiSubnetAllowList` Regex
 *
 * API version: 1.2.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package containerregistry

import (
	"encoding/json"
)

// checks if the PostRegistryOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PostRegistryOutput{}

// PostRegistryOutput struct for PostRegistryOutput
type PostRegistryOutput struct {
	Href       *string             `json:"href,omitempty"`
	Id         *string             `json:"id,omitempty"`
	Metadata   ApiResourceMetadata `json:"metadata"`
	Properties RegistryProperties  `json:"properties"`
	Type       *string             `json:"type,omitempty"`
}

// NewPostRegistryOutput instantiates a new PostRegistryOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPostRegistryOutput(metadata ApiResourceMetadata, properties RegistryProperties) *PostRegistryOutput {
	this := PostRegistryOutput{}

	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewPostRegistryOutputWithDefaults instantiates a new PostRegistryOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPostRegistryOutputWithDefaults() *PostRegistryOutput {
	this := PostRegistryOutput{}
	return &this
}

// GetHref returns the Href field value if set, zero value otherwise.
func (o *PostRegistryOutput) GetHref() string {
	if o == nil || IsNil(o.Href) {
		var ret string
		return ret
	}
	return *o.Href
}

// GetHrefOk returns a tuple with the Href field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostRegistryOutput) GetHrefOk() (*string, bool) {
	if o == nil || IsNil(o.Href) {
		return nil, false
	}
	return o.Href, true
}

// HasHref returns a boolean if a field has been set.
func (o *PostRegistryOutput) HasHref() bool {
	if o != nil && !IsNil(o.Href) {
		return true
	}

	return false
}

// SetHref gets a reference to the given string and assigns it to the Href field.
func (o *PostRegistryOutput) SetHref(v string) {
	o.Href = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *PostRegistryOutput) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostRegistryOutput) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *PostRegistryOutput) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *PostRegistryOutput) SetId(v string) {
	o.Id = &v
}

// GetMetadata returns the Metadata field value
func (o *PostRegistryOutput) GetMetadata() ApiResourceMetadata {
	if o == nil {
		var ret ApiResourceMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *PostRegistryOutput) GetMetadataOk() (*ApiResourceMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *PostRegistryOutput) SetMetadata(v ApiResourceMetadata) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *PostRegistryOutput) GetProperties() RegistryProperties {
	if o == nil {
		var ret RegistryProperties
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *PostRegistryOutput) GetPropertiesOk() (*RegistryProperties, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *PostRegistryOutput) SetProperties(v RegistryProperties) {
	o.Properties = v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *PostRegistryOutput) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PostRegistryOutput) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *PostRegistryOutput) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *PostRegistryOutput) SetType(v string) {
	o.Type = &v
}

func (o PostRegistryOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PostRegistryOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Href) {
		toSerialize["href"] = o.Href
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	toSerialize["metadata"] = o.Metadata
	toSerialize["properties"] = o.Properties
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	return toSerialize, nil
}

type NullablePostRegistryOutput struct {
	value *PostRegistryOutput
	isSet bool
}

func (v NullablePostRegistryOutput) Get() *PostRegistryOutput {
	return v.value
}

func (v *NullablePostRegistryOutput) Set(val *PostRegistryOutput) {
	v.value = val
	v.isSet = true
}

func (v NullablePostRegistryOutput) IsSet() bool {
	return v.isSet
}

func (v *NullablePostRegistryOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePostRegistryOutput(val *PostRegistryOutput) *NullablePostRegistryOutput {
	return &NullablePostRegistryOutput{value: val, isSet: true}
}

func (v NullablePostRegistryOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePostRegistryOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
