/*
 * Certificate Manager Service API
 *
 * Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates  with IONOS services and your internal connected resources.   For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic. The service provides the basic functions of uploading and deleting your certificates for this purpose.
 *
 * API version: 2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// CertificatePatch struct for CertificatePatch
type CertificatePatch struct {
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *PatchName              `json:"properties"`
}

// NewCertificatePatch instantiates a new CertificatePatch object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertificatePatch(properties PatchName) *CertificatePatch {
	this := CertificatePatch{}

	this.Properties = &properties

	return &this
}

// NewCertificatePatchWithDefaults instantiates a new CertificatePatch object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertificatePatchWithDefaults() *CertificatePatch {
	this := CertificatePatch{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *CertificatePatch) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePatch) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *CertificatePatch) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *CertificatePatch) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for PatchName will be returned
func (o *CertificatePatch) GetProperties() *PatchName {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePatch) GetPropertiesOk() (*PatchName, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *CertificatePatch) SetProperties(v PatchName) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *CertificatePatch) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o CertificatePatch) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableCertificatePatch struct {
	value *CertificatePatch
	isSet bool
}

func (v NullableCertificatePatch) Get() *CertificatePatch {
	return v.value
}

func (v *NullableCertificatePatch) Set(val *CertificatePatch) {
	v.value = val
	v.isSet = true
}

func (v NullableCertificatePatch) IsSet() bool {
	return v.isSet
}

func (v *NullableCertificatePatch) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertificatePatch(val *CertificatePatch) *NullableCertificatePatch {
	return &NullableCertificatePatch{value: val, isSet: true}
}

func (v NullableCertificatePatch) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertificatePatch) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
