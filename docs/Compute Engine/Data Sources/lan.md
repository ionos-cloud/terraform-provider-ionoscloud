---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_lan"
sidebar_current: "docs-ionoscloud-datasource-lan"
description: |-
Get information on a Ionos Cloud Lans
---

# ionoscloud\_lan

The lans data source can be used to search for and return existing lans.

## Example Usage

```hcl
data "ionoscloud_lan" "lan_example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "example_"
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
