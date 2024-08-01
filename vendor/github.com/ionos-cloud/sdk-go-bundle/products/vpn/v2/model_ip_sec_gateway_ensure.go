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

// checks if the IPSecGatewayEnsure type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecGatewayEnsure{}

// IPSecGatewayEnsure struct for IPSecGatewayEnsure
type IPSecGatewayEnsure struct {
	// The ID (UUID) of the IPSecGateway.
	Id string `json:"id"`
	// Metadata
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Properties IPSecGateway           `json:"properties"`
}

// NewIPSecGatewayEnsure instantiates a new IPSecGatewayEnsure object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecGatewayEnsure(id string, properties IPSecGateway) *IPSecGatewayEnsure {
	this := IPSecGatewayEnsure{}

	this.Id = id
	this.Properties = properties

	return &this
}

// NewIPSecGatewayEnsureWithDefaults instantiates a new IPSecGatewayEnsure object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecGatewayEnsureWithDefaults() *IPSecGatewayEnsure {
	this := IPSecGatewayEnsure{}
	return &this
}

// GetId returns the Id field value
func (o *IPSecGatewayEnsure) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayEnsure) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *IPSecGatewayEnsure) SetId(v string) {
	o.Id = v
}

// GetMetadata returns the Metadata field value if set, zero value otherwise.
func (o *IPSecGatewayEnsure) GetMetadata() map[string]interface{} {
	if o == nil || IsNil(o.Metadata) {
		var ret map[string]interface{}
		return ret
	}
	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayEnsure) GetMetadataOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.Metadata) {
		return map[string]interface{}{}, false
	}
	return o.Metadata, true
}

// HasMetadata returns a boolean if a field has been set.
func (o *IPSecGatewayEnsure) HasMetadata() bool {
	if o != nil && !IsNil(o.Metadata) {
		return true
	}

	return false
}

// SetMetadata gets a reference to the given map[string]interface{} and assigns it to the Metadata field.
func (o *IPSecGatewayEnsure) SetMetadata(v map[string]interface{}) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *IPSecGatewayEnsure) GetProperties() IPSecGateway {
	if o == nil {
		var ret IPSecGateway
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayEnsure) GetPropertiesOk() (*IPSecGateway, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *IPSecGatewayEnsure) SetProperties(v IPSecGateway) {
	o.Properties = v
}

func (o IPSecGatewayEnsure) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecGatewayEnsure) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.Metadata) {
		toSerialize["metadata"] = o.Metadata
	}
	if !IsZero(o.Properties) {
		toSerialize["properties"] = o.Properties
	}
	return toSerialize, nil
}

type NullableIPSecGatewayEnsure struct {
	value *IPSecGatewayEnsure
	isSet bool
}

func (v NullableIPSecGatewayEnsure) Get() *IPSecGatewayEnsure {
	return v.value
}

func (v *NullableIPSecGatewayEnsure) Set(val *IPSecGatewayEnsure) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecGatewayEnsure) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecGatewayEnsure) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecGatewayEnsure(val *IPSecGatewayEnsure) *NullableIPSecGatewayEnsure {
	return &NullableIPSecGatewayEnsure{value: val, isSet: true}
}

func (v NullableIPSecGatewayEnsure) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecGatewayEnsure) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
