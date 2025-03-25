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

// checks if the RegistryPagination type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &RegistryPagination{}

// RegistryPagination struct for RegistryPagination
type RegistryPagination struct {
	// The maximum number of elements to return (used together with pagination.token for pagination)
	Limit int32 `json:"limit"`
	// An opaque token used to iterate the set of results (used together with limit for pagination)
	Token string `json:"token"`
}

// NewRegistryPagination instantiates a new RegistryPagination object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewRegistryPagination(limit int32, token string) *RegistryPagination {
	this := RegistryPagination{}

	this.Limit = limit
	this.Token = token

	return &this
}

// NewRegistryPaginationWithDefaults instantiates a new RegistryPagination object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewRegistryPaginationWithDefaults() *RegistryPagination {
	this := RegistryPagination{}
	return &this
}

// GetLimit returns the Limit field value
func (o *RegistryPagination) GetLimit() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Limit
}

// GetLimitOk returns a tuple with the Limit field value
// and a boolean to check if the value has been set.
func (o *RegistryPagination) GetLimitOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Limit, true
}

// SetLimit sets field value
func (o *RegistryPagination) SetLimit(v int32) {
	o.Limit = v
}

// GetToken returns the Token field value
func (o *RegistryPagination) GetToken() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Token
}

// GetTokenOk returns a tuple with the Token field value
// and a boolean to check if the value has been set.
func (o *RegistryPagination) GetTokenOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Token, true
}

// SetToken sets field value
func (o *RegistryPagination) SetToken(v string) {
	o.Token = v
}

func (o RegistryPagination) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o RegistryPagination) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["limit"] = o.Limit
	toSerialize["token"] = o.Token
	return toSerialize, nil
}

type NullableRegistryPagination struct {
	value *RegistryPagination
	isSet bool
}

func (v NullableRegistryPagination) Get() *RegistryPagination {
	return v.value
}

func (v *NullableRegistryPagination) Set(val *RegistryPagination) {
	v.value = val
	v.isSet = true
}

func (v NullableRegistryPagination) IsSet() bool {
	return v.isSet
}

func (v *NullableRegistryPagination) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableRegistryPagination(val *RegistryPagination) *NullableRegistryPagination {
	return &NullableRegistryPagination{value: val, isSet: true}
}

func (v NullableRegistryPagination) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableRegistryPagination) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
