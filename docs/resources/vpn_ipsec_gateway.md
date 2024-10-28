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

```hcl
# Basic example

resource "ionoscloud_datacenter" "test_datacenter" {
  name = "test_vpn_gateway_basic"
  location = "de/fra"
}

resource "ionoscloud_lan" "test_lan" {
  name = "test_lan_basic"
  public = false
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
}

resource "ionoscloud_ipblock" "test_ipblock" {
  name = "test_ipblock_basic"
  location = "de/fra"
  size = 1
}

resource "ionoscloud_vpn_ipsec_gateway" "example" {
  name = "ipsec_gateway_basic"
  location = "de/fra"
  gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
  version = "IKEv2"
  description = "This gateway connects site A to VDC X."

  connections {
    datacenter_id = ionoscloud_datacenter.test_datacenter.id
    lan_id = ionoscloud_lan.test_lan.id
    ipv4_cidr = "192.168.100.10/24"
  }
}
```

```hcl
# Complete example

resource "ionoscloud_datacenter" "test_datacenter" {
  name = "vpn_gateway_test"
  location = "de/fra"
}

resource "ionoscloud_lan" "test_lan" {
  name = "test_lan"
  public = false
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  ipv6_cidr_block = local.lan_ipv6_cidr_block
}

resource "ionoscloud_ipblock" "test_ipblock" {
  name = "test_ipblock"
  location = "de/fra"
  size = 1
}

resource "ionoscloud_server" "test_server" {
  name = "test_server"
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  cores = 1
  ram = 2048
  image_name = "ubuntu:latest"
  image_password = random_password.server_image_password.result

  nic {
    lan = ionoscloud_lan.test_lan.id
    name = "test_nic"
    dhcp = true
    dhcpv6 = false
    ipv6_cidr_block = local.ipv6_cidr_block
    firewall_active   = false
  }

  volume {
    name         = "test_volume"
    disk_type    = "HDD"
    size         = 10
    licence_type = "OTHER"
  }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}

locals {
  lan_ipv6_cidr_block_parts = split("/", ionoscloud_datacenter.test_datacenter.ipv6_cidr_block)
  lan_ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "64")

  ipv4_cidr_block = format("%s/%s", ionoscloud_server.test_server.nic[0].ips[0], "24")
  ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "80")
}

resource "ionoscloud_vpn_ipsec_gateway" "example" {
	name = "ipsec-gateway"
	location = "de/fra"
	gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
	version = "IKEv2"
	description = "This gateway connects site A to VDC X."

	connections {
		datacenter_id = ionoscloud_datacenter.test_datacenter.id
		lan_id = ionoscloud_lan.test_lan.id
		ipv4_cidr = local.ipv4_cidr_block
		ipv6_cidr = local.ipv6_cidr_block
	}
    maintenance_window {
        day_of_the_week       = "Monday"
        time                  = "09:00:00"
    }
    tier = "STANDARD"
}
```

## Argument reference

* `name` - (Required)[string] The name of the IPSec Gateway.
* `location` - (Required)[string] The location of the IPSec Gateway. Supported locations: de/fra, de/txl, es/vit,
  gb/bhx, gb/lhr, us/ewr, us/las, us/mci, fr/par.
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
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `tier` - (Optional)(Computed)[string] Gateway performance options.  See product documentation for full details. Options: STANDARD, STANDARD_HA, ENHANCED, ENHANCED_HA, PREMIUM, PREMIUM_HA.

## Import

The resource can be imported using the `location` and `gateway_id`, for example:

```
terraform import ionoscloud_vpn_ipsec_gateway.example {location}:{gateway_id}
```
