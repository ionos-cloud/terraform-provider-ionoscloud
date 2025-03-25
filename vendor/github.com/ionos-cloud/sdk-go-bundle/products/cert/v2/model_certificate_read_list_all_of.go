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

// checks if the CertificateReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &CertificateReadListAllOf{}

// CertificateReadListAllOf struct for CertificateReadListAllOf
type CertificateReadListAllOf struct {
	// ID of the list of Certificate resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of Certificate resources.
	Href string `json:"href"`
	// The list of Certificate resources.
	Items []CertificateRead `json:"items,omitempty"`
}

// NewCertificateReadListAllOf instantiates a new CertificateReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertificateReadListAllOf(id string, type_ string, href string) *CertificateReadListAllOf {
	this := CertificateReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewCertificateReadListAllOfWithDefaults instantiates a new CertificateReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertificateReadListAllOfWithDefaults() *CertificateReadListAllOf {
	this := CertificateReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *CertificateReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *CertificateReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *CertificateReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *CertificateReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *CertificateReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *CertificateReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *CertificateReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *CertificateReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *CertificateReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *CertificateReadListAllOf) GetItems() []CertificateRead {
	if o == nil || IsNil(o.Items) {
		var ret []CertificateRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *CertificateReadListAllOf) GetItemsOk() ([]CertificateRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *CertificateReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []CertificateRead and assigns it to the Items field.
func (o *CertificateReadListAllOf) SetItems(v []CertificateRead) {
	o.Items = v
}

func (o CertificateReadListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o CertificateReadListAllOf) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableCertificateReadListAllOf struct {
	value *CertificateReadListAllOf
	isSet bool
}

func (v NullableCertificateReadListAllOf) Get() *CertificateReadListAllOf {
	return v.value
}

func (v *NullableCertificateReadListAllOf) Set(val *CertificateReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableCertificateReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableCertificateReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertificateReadListAllOf(val *CertificateReadListAllOf) *NullableCertificateReadListAllOf {
	return &NullableCertificateReadListAllOf{value: val, isSet: true}
}

func (v NullableCertificateReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertificateReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
