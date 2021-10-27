---
layout: "ionoscloud"
page_title: "IonosCloud: ipblock"
sidebar_current: "docs-datasource-ipblock"
description: |-
Get information on a IonosCloud Ip Block
---

# ionoscloud\_ipblock

The Ip Block data source can be used to search for and return an existing Ip Block which can 
then be used to provision a server.

## Example Usage

```hcl
datasource "ionoscloud_ipblock" "example" {
  id =` <ipblock_id>`
}
```

## Argument reference

* `id` - (Required)[string] ID of the ip block you want to search for.

## Attributes Reference
* `id` - (Optional)[string] The id of Ip Block
* `name` - (Optional)[string] The name of Ip Block
* `location` - (Required)[string] The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.
* `size` - (Required)[integer] The number of IP addresses to reserve for this block.
* `ips` - (Computed)[integer] The list of IP addresses associated with this block.
* `ip_consumers` (Computed) Read-Only attribute. Lists consumption detail of an individual ip
  * `ip`
  * `mac`
  * `nic_uuid`
  * `server_id`
  * `server_name`
  * `datacenter_id`
  * `datacenter_name`
  * `k8s_nodepool_uuid`
  * `k8s_cluster_uuid`