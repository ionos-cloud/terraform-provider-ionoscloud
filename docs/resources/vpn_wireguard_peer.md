---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: vpn_wireguard_gateway - Data Source"
sidebar_current: "docs-ionoscloud-vpn-wireguard-peer"
description: |-
  Provides information about a specific IonosCloud VPN WireGuard Peer.
---
# VPN WireGuard Gateway Resource

This page provides an overview of the `ionoscloud_vpn_wireguard_peer` resource, which allows you to manage a WireGuard Peer in your cloud infrastructure. 
This resource enables the creation, management, and deletion of a WireGuard VPN Peer, facilitating secure connections between your network resources.

## Example Usage

```hcl
resource "ionoscloud_vpn_wireguard_peer" "example" {
  location = "de/fra"
  gateway_id  = "your gateway id here"
  name        = "example-gateway"
  description = "An example WireGuard peer"
  endpoint {
    host = "1.2.3.4"
    port = 51820
  }
  allowed_ips = ["10.0.0.0/8", "192.168.1.0/24"]
  public_key  = "examplePublicKey=="
}
```

## Argument Reference

The following arguments are supported:

- `gateway_id` - (Required)[string] The ID of the WireGuard Gateway that the Peer will connect to.
- `location` - (Required)[string] The location of the WireGuard Gateway. 
- `name` - (Required)[string] The human-readable name of the WireGuard Gateway.
- `public_key` - (Required)[string] The public key for the WireGuard Gateway.
- `description` - (Optional)[string] A description of the WireGuard Gateway.
- `allowed_ips` - (Required)[list, string] A list of subnet CIDRs that are allowed to connect to the WireGuard Gateway.
- `endpoint` - (Optional)[block] An endpoint configuration block for the WireGuard Gateway. The structure of this block is as follows:
  - `host` - (Required)[string] The hostname or IPV4 address that the WireGuard Server will connect to.
  - `port` - (Optional)[int] The port that the WireGuard Server will connect to. Defaults to `51820`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `status` - The current status of the WireGuard Gateway.

## Import

WireGuard Peers can be imported using the `gateway_id` and `id`, e.g.,

```shell
terraform import ionoscloud_vpn_wireguard_peer.example <gateway_id>:<peer_id>
```