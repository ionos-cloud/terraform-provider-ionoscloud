---
subcategory: "Data Stack as a Service"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dsaas_node_pool"
sidebar_current: "docs-dsaas_node_pool"
description: |-
Get information on a DSaaS Cluster.
---

# ionoscloud\_pg_cluster

The **DSaaS Cluster data source** can be used to search for and return an existing DSaaS Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_dsaas_cluster" "example" {
  id	= <cluster_id>
}
```

### By Name

```hcl
data "ionoscloud_dsaas_cluster" "example" {
  display_name	= "DSaaS_Cluster_Example"
}
```

### By Name with Partial Match

```hcl
data "ionoscloud_dsaas_cluster" "example" {
  display_name	= "_Example"
  partial_match = true
}
```

## Argument Reference

* `id` - (Optional) ID of the cluster you want to search for.
* `name` - (Optional) Name or an existing cluster that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `display_name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the cluster.
* `datacenter_id` - The UUID of the virtual data center (VDC) the cluster is provisioned.
* `name` - The name of your cluster.
* `data_platform_version` - The version of the DataPlatform.
* `maintenance_window` - Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  * `time` - Time at which the maintenance should start. 
  * `day_of_the_week` - 
  