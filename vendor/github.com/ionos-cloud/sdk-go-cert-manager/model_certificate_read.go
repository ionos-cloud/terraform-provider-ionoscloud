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

// CertificateRead struct for CertificateRead
type CertificateRead struct {
	// The ID (UUID) of the Certificate.
	Id *string `json:"id"`
	// The type of the resource.
	Type *string `json:"type"`
	// The URL of the Certificate.
	Href       *string                             `json:"href"`
	Metadata   *MetadataWithCertificateInformation `json:"metadata"`
	Properties *Certificate                        `json:"properties"`
}

// NewCertificateRead instantiates a new CertificateRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertificateRead(id string, type_ string, href string, metadata MetadataWithCertificateInformation, properties Certificate) *CertificateRead {
	this := CertificateRead{}

	this.Id = &id
	this.Type = &type_
	this.Href = &href
	this.Metadata = &metadata
	this.Properties = &properties

	return &this
}

// NewCertificateReadWithDefaults instantiates a new CertificateRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertificateReadWithDefaults() *CertificateRead {
	this := CertificateRead{}
	return &this
}

// GetId returns the Id field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificateRead) GetId() *string {
	if o == nil {
		return nil
	}

	return o.Id

}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Id, true
}

// SetId sets field value
func (o *CertificateRead) SetId(v string) {

	o.Id = &v

}

// HasId returns a boolean if a field has been set.
func (o *CertificateRead) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// GetType returns the Type field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificateRead) GetType() *string {
	if o == nil {
		return nil
	}

	return o.Type

}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Type, true
}

// SetType sets field value
func (o *CertificateRead) SetType(v string) {

	o.Type = &v

}

// HasType returns a boolean if a field has been set.
func (o *CertificateRead) HasType() bool {
	if o != nil && o.Type != nil {
		return true
	}

	return false
}

// GetHref returns the Href field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificateRead) GetHref() *string {
	if o == nil {
		return nil
	}

	return o.Href

}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Href, true
}

// SetHref sets field value
func (o *CertificateRead) SetHref(v string) {

	o.Href = &v

}

// HasHref returns a boolean if a field has been set.
func (o *CertificateRead) HasHref() bool {
	if o != nil && o.Href != nil {
		return true
	}

	return false
}

// GetMetadata returns the Metadata field value
// If the value is explicit nil, the zero value for MetadataWithCertificateInformation will be returned
func (o *CertificateRead) GetMetadata() *MetadataWithCertificateInformation {
	if o == nil {
		return nil
	}

	return o.Metadata

}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateRead) GetMetadataOk() (*MetadataWithCertificateInformation, bool) {
	if o == nil {
		return nil, false
	}

	return o.Metadata, true
}

// SetMetadata sets field value
func (o *CertificateRead) SetMetadata(v MetadataWithCertificateInformation) {

	o.Metadata = &v

}

// HasMetadata returns a boolean if a field has been set.
func (o *CertificateRead) HasMetadata() bool {
	if o != nil && o.Metadata != nil {
		return true
	}

	return false
}

// GetProperties returns the Properties field value
// If the value is explicit nil, the zero value for Certificate will be returned
func (o *CertificateRead) GetProperties() *Certificate {
	if o == nil {
		return nil
	}

	return o.Properties

}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificateRead) GetPropertiesOk() (*Certificate, bool) {
	if o == nil {
		return nil, false
	}

	return o.Properties, true
}

// SetProperties sets field value
func (o *CertificateRead) SetProperties(v Certificate) {

	o.Properties = &v

}

// HasProperties returns a boolean if a field has been set.
func (o *CertificateRead) HasProperties() bool {
	if o != nil && o.Properties != nil {
		return true
	}

	return false
}

func (o CertificateRead) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Href != nil {
		toSerialize["href"] = o.Href
	}

	if o.Metadata != nil {
		toSerialize["metadata"] = o.Metadata
	}

	if o.Properties != nil {
		toSerialize["properties"] = o.Properties
	}

	return json.Marshal(toSerialize)
}

type NullableCertificateRead struct {
	value *CertificateRead
	isSet bool
}

func (v NullableCertificateRead) Get() *CertificateRead {
	return v.value
}

func (v *NullableCertificateRead) Set(val *CertificateRead) {
	v.value = val
	v.isSet = true
}

func (v NullableCertificateRead) IsSet() bool {
	return v.isSet
}

func (v *NullableCertificateRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertificateRead(val *CertificateRead) *NullableCertificateRead {
	return &NullableCertificateRead{value: val, isSet: true}
}

func (v NullableCertificateRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertificateRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}