---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: vpn_wireguard_gateway - Data Source"
sidebar_current: "docs-ionoscloud-data-source-vpn-wireguard-gateway"
description: |-
  Provides information about a specific IonosCloud VPN WireGuard Gateway.
---

# Data Source: ionoscloud_vpn_wireguard_gateway

The `ionoscloud_vpn_wireguard_gateway` data source provides information about a specific IonosCloud VPN WireGuard Gateway. You can use this data source to retrieve details of a WireGuard Gateway for use in other resources and configurations.

## Example Usage

```hcl
data "ionoscloud_vpn_wireguard_gateway" "example" {
  location = "de/fra"
  name = "example-gateway"
}

output "vpn_wireguard_gateway_public_key" {
  value = data.vpn_wireguard_gateway.example.public_key
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional)[String] The name of the WireGuard Gateway.
- `id` - (Optional)[String] The ID of the WireGuard Gateway.
- `location` - (Required)[String] The location of the WireGuard Gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `description` - The description of the WireGuard Gateway.
- `public_key` - The public key for the WireGuard Gateway.
- `status` - The current status of the WireGuard Gateway.
- `gateway_ip` - The IP address of the WireGuard Gateway.
- `interface_ipv4_cidr` - The IPv4 CIDR for the WireGuard Gateway interface.
- `interface_ipv6_cidr` - The IPv6 CIDR for the WireGuard Gateway interface.
- `connections` - A list of connection configurations for the WireGuard Gateway. Each `connections` block contains:
    - `datacenter_id` - The ID of the datacenter where the WireGuard Gateway is located.
    - `lan_id` - The ID of the LAN where the WireGuard Gateway is connected.
    - `ipv4_cidr` - The IPv4 CIDR for the WireGuard Gateway connection.
    - `ipv6_cidr` - The IPv6 CIDR for the WireGuard Gateway connection.
- `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur.
  - `time` - Start of the maintenance window in UTC time.
  - `day_of_the_week` - The name of the week day.
- `tier` - Gateway performance options.
