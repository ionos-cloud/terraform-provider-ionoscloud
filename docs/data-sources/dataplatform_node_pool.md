---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dataplatform_node_pool"
sidebar_current: "docs-dataplatform_node_pool"
description: |-
Get information on a Dataplatform Node Pool.
---

# ionoscloud\_dataplatform_node_pool

The **Dataplatform Node Pool Data Source** can be used to search for and return an existing Dataplatform Node Pool under a Dataplatform Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_dataplatform_node_pool" "example" {
  cluster_id  = <cluster_id>
  id	      = <node_pool_id>
}
```

### By Name

```hcl
data "ionoscloud_dataplatform_node_pool" "example" {
  cluster_id    = <cluster_id>
  name      	= "Dataplatform_Node_Pool_Example"
}
```

### By Name with Partial Match

```hcl
data "ionoscloud_dataplatform_node_pool" "example" {
  cluster_id    = <cluster_id>
  name      	= "_Example"
  partial_match = true
}
```

## Argument Reference

* `cluster_id` - (Required) ID of the cluster the searched node pool is part of.
* `id` - (Optional) ID of the node pool you want to search for.
* `name` - (Optional) Name or an existing cluster that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `display_name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_id` - ID of the cluster the searched node pool is part of.
* `datacenter_id` - The UUID of the virtual data center (VDC) the cluster is provisioned.
* `id` - ID of your node pool.
* `name` - The name of your node pool
* `data_platform_version` - The version of the DataPlatform.
* `node_count` - The number of nodes that make up the node pool.
* `cpu_family` - A CPU family.
* `cores_count` - The number of CPU cores per node. 
* `ram_size` - The RAM size for one node in MB. 
* `availability_zone` - The availability zone of the virtual datacenter region where the node pool resources should be provisioned. 
* `storage_type` - The type of hardware for the volume. 
* `storage_size` - The size of the volume in GB. 
* `maintenance_window` - Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  * `time` - Time at which the maintenance should start. 
  * `day_of_the_week` 
* `labels` - Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
* `annotations` - Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/).
