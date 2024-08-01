/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package vpn

import (
	"encoding/json"
)

// checks if the IPSecGatewayRead type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecGatewayRead{}

// IPSecGatewayRead struct for IPSecGatewayRead
type IPSecGatewayRead struct {
	// The ID (UUID) of the IPSecGateway.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the IPSecGateway.
	Href       string               `json:"href"`
	Metadata   IPSecGatewayMetadata `json:"metadata"`
	Properties IPSecGateway         `json:"properties"`
}

// NewIPSecGatewayRead instantiates a new IPSecGatewayRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecGatewayRead(id string, type_ string, href string, metadata IPSecGatewayMetadata, properties IPSecGateway) *IPSecGatewayRead {
	this := IPSecGatewayRead{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewIPSecGatewayReadWithDefaults instantiates a new IPSecGatewayRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecGatewayReadWithDefaults() *IPSecGatewayRead {
	this := IPSecGatewayRead{}
	return &this
}

// GetId returns the Id field value
func (o *IPSecGatewayRead) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *IPSecGatewayRead) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *IPSecGatewayRead) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *IPSecGatewayRead) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *IPSecGatewayRead) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *IPSecGatewayRead) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *IPSecGatewayRead) GetMetadata() IPSecGatewayMetadata {
	if o == nil {
		var ret IPSecGatewayMetadata
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayRead) GetMetadataOk() (*IPSecGatewayMetadata, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *IPSecGatewayRead) SetMetadata(v IPSecGatewayMetadata) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *IPSecGatewayRead) GetProperties() IPSecGateway {
	if o == nil {
		var ret IPSecGateway
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayRead) GetPropertiesOk() (*IPSecGateway, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *IPSecGatewayRead) SetProperties(v IPSecGateway) {
	o.Properties = v
}

func (o IPSecGatewayRead) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecGatewayRead) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsZero(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsZero(o.Href) {
		toSerialize["href"] = o.Href
	}
	if !IsZero(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsZero(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableIPSecGatewayRead struct {
	value *IPSecGatewayRead
	isSet bool
}

func (v NullableIPSecGatewayRead) Get() *IPSecGatewayRead {
	return v.value
}

func (v *NullableIPSecGatewayRead) Set(val *IPSecGatewayRead) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecGatewayRead) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecGatewayRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecGatewayRead(val *IPSecGatewayRead) *NullableIPSecGatewayRead {
	return &NullableIPSecGatewayRead{value: val, isSet: true}
}

func (v NullableIPSecGatewayRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecGatewayRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
