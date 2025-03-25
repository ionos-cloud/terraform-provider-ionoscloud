/*
 * Certificate Manager Service API
 *
 * Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates  with IONOS services and your internal connected resources.   For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic. The service provides the basic functions of uploading and deleting your certificates for this purpose.
 *
 * API version: 2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cert

import (
	"encoding/json"
)

// checks if the MetadataWithAutoCertificateInformationAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MetadataWithAutoCertificateInformationAllOf{}

// MetadataWithAutoCertificateInformationAllOf struct for MetadataWithAutoCertificateInformationAllOf
type MetadataWithAutoCertificateInformationAllOf struct {
	// The ID of the last certificate that was issued.
	LastIssuedCertificate *string `json:"lastIssuedCertificate,omitempty"`
}

// NewMetadataWithAutoCertificateInformationAllOf instantiates a new MetadataWithAutoCertificateInformationAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMetadataWithAutoCertificateInformationAllOf() *MetadataWithAutoCertificateInformationAllOf {
	this := MetadataWithAutoCertificateInformationAllOf{}

	return &this
}

// NewMetadataWithAutoCertificateInformationAllOfWithDefaults instantiates a new MetadataWithAutoCertificateInformationAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMetadataWithAutoCertificateInformationAllOfWithDefaults() *MetadataWithAutoCertificateInformationAllOf {
	this := MetadataWithAutoCertificateInformationAllOf{}
	return &this
}

// GetLastIssuedCertificate returns the LastIssuedCertificate field value if set, zero value otherwise.
func (o *MetadataWithAutoCertificateInformationAllOf) GetLastIssuedCertificate() string {
	if o == nil || IsNil(o.LastIssuedCertificate) {
		var ret string
		return ret
	}
	return *o.LastIssuedCertificate
}

// GetLastIssuedCertificateOk returns a tuple with the LastIssuedCertificate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *MetadataWithAutoCertificateInformationAllOf) GetLastIssuedCertificateOk() (*string, bool) {
	if o == nil || IsNil(o.LastIssuedCertificate) {
		return nil, false
	}
	return o.LastIssuedCertificate, true
}

// HasLastIssuedCertificate returns a boolean if a field has been set.
func (o *MetadataWithAutoCertificateInformationAllOf) HasLastIssuedCertificate() bool {
	if o != nil && !IsNil(o.LastIssuedCertificate) {
		return true
	}

	return false
}

// SetLastIssuedCertificate gets a reference to the given string and assigns it to the LastIssuedCertificate field.
func (o *MetadataWithAutoCertificateInformationAllOf) SetLastIssuedCertificate(v string) {
	o.LastIssuedCertificate = &v
}

func (o MetadataWithAutoCertificateInformationAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MetadataWithAutoCertificateInformationAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.LastIssuedCertificate) {
		toSerialize["lastIssuedCertificate"] = o.LastIssuedCertificate
	}
	return toSerialize, nil
}

type NullableMetadataWithAutoCertificateInformationAllOf struct {
	value *MetadataWithAutoCertificateInformationAllOf
	isSet bool
}

func (v NullableMetadataWithAutoCertificateInformationAllOf) Get() *MetadataWithAutoCertificateInformationAllOf {
	return v.value
}

func (v *NullableMetadataWithAutoCertificateInformationAllOf) Set(val *MetadataWithAutoCertificateInformationAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataWithAutoCertificateInformationAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataWithAutoCertificateInformationAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataWithAutoCertificateInformationAllOf(val *MetadataWithAutoCertificateInformationAllOf) *NullableMetadataWithAutoCertificateInformationAllOf {
	return &NullableMetadataWithAutoCertificateInformationAllOf{value: val, isSet: true}
}

func (v NullableMetadataWithAutoCertificateInformationAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataWithAutoCertificateInformationAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
