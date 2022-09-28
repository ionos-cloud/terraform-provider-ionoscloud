---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipblock"
sidebar_current: "docs-datasource-ipblock"
description: |-
  Get information on a IonosCloud Ip Block
---

# ionoscloud\_ipblock

The **IP Block data source** can be used to search for and return an existing Ip Block.
You can provide a string for the id, the name or the location parameters which will be compared with the provisioned Ip Blocks.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_ipblock" "example" {
  id = <ipblock_id>
}
``` 

### By Name
```hcl
data "ionoscloud_ipblock" "example" {
  name = "IP Block Example"
}
``` 

### By Name with Partial Match
```hcl
data "ionoscloud_ipblock" "example" {
  name          = "Example"
  partial_match = true
}
``` 

### By Location
```hcl
data "ionoscloud_ipblock" "example" {
  location = "us/las"
}
``` 

### By Name & Location
``` 
data "ionoscloud_ipblock" "example" {
  name      = "IP Block Name"
  location  = "us/las"
}
```

## Argument reference

* `id` - (Optional) ID of an existing Ip Block that you want to search for.
* `name` - (Optional) Name of an existing Ip Block that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.
* `location` - (Optional) The regional location for this IP Block: us/las, us/ewr, de/fra, de/fkb.

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