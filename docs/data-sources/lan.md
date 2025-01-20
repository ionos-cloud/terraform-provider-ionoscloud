---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_lan"
sidebar_current: "docs-ionoscloud-datasource-lan"
description: |-
  Get information on a Ionos Cloud Lans
---

# ionoscloud_lan

The **LAN data source** can be used to search for and return existing lans.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_lan" "example" {
  datacenter_id = "datacenter_id"
  id			= "lan_id"
}
```

### By Name
```hcl
data "ionoscloud_lan" "example" {
  datacenter_id = "datacenter_id"
  name			= "Lan Example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing lan that you want to search for.
* `id` - (Optional) ID of the lan you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the LAN.
* `name` - The name of the LAN.
* `datacenter_id` - The ID of lan's Virtual Data Center.
* `ip_failover` - list of
    * `nic_uuid`
    * `ip`
* `pcc` - The unique id of a `ionoscloud_private_crossconnect` resource, in order.
* `public` - Indicates if the LAN faces the public Internet (true) or not (false).
* `ipv4_cidr_block` - For public LANs this property is null, for private LANs it contains the private IPv4 CIDR range.
* `ipv6_cidr_block` - Contains the LAN's /64 IPv6 CIDR block if this LAN is IPv6 enabled.
