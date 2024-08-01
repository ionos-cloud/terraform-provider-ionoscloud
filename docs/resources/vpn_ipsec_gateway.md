---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_vpn_ipsec_gateway"
sidebar_current: "docs-resource-vpn-ipsec-gateway"
description: |-
  IPSec Gateway
---

# ionoscloud_vpn_ipsec_gateway

An IPSec Gateway resource manages the creation, management, and deletion of VPN IPSec Gateways within the IONOS Cloud
infrastructure. This resource facilitates the creation of VPN IPSec Gateways, enabling secure connections between your
network resources.

## Usage example

```
resource "ionoscloud_vpn_ipsec_gateway" "example" {
	name = "example-gateway"
	location = <gateway_location>
	gateway_ip = <gateway_public_ip>
	version = "IKEv2"
	description = "This gateway connects site A to VDC X."

	connections {
		datacenter_id = <gateway_datacenter_id>
		lan_id = <lan_id>
		ipv4_cidr = <lan_ipv4_cidr>
		ipv6_cidr = <lan_ipv6_cidr>
	}
}
```

## Argument reference

* `name` - (Required)[string] The name of the IPSec Gateway.
* `location` - (Required)[string] The location of the IPSec Gateway. Supported locations: de/fra, de/txl, es/vit,
  gb/lhr, us/ewr, us/las, us/mci, fr/par
* `gateway_ip` - (Required)[string] Public IP address to be assigned to the gateway. Note: This must be an IP address in
  the same datacenter as the connections.
* `description` - (Optional)[string] The human-readable description of the IPSec Gateway.
* `connections` - (Required)[list] The network connection for your gateway. **Note**: all connections must belong to the
  same datacenter. Minimum items: 1. Maximum items: 10.
    * `datacenter_id` - (Required)[string] The datacenter to connect your VPN Gateway to.
    * `lan_id` - (Required)[string] The numeric LAN ID to connect your VPN Gateway to.
    * `ipv4_cidr` - (Required)[string] Describes the private ipv4 subnet in your LAN that should be accessible by the
      VPN Gateway. Note: this should be the subnet already assigned to the LAN
    * `ipv6_cidr` - (Optional)[string] Describes the ipv6 subnet in your LAN that should be accessible by the VPN
      Gateway. **Note**: this should be the subnet already assigned to the LAN
* `version` - (Required)[string] The IKE version that is permitted for the VPN tunnels. Default: `IKEv2`. Possible
  values: `IKEv2`.

## Import

The resource can be imported using the `location` and `gateway_id`, for example:

```
terraform import ionoscloud_vpn_ipsec_gateway.example {location}:{gateway_id}
```
