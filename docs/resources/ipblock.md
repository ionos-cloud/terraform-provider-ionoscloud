---
layout: "ionoscloud"
page_title: "IonosCloud: ipblock"
sidebar_current: "docs-resource-ipblock"
description: |-
  Creates and manages IP Block objects.
---

# ionoscloud\_ipblock

Manages IP Blocks on IonosCloud. IP Blocks contain reserved public IP addresses that can be assigned servers or other resources.

## Example Usage

```hcl
resource "ionoscloud_ipblock" "example" {
  location = "${ionoscloud_datacenter.example.location}"
  size     = 1
}
```

## Argument reference

* `name` - (Optional)[string] The name of Ip Block
* `location` - (Required)[string] The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.
* `size` - (Required)[integer] The number of IP addresses to reserve for this block.
* `ips` - (Computed)[integer] The list of IP addresses associated with this block.

## Import

Resource Ipblock can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_ipblock.myipblock {ipblock uuid}
```
