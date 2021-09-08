/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// LocationProperties struct for LocationProperties
type LocationProperties struct {
	// A name of that resource
	Name *string `json:"name,omitempty"`
	// List of features supported by the location
	Features *[]string `json:"features,omitempty"`
	// List of image aliases available for the location
	ImageAliases *[]string `json:"imageAliases,omitempty"`
	// Array of features and CPU families available in a location
	CpuArchitecture *[]CpuArchitectureProperties `json:"cpuArchitecture,omitempty"`
}



// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *LocationProperties) GetName() *string {
	if o == nil {
		return nil
	}


	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LocationProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Name, true
}

// SetName sets field value
func (o *LocationProperties) SetName(v string) {


	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *LocationProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}



// GetFeatures returns the Features field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *LocationProperties) GetFeatures() *[]string {
	if o == nil {
		return nil
	}


	return o.Features

}

// GetFeaturesOk returns a tuple with the Features field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LocationProperties) GetFeaturesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Features, true
}

// SetFeatures sets field value
func (o *LocationProperties) SetFeatures(v []string) {


	o.Features = &v

}

// HasFeatures returns a boolean if a field has been set.
func (o *LocationProperties) HasFeatures() bool {
	if o != nil && o.Features != nil {
		return true
	}

	return false
}



// GetImageAliases returns the ImageAliases field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *LocationProperties) GetImageAliases() *[]string {
	if o == nil {
		return nil
	}


	return o.ImageAliases

}

// GetImageAliasesOk returns a tuple with the ImageAliases field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LocationProperties) GetImageAliasesOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}


	return o.ImageAliases, true
}

// SetImageAliases sets field value
func (o *LocationProperties) SetImageAliases(v []string) {


	o.ImageAliases = &v

}

// HasImageAliases returns a boolean if a field has been set.
func (o *LocationProperties) HasImageAliases() bool {
	if o != nil && o.ImageAliases != nil {
		return true
	}

	return false
}



// GetCpuArchitecture returns the CpuArchitecture field value
// If the value is explicit nil, the zero value for []CpuArchitectureProperties will be returned
func (o *LocationProperties) GetCpuArchitecture() *[]CpuArchitectureProperties {
	if o == nil {
		return nil
	}


	return o.CpuArchitecture

}

// GetCpuArchitectureOk returns a tuple with the CpuArchitecture field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *LocationProperties) GetCpuArchitectureOk() (*[]CpuArchitectureProperties, bool) {
	if o == nil {
		return nil, false
	}


	return o.CpuArchitecture, true
}

// SetCpuArchitecture sets field value
func (o *LocationProperties) SetCpuArchitecture(v []CpuArchitectureProperties) {


	o.CpuArchitecture = &v

}

// HasCpuArchitecture returns a boolean if a field has been set.
func (o *LocationProperties) HasCpuArchitecture() bool {
	if o != nil && o.CpuArchitecture != nil {
		return true
	}

	return false
}


func (o LocationProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	

	if o.Features != nil {
		toSerialize["features"] = o.Features
	}
	

	if o.ImageAliases != nil {
		toSerialize["imageAliases"] = o.ImageAliases
	}
	

	if o.CpuArchitecture != nil {
		toSerialize["cpuArchitecture"] = o.CpuArchitecture
	}
	
	return json.Marshal(toSerialize)
}

type NullableLocationProperties struct {
	value *LocationProperties
	isSet bool
}

func (v NullableLocationProperties) Get() *LocationProperties {
	return v.value
}

func (v *NullableLocationProperties) Set(val *LocationProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableLocationProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableLocationProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableLocationProperties(val *LocationProperties) *NullableLocationProperties {
	return &NullableLocationProperties{value: val, isSet: true}
}

func (v NullableLocationProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableLocationProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


