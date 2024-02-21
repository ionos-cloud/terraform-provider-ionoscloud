/*
 * Container Registry service
 *
 * ## Overview Container Registry service enables IONOS clients to manage docker and OCI compliant registries for use by their managed Kubernetes clusters. Use a Container Registry to ensure you have a privately accessed registry to efficiently support image pulls. ## Changelog ### 1.1.0  - Added new endpoints for Repositories  - Added new endpoints for Artifacts  - Added new endpoints for Vulnerabilities  - Added registry vulnerabilityScanning feature
 *
 * API version: 1.1.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Purl struct for Purl
type Purl struct {
	// The affected package type
	Type *string `json:"type"`
	// The affected package name
	Name *string `json:"name"`
	// The affected package version
	Version *string `json:"version"`
}

// NewPurl instantiates a new Purl object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPurl(type_ string, name string, version string) *Purl {
	this := Purl{}

	this.Type = &type_
	this.Name = &name
	this.Version = &version

	return &this
}

// NewPurlWithDefaults instantiates a new Purl object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPurlWithDefaults() *Purl {
	this := Purl{}
	return &this
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Purl) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Purl) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *Purl) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *Purl) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Purl) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Purl) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *Purl) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *Purl) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetVersion returns the Version field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Purl) GetVersion() *string {
	if o == nil {
		return nil
	}

	return o.Version

}

// GetVersionOk returns a tuple with the Version field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Purl) GetVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Version, true
}

// SetVersion sets field value
func (o *Purl) SetVersion(v string) {

	o.Version = &v

}

// HasVersion returns a boolean if a field has been set.
func (o *Purl) HasVersion() bool {
	if o != nil && o.Version != nil {
		return true
	}

	return false
}

func (o Purl) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Version != nil {
		toSerialize["version"] = o.Version
	}

	return json.Marshal(toSerialize)
}

type NullablePurl struct {
	value *Purl
	isSet bool
}

func (v NullablePurl) Get() *Purl {
	return v.value
}

func (v *NullablePurl) Set(val *Purl) {
	v.value = val
	v.isSet = true
}

func (v NullablePurl) IsSet() bool {
	return v.isSet
}

func (v *NullablePurl) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePurl(val *Purl) *NullablePurl {
	return &NullablePurl{value: val, isSet: true}
}

func (v NullablePurl) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePurl) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}