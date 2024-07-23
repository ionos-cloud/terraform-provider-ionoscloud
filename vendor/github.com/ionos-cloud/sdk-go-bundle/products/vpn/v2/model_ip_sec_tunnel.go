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

// checks if the IPSecTunnel type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &IPSecTunnel{}

// IPSecTunnel Properties with all data needed to create a new IPSec Gateway Tunnel.\\ __Note__: there is a limit of 20 tunnels per IPSec Gateway.
type IPSecTunnel struct {
	// The human readable name of your IPSec Gateway Tunnel.
	Name string `json:"name"`
	// Human readable description of the IPSec Gateway Tunnel.
	Description *string `json:"description,omitempty"`
	// The remote peer host fully qualified domain name or IPV4 IP to connect to. * __Note__: This should be the public IP of the remote peer. * Tunnels only support IPV4 or hostname (fully qualified DNS name).
	RemoteHost string          `json:"remoteHost"`
	Auth       IPSecTunnelAuth `json:"auth"`
	Ike        IKEEncryption   `json:"ike"`
	Esp        ESPEncryption   `json:"esp"`
	// The network CIDRs on the \"Left\" side that are allowed to connect to the IPSec tunnel, i.e the CIDRs within your IONOS Cloud LAN.  Specify \"0.0.0.0/0\" or \"::/0\" for all addresses.
	CloudNetworkCIDRs []string `json:"cloudNetworkCIDRs"`
	// The network CIDRs on the \"Right\" side that are allowed to connect to the IPSec tunnel.  Specify \"0.0.0.0/0\" or \"::/0\" for all addresses.
	PeerNetworkCIDRs []string `json:"peerNetworkCIDRs"`
}

// NewIPSecTunnel instantiates a new IPSecTunnel object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewIPSecTunnel(name string, remoteHost string, auth IPSecTunnelAuth, ike IKEEncryption, esp ESPEncryption, cloudNetworkCIDRs []string, peerNetworkCIDRs []string) *IPSecTunnel {
	this := IPSecTunnel{}

	this.Name = name
	this.RemoteHost = remoteHost
	this.Auth = auth
	this.Ike = ike
	this.Esp = esp
	this.CloudNetworkCIDRs = cloudNetworkCIDRs
	this.PeerNetworkCIDRs = peerNetworkCIDRs

	return &this
}

// NewIPSecTunnelWithDefaults instantiates a new IPSecTunnel object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewIPSecTunnelWithDefaults() *IPSecTunnel {
	this := IPSecTunnel{}
	return &this
}

// GetName returns the Name field value
func (o *IPSecTunnel) GetName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Name
}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Name, true
}

// SetName sets field value
func (o *IPSecTunnel) SetName(v string) {
	o.Name = v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *IPSecTunnel) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *IPSecTunnel) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *IPSecTunnel) SetDescription(v string) {
	o.Description = &v
}

// GetRemoteHost returns the RemoteHost field value
func (o *IPSecTunnel) GetRemoteHost() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RemoteHost
}

// GetRemoteHostOk returns a tuple with the RemoteHost field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetRemoteHostOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RemoteHost, true
}

// SetRemoteHost sets field value
func (o *IPSecTunnel) SetRemoteHost(v string) {
	o.RemoteHost = v
}

// GetAuth returns the Auth field value
func (o *IPSecTunnel) GetAuth() IPSecTunnelAuth {
	if o == nil {
		var ret IPSecTunnelAuth
		return ret
	}

	return o.Auth
}

// GetAuthOk returns a tuple with the Auth field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetAuthOk() (*IPSecTunnelAuth, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Auth, true
}

// SetAuth sets field value
func (o *IPSecTunnel) SetAuth(v IPSecTunnelAuth) {
	o.Auth = v
}

// GetIke returns the Ike field value
func (o *IPSecTunnel) GetIke() IKEEncryption {
	if o == nil {
		var ret IKEEncryption
		return ret
	}

	return o.Ike
}

// GetIkeOk returns a tuple with the Ike field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetIkeOk() (*IKEEncryption, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ike, true
}

// SetIke sets field value
func (o *IPSecTunnel) SetIke(v IKEEncryption) {
	o.Ike = v
}

// GetEsp returns the Esp field value
func (o *IPSecTunnel) GetEsp() ESPEncryption {
	if o == nil {
		var ret ESPEncryption
		return ret
	}

	return o.Esp
}

// GetEspOk returns a tuple with the Esp field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetEspOk() (*ESPEncryption, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Esp, true
}

// SetEsp sets field value
func (o *IPSecTunnel) SetEsp(v ESPEncryption) {
	o.Esp = v
}

// GetCloudNetworkCIDRs returns the CloudNetworkCIDRs field value
func (o *IPSecTunnel) GetCloudNetworkCIDRs() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.CloudNetworkCIDRs
}

// GetCloudNetworkCIDRsOk returns a tuple with the CloudNetworkCIDRs field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetCloudNetworkCIDRsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.CloudNetworkCIDRs, true
}

// SetCloudNetworkCIDRs sets field value
func (o *IPSecTunnel) SetCloudNetworkCIDRs(v []string) {
	o.CloudNetworkCIDRs = v
}

// GetPeerNetworkCIDRs returns the PeerNetworkCIDRs field value
func (o *IPSecTunnel) GetPeerNetworkCIDRs() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.PeerNetworkCIDRs
}

// GetPeerNetworkCIDRsOk returns a tuple with the PeerNetworkCIDRs field value
// and a boolean to check if the value has been set.
func (o *IPSecTunnel) GetPeerNetworkCIDRsOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.PeerNetworkCIDRs, true
}

// SetPeerNetworkCIDRs sets field value
func (o *IPSecTunnel) SetPeerNetworkCIDRs(v []string) {
	o.PeerNetworkCIDRs = v
}

func (o IPSecTunnel) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o IPSecTunnel) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsZero(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsZero(o.RemoteHost) {
		toSerialize["remoteHost"] = o.RemoteHost
	}
	if !IsZero(o.Auth) {
		toSerialize["auth"] = o.Auth
	}
	if !IsZero(o.Ike) {
		toSerialize["ike"] = o.Ike
	}
	if !IsZero(o.Esp) {
		toSerialize["esp"] = o.Esp
	}
	if !IsZero(o.CloudNetworkCIDRs) {
		toSerialize["cloudNetworkCIDRs"] = o.CloudNetworkCIDRs
	}
	if !IsZero(o.PeerNetworkCIDRs) {
		toSerialize["peerNetworkCIDRs"] = o.PeerNetworkCIDRs
	}
	return toSerialize, nil
}

type NullableIPSecTunnel struct {
	value *IPSecTunnel
	isSet bool
}

func (v NullableIPSecTunnel) Get() *IPSecTunnel {
	return v.value
}

func (v *NullableIPSecTunnel) Set(val *IPSecTunnel) {
	v.value = val
	v.isSet = true
}

func (v NullableIPSecTunnel) IsSet() bool {
	return v.isSet
}

func (v *NullableIPSecTunnel) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableIPSecTunnel(val *IPSecTunnel) *NullableIPSecTunnel {
	return &NullableIPSecTunnel{value: val, isSet: true}
}

func (v NullableIPSecTunnel) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableIPSecTunnel) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
