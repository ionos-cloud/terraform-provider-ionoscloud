---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dataplatform_node_pools"
sidebar_current: "docs-dataplatform_node_pools"
description: |-
  Get information on a Dataplatform Node Pool list under a Dataplatform Cluster.
---

# ionoscloud\_dataplatform_node_pools

⚠️ **Note:** Data Platform is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

The **Dataplatform Node Pools Data Source** can be used to search for and return a list of existing Dataplatform Node Pools under a Dataplatform Cluster.

## Example Usage

### All Node Pools under a Cluster ID
```hcl
data "ionoscloud_dataplatform_node_pools" "example" {
  cluster_id  = <cluster_id>
}
```

### By Name

```hcl
data "ionoscloud_dataplatform_node_pools" "example" {
  cluster_id    = <cluster_id>
  name      	= "Dataplatform_Node_Pool_Example"
}
```

### By Name with Partial Match

```hcl
data "ionoscloud_dataplatform_node_pools" "example" {
  cluster_id    = <cluster_id>
  name    	    = "_Example"
  partial_match = true
}
```

## Argument Reference

* `cluster_id` - (Required) ID of the cluster the searched node pool is part of.
* `name` - (Optional) Name of an existing cluster that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_id` - ID of the cluster the searched node pool is part of.
* `node_pools` - List of Node Pools - See the [Node Pool](dataplatform_node_pool.md) section.