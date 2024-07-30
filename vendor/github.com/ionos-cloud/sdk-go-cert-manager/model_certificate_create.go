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

// CertificateCreate struct for CertificateCreate
type CertificateCreate struct {
	// Metadata
	Metadata   *map[string]interface{} `json:"metadata,omitempty"`
	Properties *Certificate            `json:"properties"`
}

// NewCertificateCreate instantiates a new CertificateCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertificateCreate(properties Certificate) *CertificateCreate {
	this := CertificateCreate{}

	this.Properties = &properties

	return &this
}

// NewCertificateCreateWithDefaults instantiates a new CertificateCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertificateCreateWithDefaults() *CertificateCreate {
	this := CertificateCreate{}
	return &this
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for map[string]interface{} will be returned
func (o *CertificateCreate) GetMetadata() *map[string]interface{} {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateCreate) GetMetadataOk() (*map[string]interface{}, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *CertificateCreate) SetMetadata(v map[string]interface{}) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *CertificateCreate) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Certificate will be returned
func (o *CertificateCreate) GetProperties() *Certificate {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateCreate) GetPropertiesOk() (*Certificate, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *CertificateCreate) SetProperties(v Certificate) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *CertificateCreate) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o CertificateCreate) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableCertificateCreate struct {
	value *CertificateCreate
	isSet bool
}

func (v NullableCertificateCreate) Get() *CertificateCreate {
	return v.value
}

func (v *NullableCertificateCreate) Set(val *CertificateCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableCertificateCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableCertificateCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertificateCreate(val *CertificateCreate) *NullableCertificateCreate {
	return &NullableCertificateCreate{value: val, isSet: true}
}

func (v NullableCertificateCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertificateCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
