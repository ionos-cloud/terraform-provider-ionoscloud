---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_k8s_node_pool"
sidebar_current: "docs-ionoscloud-datasource-k8s-node-pool"
description: |-
Get information on a IonosCloud K8s Node Pool
---

# ionoscloud\_k8s\_node\_pool

The k8s node pool data source can be used to search for and return existing k8s node pools.

## Example Usage

```hcl
data "ionoscloud_k8s_node_pool" "k8s_node_pool_example" {
  name     = "My_Node_Pool"
  k8s_cluster_id 	= ionoscloud_k8s_cluster.k8s_cluster_example.id
}
```

## Argument Reference

* `k8s_cluster_id` (Required) K8s Cluster' UUID
* `name` - (Optional) Name of an existing node pool that you want to search for.
* `id` - (Optional) ID of the node pool you want to search for.

`k8s_cluster_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - id of the node pool
* `name` - name of the node pool
* `k8s_cluster_id` - ID of the cluster this node pool is part of
* `datacenter_id` - The UUID of the VDC
* `state` - one of "AVAILABLE",
  "INACTIVE",
  "BUSY",
  "DEPLOYING",
  "ACTIVE",
  "FAILED",
  "SUSPENDED",
  "FAILED_SUSPENDED",
  "UPDATING",
  "FAILED_UPDATING",
  "DESTROYING",
  "FAILED_DESTROYING",
  "TERMINATED"
* `node_count` - The number of nodes in this node pool
* `cpu_family` - CPU Family
* `cores_count` - CPU cores count
* `ram_size` - The amount of RAM in MB
* `availability_zone` - The compute availability zone in which the nodes should exist
* `storage_type` - HDD or SDD
* `storage_size` - Size of storage
* `k8s_version` - The kubernetes version
* `maintenance_window` - A maintenance window comprise of a day of the week and a time for maintenance to be allowed
  * `time` - A clock time in the day when maintenance is allowed
  * `day_of_the_week` - Day of the week when maintenance is allowed
* `auto_scaling` - The range defining the minimum and maximum number of worker nodes that the managed node group can scale in
  * `min_node_count` - The minimum number of worker nodes the node pool can scale down to
  * `max_node_count` - The maximum number of worker nodes that the node pool can scale to
* `lans` - A list of Local Area Networks the node pool is a part of
* `labels` - A map of labels in the form of key -> value
* `annotations` - A map of annotations in the form of key -> value
* `available_upgrade_versions` - A list of kubernetes versions available for upgrade
* `public_ips` - The list of fixed IPs associated with this node pool