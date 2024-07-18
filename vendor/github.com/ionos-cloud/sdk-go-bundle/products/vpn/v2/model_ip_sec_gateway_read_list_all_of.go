/*
 * VPN Gateways
 *
 * POC Docs for VPN gateway as service
 *
 * API version: 0.0.1
 * Contact: support@cloud.ionos.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// checks if the IPSecGatewayReadListAllOf type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecGatewayReadListAllOf{}

// IPSecGatewayReadListAllOf struct for IPSecGatewayReadListAllOf
type IPSecGatewayReadListAllOf struct {
	// ID of the list of IPSecGateway resources.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the list of IPSecGateway resources.
	Href string `json:"href"`
	// The list of IPSecGateway resources.
	Items []IPSecGatewayRead `json:"items,omitempty"`
}

// NewIPSecGatewayReadListAllOf instantiates a new IPSecGatewayReadListAllOf object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecGatewayReadListAllOf(id string, type_ string, href string) *IPSecGatewayReadListAllOf {
	this := IPSecGatewayReadListAllOf{}

	this.Id = id
	this.Type = type_
	this.Href = href

	return &this
}

// NewIPSecGatewayReadListAllOfWithDefaults instantiates a new IPSecGatewayReadListAllOf object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecGatewayReadListAllOfWithDefaults() *IPSecGatewayReadListAllOf {
	this := IPSecGatewayReadListAllOf{}
	return &this
}

// GetId returns the Id field value
func (o *IPSecGatewayReadListAllOf) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayReadListAllOf) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *IPSecGatewayReadListAllOf) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *IPSecGatewayReadListAllOf) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayReadListAllOf) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *IPSecGatewayReadListAllOf) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *IPSecGatewayReadListAllOf) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *IPSecGatewayReadListAllOf) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *IPSecGatewayReadListAllOf) SetHref(v string) {
	o.Href = v
}

// GetItems returns the Items field value if set, zero value otherwise.
func (o *IPSecGatewayReadListAllOf) GetItems() []IPSecGatewayRead {
	if o == nil || IsNil(o.Items) {
		var ret []IPSecGatewayRead
		return ret
	}
	return o.Items
}

// GetItemsOk returns a tuple with the Items field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecGatewayReadListAllOf) GetItemsOk() ([]IPSecGatewayRead, bool) {
	if o == nil || IsNil(o.Items) {
		return nil, false
	}
	return o.Items, true
}

// HasItems returns a boolean if a field has been set.
func (o *IPSecGatewayReadListAllOf) HasItems() bool {
	if o != nil && !IsNil(o.Items) {
		return true
	}

	return false
}

// SetItems gets a reference to the given []IPSecGatewayRead and assigns it to the Items field.
func (o *IPSecGatewayReadListAllOf) SetItems(v []IPSecGatewayRead) {
	o.Items = v
}

func (o IPSecGatewayReadListAllOf) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecGatewayReadListAllOf) ToMap() (map[string]interface{}, error) {
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
	if !IsNil(o.Items) {
		toSerialize["items"] = o.Items
	}
	return toSerialize, nil
}

type NullableIPSecGatewayReadListAllOf struct {
	value *IPSecGatewayReadListAllOf
	isSet bool
}

func (v NullableIPSecGatewayReadListAllOf) Get() *IPSecGatewayReadListAllOf {
	return v.value
}

func (v *NullableIPSecGatewayReadListAllOf) Set(val *IPSecGatewayReadListAllOf) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecGatewayReadListAllOf) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecGatewayReadListAllOf) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecGatewayReadListAllOf(val *IPSecGatewayReadListAllOf) *NullableIPSecGatewayReadListAllOf {
	return &NullableIPSecGatewayReadListAllOf{value: val, isSet: true}
}

func (v NullableIPSecGatewayReadListAllOf) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecGatewayReadListAllOf) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
