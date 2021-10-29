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

* `id` - (Optional) ID of an existing Ip Block that you want to search for.
* `name` - (Optional) Name of an existing Ip Block that you want to search for.
* `location` - (Optional) ID of the existing Ip Block location.

## Attributes Reference
* `id` - The id of Ip Block
* `name` - The name of Ip Block
* `location` - The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.
* `size` - The number of IP addresses to reserve for this block.
* `ips` - The list of IP addresses associated with this block.
* `ip_consumers` Read-Only attribute. Lists consumption detail of an individual ip
    * `ip`
    * `mac`
    * `nic_uuid`
    * `server_id`
    * `server_name`
    * `datacenter_id`
    * `datacenter_name`
    * `k8s_nodepool_uuid`
    * `k8s_cluster_uuid`