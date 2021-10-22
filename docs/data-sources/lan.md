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
  name     = "My Lan"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing lan that you want to search for.
* `id` - (Optional) ID of the lan you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

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
