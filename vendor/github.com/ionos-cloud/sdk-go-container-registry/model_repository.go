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

// Repository struct for Repository
type Repository struct {
	Name *string `json:"name"`
}

// NewRepository instantiates a new Repository object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRepository(name string) *Repository {
	this := Repository{}

	this.Name = &name

	return &this
}

// NewRepositoryWithDefaults instantiates a new Repository object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRepositoryWithDefaults() *Repository {
	this := Repository{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *Repository) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Repository) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *Repository) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *Repository) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

func (o Repository) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	return json.Marshal(toSerialize)
}

type NullableRepository struct {
	value *Repository
	isSet bool
}

func (v NullableRepository) Get() *Repository {
	return v.value
}

func (v *NullableRepository) Set(val *Repository) {
	v.value = val
	v.isSet = true
}

func (v NullableRepository) IsSet() bool {
	return v.isSet
}

func (v *NullableRepository) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRepository(val *Repository) *NullableRepository {
	return &NullableRepository{value: val, isSet: true}
}

func (v NullableRepository) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRepository) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
