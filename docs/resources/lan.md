---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: lan"
sidebar_current: "docs-resource-lan"
description: |-
  Creates and manages LAN objects.
---

# ionoscloud\_lan

Manages a LAN on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = true
  pcc           = ionoscloud_private_crossconnect.example.id
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the LAN.
* `public` - (Optional)[Boolean] Indicates if the LAN faces the public Internet (true) or not (false).
* `pcc` - (Optional)[String] The unique id of a `ionoscloud_private_crossconnect` resource, in order
* `ip_failover` - (Computed) IP failover configurations for lan
  * `ip`
  * `nic_uuid`
  
## Import

Resource Lan can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_lan.mylan {datacenter uuid}/{lan id}
```

## Important Notes

- Please note that only LANS datacenters found in the same physical location can be connected through a private cross-connect
- A LAN cannot be a part of two private cross-connects