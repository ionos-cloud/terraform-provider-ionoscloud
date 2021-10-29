/*
 * CLOUD API
 *
 * IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0-SDK.3
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// ResourceProperties struct for ResourceProperties
type ResourceProperties struct {
	// name of the resource
	Name *string `json:"name,omitempty"`
	// Boolean value representing if the resource is multi factor protected or not e.g. using two factor protection. Currently only data centers and Snapshots are allowed to be multi factor protected, The value of attribute if null is intentional and it means that the resource doesn't support multi factor protection at all.
	SecAuthProtection *bool `json:"secAuthProtection,omitempty"`
}



// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *ResourceProperties) GetName() *string {
	if o == nil {
		return nil
	}


	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceProperties) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}


	return o.Name, true
}

// SetName sets field value
func (o *ResourceProperties) SetName(v string) {


	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *ResourceProperties) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}



// GetSecAuthProtection returns the SecAuthProtection field value
// If the value is explicit nil, the zero value for bool will be returned
func (o *ResourceProperties) GetSecAuthProtection() *bool {
	if o == nil {
		return nil
	}


	return o.SecAuthProtection

}

// GetSecAuthProtectionOk returns a tuple with the SecAuthProtection field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ResourceProperties) GetSecAuthProtectionOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}


	return o.SecAuthProtection, true
}

// SetSecAuthProtection sets field value
func (o *ResourceProperties) SetSecAuthProtection(v bool) {


	o.SecAuthProtection = &v

}

// HasSecAuthProtection returns a boolean if a field has been set.
func (o *ResourceProperties) HasSecAuthProtection() bool {
	if o != nil && o.SecAuthProtection != nil {
		return true
	}

	return false
}


func (o ResourceProperties) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	

	if o.SecAuthProtection != nil {
		toSerialize["secAuthProtection"] = o.SecAuthProtection
	}
	
	return json.Marshal(toSerialize)
}

type NullableResourceProperties struct {
	value *ResourceProperties
	isSet bool
}

func (v NullableResourceProperties) Get() *ResourceProperties {
	return v.value
}

func (v *NullableResourceProperties) Set(val *ResourceProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableResourceProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableResourceProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableResourceProperties(val *ResourceProperties) *NullableResourceProperties {
	return &NullableResourceProperties{value: val, isSet: true}
}

func (v NullableResourceProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableResourceProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


