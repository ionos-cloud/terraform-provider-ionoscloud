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

// checks if the WireguardGateway type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &WireguardGateway{}

// WireguardGateway Properties with all data needed to create a new WireGuard Gateway.
type WireguardGateway struct {
	// The human readable name of your WireguardGateway.
	Name string `json:"name"`
	// Human readable description of the WireguardGateway.
	Description *string `json:"description,omitempty"`
	// Public IP address to be assigned to the gateway. __Note__: This must be an IP address in the same datacenter as the connections.
	GatewayIP string `json:"gatewayIP"`
	// Describes a range of IP V4 addresses in CIDR notation.
	InterfaceIPv4CIDR *string `json:"interfaceIPv4CIDR,omitempty"`
	// Describes a range of IP V6 addresses in CIDR notation.
	InterfaceIPv6CIDR *string `json:"interfaceIPv6CIDR,omitempty"`
	// The network connection for your gateway. __Note__: all connections must belong to the same datacenterId.
	Connections []Connection `json:"connections"`
	// PrivateKey used for WireGuard Server
	PrivateKey string `json:"privateKey"`
	// IP port number
	ListenPort *int32 `json:"listenPort,omitempty"`
}

// NewWireguardGateway instantiates a new WireguardGateway object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewWireguardGateway(name string, gatewayIP string, connections []Connection, privateKey string) *WireguardGateway {
	this := WireguardGateway{}

	this.Name = name
	this.GatewayIP = gatewayIP
	this.Connections = connections
	this.PrivateKey = privateKey
	var listenPort int32 = 51820
	this.ListenPort = &listenPort

	return &this
}

// NewWireguardGatewayWithDefaults instantiates a new WireguardGateway object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewWireguardGatewayWithDefaults() *WireguardGateway {
	this := WireguardGateway{}
	var listenPort int32 = 51820
	this.ListenPort = &listenPort
	return &this
}

// GetName returns the Name field value
func (o *WireguardGateway) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *WireguardGateway) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *WireguardGateway) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *WireguardGateway) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *WireguardGateway) SetDescription(v string) {
	o.Description = &v
}

// GetGatewayIP returns the GatewayIP field value
func (o *WireguardGateway) GetGatewayIP() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.GatewayIP
}

// GetGatewayIPOk returns a tuple with the GatewayIP field value
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetGatewayIPOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.GatewayIP, true
}

// SetGatewayIP sets field value
func (o *WireguardGateway) SetGatewayIP(v string) {
	o.GatewayIP = v
}

// GetInterfaceIPv4CIDR returns the InterfaceIPv4CIDR field value if set, zero value otherwise.
func (o *WireguardGateway) GetInterfaceIPv4CIDR() string {
	if o == nil || IsNil(o.InterfaceIPv4CIDR) {
		var ret string
		return ret
	}
	return *o.InterfaceIPv4CIDR
}

// GetInterfaceIPv4CIDROk returns a tuple with the InterfaceIPv4CIDR field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetInterfaceIPv4CIDROk() (*string, bool) {
	if o == nil || IsNil(o.InterfaceIPv4CIDR) {
		return nil, false
	}
	return o.InterfaceIPv4CIDR, true
}

// HasInterfaceIPv4CIDR returns a boolean if a field has been set.
func (o *WireguardGateway) HasInterfaceIPv4CIDR() bool {
	if o != nil && !IsNil(o.InterfaceIPv4CIDR) {
		return true
	}

	return false
}

// SetInterfaceIPv4CIDR gets a reference to the given string and assigns it to the InterfaceIPv4CIDR field.
func (o *WireguardGateway) SetInterfaceIPv4CIDR(v string) {
	o.InterfaceIPv4CIDR = &v
}

// GetInterfaceIPv6CIDR returns the InterfaceIPv6CIDR field value if set, zero value otherwise.
func (o *WireguardGateway) GetInterfaceIPv6CIDR() string {
	if o == nil || IsNil(o.InterfaceIPv6CIDR) {
		var ret string
		return ret
	}
	return *o.InterfaceIPv6CIDR
}

