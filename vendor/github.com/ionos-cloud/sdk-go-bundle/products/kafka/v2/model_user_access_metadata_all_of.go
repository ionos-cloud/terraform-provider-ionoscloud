/*
 * Kafka as a Service API
 *
 * An managed Apache Kafka cluster is designed to be highly fault-tolerant and scalable, allowing large volumes of data to be ingested, stored, and processed in real-time. By distributing data across multiple brokers, Kafka achieves high throughput and low latency, making it suitable for applications requiring real-time data processing and analytics.
 *
 * API version: 1.7.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package kafka

import (
	"encoding/json"
)

// checks if the UserAccessMetadataAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UserAccessMetadataAllOf{}

// UserAccessMetadataAllOf struct for UserAccessMetadataAllOf
type UserAccessMetadataAllOf struct {
	// PEM for the certificate authority.
	CertificateAuthority *string `json:"certificateAuthority,omitempty"`
	// PEM for the private key.
	PrivateKey *string `json:"privateKey,omitempty"`
	// PEM for the certificate.
	Certificate *string `json:"certificate,omitempty"`
}

// NewUserAccessMetadataAllOf instantiates a new UserAccessMetadataAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUserAccessMetadataAllOf() *UserAccessMetadataAllOf {
	this := UserAccessMetadataAllOf{}

	return &this
}

// NewUserAccessMetadataAllOfWithDefaults instantiates a new UserAccessMetadataAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUserAccessMetadataAllOfWithDefaults() *UserAccessMetadataAllOf {
	this := UserAccessMetadataAllOf{}
	return &this
}

// GetCertificateAuthority returns the CertificateAuthority field value if set, zero value otherwise.
func (o *UserAccessMetadataAllOf) GetCertificateAuthority() string {
	if o == nil || IsNil(o.CertificateAuthority) {
		var ret string
		return ret
	}
	return *o.CertificateAuthority
}

// GetCertificateAuthorityOk returns a tuple with the CertificateAuthority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessMetadataAllOf) GetCertificateAuthorityOk() (*string, bool) {
	if o == nil || IsNil(o.CertificateAuthority) {
		return nil, false
	}
	return o.CertificateAuthority, true
}

// HasCertificateAuthority returns a boolean if a field has been set.
func (o *UserAccessMetadataAllOf) HasCertificateAuthority() bool {
	if o != nil && !IsNil(o.CertificateAuthority) {
		return true
	}

	return false
}

// SetCertificateAuthority gets a reference to the given string and assigns it to the CertificateAuthority field.
func (o *UserAccessMetadataAllOf) SetCertificateAuthority(v string) {
	o.CertificateAuthority = &v
}

// GetPrivateKey returns the PrivateKey field value if set, zero value otherwise.
func (o *UserAccessMetadataAllOf) GetPrivateKey() string {
	if o == nil || IsNil(o.PrivateKey) {
		var ret string
		return ret
	}
	return *o.PrivateKey
}

// GetPrivateKeyOk returns a tuple with the PrivateKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessMetadataAllOf) GetPrivateKeyOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateKey) {
		return nil, false
	}
	return o.PrivateKey, true
}

// HasPrivateKey returns a boolean if a field has been set.
func (o *UserAccessMetadataAllOf) HasPrivateKey() bool {
	if o != nil && !IsNil(o.PrivateKey) {
		return true
	}

	return false
}

// SetPrivateKey gets a reference to the given string and assigns it to the PrivateKey field.
func (o *UserAccessMetadataAllOf) SetPrivateKey(v string) {
	o.PrivateKey = &v
}

// GetCertificate returns the Certificate field value if set, zero value otherwise.
func (o *UserAccessMetadataAllOf) GetCertificate() string {
	if o == nil || IsNil(o.Certificate) {
		var ret string
		return ret
	}
	return *o.Certificate
}

// GetCertificateOk returns a tuple with the Certificate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UserAccessMetadataAllOf) GetCertificateOk() (*string, bool) {
	if o == nil || IsNil(o.Certificate) {
		return nil, false
	}
	return o.Certificate, true
}

// HasCertificate returns a boolean if a field has been set.
func (o *UserAccessMetadataAllOf) HasCertificate() bool {
	if o != nil && !IsNil(o.Certificate) {
		return true
	}

	return false
}

// SetCertificate gets a reference to the given string and assigns it to the Certificate field.
func (o *UserAccessMetadataAllOf) SetCertificate(v string) {
	o.Certificate = &v
}

func (o UserAccessMetadataAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UserAccessMetadataAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CertificateAuthority) {
		toSerialize["certificateAuthority"] = o.CertificateAuthority
	}
	if !IsNil(o.PrivateKey) {
		toSerialize["privateKey"] = o.PrivateKey
	}
	if !IsNil(o.Certificate) {
		toSerialize["certificate"] = o.Certificate
	}
	return toSerialize, nil
}

type NullableUserAccessMetadataAllOf struct {
	value *UserAccessMetadataAllOf
	isSet bool
}

func (v NullableUserAccessMetadataAllOf) Get() *UserAccessMetadataAllOf {
	return v.value
}

func (v *NullableUserAccessMetadataAllOf) Set(val *UserAccessMetadataAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableUserAccessMetadataAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableUserAccessMetadataAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUserAccessMetadataAllOf(val *UserAccessMetadataAllOf) *NullableUserAccessMetadataAllOf {
	return &NullableUserAccessMetadataAllOf{value: val, isSet: true}
}

func (v NullableUserAccessMetadataAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUserAccessMetadataAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
