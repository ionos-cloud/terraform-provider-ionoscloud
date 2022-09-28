---
subcategory: "Managed Kubernetes"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_k8s_node_pool"
sidebar_current: "docs-ionoscloud-datasource-k8s-node-pool"
description: |-
  Get information on a IonosCloud K8s Node Pool
---

# ionoscloud\_k8s\_node\_pool

The **k8s Node Pool** data source can be used to search for and return existing k8s Node Pools.
You can provide a string for either id or name parameters which will be compared with provisioned K8s Node Pools.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_k8s_node_pool" "example" {
  id                = <k8s_nodepool_id>
  k8s_cluster_id 	= <k8s_cluster_id>
}
```

### By Name
```hcl
data "ionoscloud_k8s_node_pool" "example" {
  name              = "k8s NodePool Example"
  k8s_cluster_id 	= <k8s_cluster_id>
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_k8s_node_pool" "example" {
  k8s_cluster_id 	= <k8s_cluster_id>
  name              = "Example"
  partial_match 	= true
}
```

## Argument Reference

* `k8s_cluster_id` (Required) K8s Cluster' UUID
* `id` - (Optional) ID of the node pool you want to search for.
* `name` - (Optional) Name of an existing node pool that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

`k8s_cluster_id` and either `id` or `name` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

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
* `storage_size` - The size of the volume in GB. The size should be greater than 10GB.
* `k8s_version` - The kubernetes version
* `maintenance_window` - A maintenance window comprise of a day of the week and a time for maintenance to be allowed
  * `time` - A clock time in the day when maintenance is allowed
  * `day_of_the_week` - Day of the week when maintenance is allowed
* `auto_scaling` - The range defining the minimum and maximum number of worker nodes that the managed node group can scale in
  * `min_node_count` - The minimum number of worker nodes the node pool can scale down to
  * `max_node_count` - The maximum number of worker nodes that the node pool can scale to
* `lans` - A list of Local Area Networks the node pool is a part of
  * `id` - The LAN ID of an existing LAN at the related datacenter
  * `dhcp` - Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP
  * `routes` - An array of additional LANs attached to worker nodes
    * `network` - IPv4 or IPv6 CIDR to be routed via the interface
    * `gateway_ip` - IPv4 or IPv6 Gateway IP for the route
* `labels` - A map of labels in the form of key -> value
* `annotations` - A map of annotations in the form of key -> value
* `available_upgrade_versions` - A list of kubernetes versions available for upgrade
* `public_ips` - The list of fixed IPs associated with this node pool
