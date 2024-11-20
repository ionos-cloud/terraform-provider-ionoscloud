---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: vpn_wireguard_gateway_peer"
sidebar_current: "docs-ionoscloud-data-source-vpn-wireguard-peer"
description: |-
  Provides information about a specific IonosCloud VPN WireGuard Gateway.
---

# Data Source: ionoscloud_vpn_wireguard_gateway

The `ionoscloud_vpn_wireguard_gateway` data source provides information about a specific IonosCloud VPN WireGuard Gateway. You can use this data source to retrieve details of a WireGuard Gateway for use in other resources and configurations.

## Example Usage

```hcl
data "ionoscloud_vpn_wireguard_peer" "example" {
  location = "de/fra"
  gateway_id = "example-gateway"
  name = "example-peer"
}

output "vpn_wireguard_peer_public_key" {
  value = data.vpn_wireguard_peer.example.public_key
}
```

## Argument Reference

The following arguments are supported:

- `gateway_id` - (Required)[String] The ID of the WireGuard Gateway.
- `location` - (Optional)[String] The location of the WireGuard Gateway.
- `name` - (Optional)[String] The name of the WireGuard Peer.
- `id` - (Optional)[String] The ID of the WireGuard Peer.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The unique ID of the WireGuard Peer.
- `name` - The name of the WireGuard Peer.
- `description` - The description of the WireGuard Peer.
- `public_key` - WireGuard public key of the connecting peer.
- `status` - The current status of the WireGuard Peer.
- `endpoint` - The endpoint of the WireGuard Peer.
  - `host` - Hostname or IPV4 address that the WireGuard Server will connect to.
  - `port` - Port that the WireGuard Server will connect to. Default: 51820
- `allowed_ips` -  The subnet CIDRs that are allowed to connect to the WireGuard Gateway.