// GetInterfaceIPv6CIDROk returns a tuple with the InterfaceIPv6CIDR field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetInterfaceIPv6CIDROk() (*string, bool) {
	if o == nil || IsNil(o.InterfaceIPv6CIDR) {
		return nil, false
	}
	return o.InterfaceIPv6CIDR, true
}

// HasInterfaceIPv6CIDR returns a boolean if a field has been set.
func (o *WireguardGateway) HasInterfaceIPv6CIDR() bool {
	if o != nil && !IsNil(o.InterfaceIPv6CIDR) {
		return true
	}

	return false
}

// SetInterfaceIPv6CIDR gets a reference to the given string and assigns it to the InterfaceIPv6CIDR field.
func (o *WireguardGateway) SetInterfaceIPv6CIDR(v string) {
	o.InterfaceIPv6CIDR = &v
}

// GetConnections returns the Connections field value
func (o *WireguardGateway) GetConnections() []Connection {
	if o == nil {
		var ret []Connection
		return ret
	}

	return o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetConnectionsOk() ([]Connection, bool) {
	if o == nil {
		return nil, false
	}
	return o.Connections, true
}

// SetConnections sets field value
func (o *WireguardGateway) SetConnections(v []Connection) {
	o.Connections = v
}

// GetPrivateKey returns the PrivateKey field value
func (o *WireguardGateway) GetPrivateKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PrivateKey
}

// GetPrivateKeyOk returns a tuple with the PrivateKey field value
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetPrivateKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PrivateKey, true
}

// SetPrivateKey sets field value
func (o *WireguardGateway) SetPrivateKey(v string) {
	o.PrivateKey = v
}

// GetListenPort returns the ListenPort field value if set, zero value otherwise.
func (o *WireguardGateway) GetListenPort() int32 {
	if o == nil || IsNil(o.ListenPort) {
		var ret int32
		return ret
	}
	return *o.ListenPort
}

// GetListenPortOk returns a tuple with the ListenPort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *WireguardGateway) GetListenPortOk() (*int32, bool) {
	if o == nil || IsNil(o.ListenPort) {
		return nil, false
	}
	return o.ListenPort, true
}

// HasListenPort returns a boolean if a field has been set.
func (o *WireguardGateway) HasListenPort() bool {
	if o != nil && !IsNil(o.ListenPort) {
		return true
	}

	return false
}

// SetListenPort gets a reference to the given int32 and assigns it to the ListenPort field.
func (o *WireguardGateway) SetListenPort(v int32) {
	o.ListenPort = &v
}

func (o WireguardGateway) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o WireguardGateway) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsZero(o.GatewayIP) {
		toSerialize["gatewayIP"] = o.GatewayIP
	}
	if !IsNil(o.InterfaceIPv4CIDR) {
		toSerialize["interfaceIPv4CIDR"] = o.InterfaceIPv4CIDR
	}
	if !IsNil(o.InterfaceIPv6CIDR) {
		toSerialize["interfaceIPv6CIDR"] = o.InterfaceIPv6CIDR
	}
	if !IsZero(o.Connections) {
		toSerialize["connections"] = o.Connections
	}
	if !IsZero(o.PrivateKey) {
		toSerialize["privateKey"] = o.PrivateKey
	}
	if !IsNil(o.ListenPort) {
		toSerialize["listenPort"] = o.ListenPort
	}
	return toSerialize, nil
}

type NullableWireguardGateway struct {
	value *WireguardGateway
	isSet bool
}

func (v NullableWireguardGateway) Get() *WireguardGateway {
	return v.value
}

func (v *NullableWireguardGateway) Set(val *WireguardGateway) {
	v.value = val
	v.isSet = true
}

func (v NullableWireguardGateway) IsSet() bool {
	return v.isSet
}

func (v *NullableWireguardGateway) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableWireguardGateway(val *WireguardGateway) *NullableWireguardGateway {
	return &NullableWireguardGateway{value: val, isSet: true}
}

func (v NullableWireguardGateway) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableWireguardGateway) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
