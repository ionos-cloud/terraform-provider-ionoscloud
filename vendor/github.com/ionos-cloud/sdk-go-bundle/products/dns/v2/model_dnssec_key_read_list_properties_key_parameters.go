/*
 * IONOS Cloud - DNS API
 *
 * Cloud DNS service helps IONOS Cloud customers to automate DNS Zone and Record management.
 *
 * API version: 1.17.0
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dns

import (
	"encoding/json"
)

// checks if the DnssecKeyReadListPropertiesKeyParameters type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DnssecKeyReadListPropertiesKeyParameters{}

// DnssecKeyReadListPropertiesKeyParameters struct for DnssecKeyReadListPropertiesKeyParameters
type DnssecKeyReadListPropertiesKeyParameters struct {
	Algorithm *Algorithm `json:"algorithm,omitempty"`
}

// NewDnssecKeyReadListPropertiesKeyParameters instantiates a new DnssecKeyReadListPropertiesKeyParameters object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnssecKeyReadListPropertiesKeyParameters() *DnssecKeyReadListPropertiesKeyParameters {
	this := DnssecKeyReadListPropertiesKeyParameters{}

	return &this
}

// NewDnssecKeyReadListPropertiesKeyParametersWithDefaults instantiates a new DnssecKeyReadListPropertiesKeyParameters object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnssecKeyReadListPropertiesKeyParametersWithDefaults() *DnssecKeyReadListPropertiesKeyParameters {
	this := DnssecKeyReadListPropertiesKeyParameters{}
	return &this
}

// GetAlgorithm returns the Algorithm field value if set, zero value otherwise.
func (o *DnssecKeyReadListPropertiesKeyParameters) GetAlgorithm() Algorithm {
	if o == nil || IsNil(o.Algorithm) {
		var ret Algorithm
		return ret
	}
	return *o.Algorithm
}

// GetAlgorithmOk returns a tuple with the Algorithm field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnssecKeyReadListPropertiesKeyParameters) GetAlgorithmOk() (*Algorithm, bool) {
	if o == nil || IsNil(o.Algorithm) {
		return nil, false
	}
	return o.Algorithm, true
}

// HasAlgorithm returns a boolean if a field has been set.
func (o *DnssecKeyReadListPropertiesKeyParameters) HasAlgorithm() bool {
	if o != nil && !IsNil(o.Algorithm) {
		return true
	}

	return false
}

// SetAlgorithm gets a reference to the given Algorithm and assigns it to the Algorithm field.
func (o *DnssecKeyReadListPropertiesKeyParameters) SetAlgorithm(v Algorithm) {
	o.Algorithm = &v
}

func (o DnssecKeyReadListPropertiesKeyParameters) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DnssecKeyReadListPropertiesKeyParameters) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Algorithm) {
		toSerialize["algorithm"] = o.Algorithm
	}
	return toSerialize, nil
}

type NullableDnssecKeyReadListPropertiesKeyParameters struct {
	value *DnssecKeyReadListPropertiesKeyParameters
	isSet bool
}

func (v NullableDnssecKeyReadListPropertiesKeyParameters) Get() *DnssecKeyReadListPropertiesKeyParameters {
	return v.value
}

func (v *NullableDnssecKeyReadListPropertiesKeyParameters) Set(val *DnssecKeyReadListPropertiesKeyParameters) {
	v.value = val
	v.isSet = true
}

func (v NullableDnssecKeyReadListPropertiesKeyParameters) IsSet() bool {
	return v.isSet
}

func (v *NullableDnssecKeyReadListPropertiesKeyParameters) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnssecKeyReadListPropertiesKeyParameters(val *DnssecKeyReadListPropertiesKeyParameters) *NullableDnssecKeyReadListPropertiesKeyParameters {
	return &NullableDnssecKeyReadListPropertiesKeyParameters{value: val, isSet: true}
}

func (v NullableDnssecKeyReadListPropertiesKeyParameters) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnssecKeyReadListPropertiesKeyParameters) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
