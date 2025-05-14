---
subcategory: "Managed Kubernetes"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_k8s_node_pool_nodes"
sidebar_current: "docs-ionoscloud-datasource-k8s-node-pool-nodes"
description: |-
  Get information on the list of IonosCloud K8s Nodes that make a nodepool
---

# ionoscloud_k8s_node_pool_nodes

The **k8s Node Pool Nodes** data source can be used to search for and return a list of existing k8s Node Pool nodes.
## Example Usage

### By IDs
```hcl
data "ionoscloud_k8s_node_pool_nodes" "example" {
  node_pool_id      = "k8s_nodepool_id"
  k8s_cluster_id 	= "k8s_cluster_id"
}
```


## Argument Reference

* `k8s_cluster_id` (Required) K8s Cluster' UUID
* `id` - (Optional) ID of the node pool you want to search for.

`k8s_cluster_id` and `node_pool_id` must be provided.

## Attributes Reference

The following attributes are returned by the datasource:
* `nodes` - a list of the nodes that are in the nodepool 
  * `id` - id of the node in the nodepool
  * `name` - name of the node
  * `k8s_version` - The kubernetes version
  * `public_ip` - public ip of the node. Only present if the k8s clsuter is public
  * `private_ip` - private ip of the node. Only present if the k8s cluster is private.