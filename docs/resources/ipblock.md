---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ipblock"
sidebar_current: "docs-resource-ipblock"
description: |-
  Creates and manages IP Block objects.
---

# ionoscloud_ipblock

Manages **IP Blocks** on IONOS CLOUD. IP Blocks contain reserved public IP addresses that can be assigned servers or other resources.

## Example Usage

```hcl
resource "ionoscloud_ipblock" "example" {
  location  = "us/las"
  size      = 1
  name      = "IP Block Example"
}
```

## Argument reference

* `name` - (Optional)[string] The name of Ip Block
* `location` - (Required)[string] The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.
* `size` - (Required)[integer] The number of IP addresses to reserve for this block.
* `ips` - (Computed)[list of string] The list of IP addresses associated with this block.
* `ip_consumers` (Computed) Read-Only attribute. Lists consumption detail of an individual ip
  * `ip`
  * `mac`
  * `nic_id`
  * `server_id`
  * `server_name`
  * `datacenter_id`
  * `datacenter_name`
  * `k8s_nodepool_uuid`
  * `k8s_cluster_uuid`
  
## Import

Resource Ipblock can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_ipblock.myipblock ipblock uuid
```

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can also be used with the `identity` attribute:

```hcl
import {
  to = ionoscloud_ipblock.example
  identity = {
    id = "ipblock uuid"
  }
}

resource "ionoscloud_ipblock" "example" {
  ### Configuration omitted for brevity ###
}
```

### Identity Schema

#### Required

* `id` (String) The UUID of the IP block.

#### Optional

* `location` (String) The location the IP block lives in (e.g. `us/las`). Only needed when the Cloud API endpoint is overridden per location.

## Query (List Resource)

IP blocks can be listed using `terraform query` (requires Terraform 1.14+). List blocks must be placed in a dedicated query file, whose name ends in `.tfquery.hcl` (for example `queries.tfquery.hcl`).

```hcl
list "ionoscloud_ipblock" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

See the [`ionoscloud_ipblock` list resource documentation](../list-resources/ipblock.md) for filters and the full attribute reference.
