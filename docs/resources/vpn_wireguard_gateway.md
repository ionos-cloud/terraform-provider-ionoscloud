---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-server"
description: |-
  Creates and manages IonosCloud VPN Wireguard Gateway objects.
---

# VPN WireGuard Gateway Resource

## Overview

The `ionoscloud_vpn_wireguard_peer` resource manages a WireGuard Gateway within the IONOS Cloud infrastructure. 
This resource facilitates the creation, management, and deletion of WireGuard VPN Gateways, enabling secure connections between your network resources.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "es/vit"
}
resource "ionoscloud_ipblock" "ipblock_example" {
  location = "es/vit"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource ionoscloud_vpn_wireguard_gateway "gateway" {
  name = "my vpn test gateway"
  description = "description"
  private_key = "private"

  gateway_ip = ionoscloud_ipblock.ipblock_example.ips[0]
  interface_ipv4_cidr =  "192.168.1.100/24"
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.108/24"
  }
}
```

## Argument Reference

The following arguments are supported by the `vpn_wireguard_gateway` resource:

- `name` - (Required)[String] The name of the WireGuard Gateway.
- `description` - (Optional)[String] A description of the WireGuard Gateway.
- `endpoint` - (Optional, Block) The endpoint configuration for the WireGuard Gateway. This block supports fields documented below.
- `private_key` - (Required)[String] The private key for the WireGuard Gateway. To be created with the wg utility.
- `gateway_ip` - (Required)[String] The IP address of the WireGuard Gateway.
- `interface_ipv4_cidr` - (Optional)[String] The IPv4 CIDR for the WireGuard Gateway interface.
- `interface_ipv6_cidr` - (Optional)[String] The IPv6 CIDR for the WireGuard Gateway interface.
- `connections` - (Optional)[Block] The connection configuration for the WireGuard Gateway. This block supports fields documented below.
  - `datacenter_id` - (Required)[String] The ID of the datacenter where the WireGuard Gateway is located.
  - `lan_id` - (Required)[String] The ID of the LAN where the WireGuard Gateway is connected.
  - `ipv4_cidr` - (Required)[String] The IPv4 CIDR for the WireGuard Gateway connection.
  - `ipv6_cidr` - (Optional)[String] The IPv6 CIDR for the WireGuard Gateway connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `status` - (Computed)[String] The current status of the WireGuard Gateway.
- `public_key` - (Computed)[String] The public key for the WireGuard Gateway.

## Import

WireGuard Gateways can be imported using their ID:

```shell
terraform import ionoscloud_vpn_wireguard_gateway.example_gateway <id>
```