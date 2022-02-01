---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: nic"
sidebar_current: "docs-ionoscloud-datasource-nic"
description: |-
  Get information on a Ionos Cloud NIC
---

# ionoscloud_nic

The nic data source can be used to search for and return existing nics.

## Example Usage

```hcl
data "ionoscloud_nic" "lan_example" {
  datacenter_id = ionoscloud_datacenter.example.id
  server_id = ionoscloud_server.example.id
  id			= "nic_id"
}

data "ionoscloud_nic" "lan_example" {
  datacenter_id = ionoscloud_datacenter.example.id
  server_id = ionoscloud_server.example.id
  name			= "nic_name"
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `server_id` - (Required)[string] The ID of a server.
- `name` - (Optional)[string] The name of the LAN.
* `id` - (Optional) ID of the nic you want to search for.

`datacenter_id` and either `name` or `id` must be provided.
If none, are provided, the datasource will return an error.

## Import

Resource Nic can be imported using the `resource id`, e.g.

## Attributes Reference

The following attributes are returned by the datasource:
* `id` - The id of the NIC.
* `datacenter_id` - The ID of a Virtual Data Center.
* `server_id` - The ID of a server.
* `lan` - The LAN ID the NIC will sit on.
* `name` - The name of the LAN.
* `dhcp` - Indicates if the NIC should get an IP address using DHCP (true) or not (false).
* `ips` - Collection of IP addresses assigned to a nic. Explicitly assigned public IPs need to come from reserved IP blocks, Passing value null or empty array will assign an IP address automatically.
* `firewall_active` - If this resource is set to true and is nested under a server resource firewall, with open SSH port, resource must be nested under the NIC.
* `mac` - The MAC address of the NIC.
* `nat` - Boolean value indicating if the private IP address has outbound access to the public internet.
