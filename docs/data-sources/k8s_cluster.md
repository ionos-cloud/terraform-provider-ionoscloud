---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_k8s_cluster"
sidebar_current: "docs-ionoscloud-datasource-k8s-cluster"
description: |-
Get information on a IonosCloud K8s Cluster
---

# ionoscloud\_k8s\_cluster

The k8s cluster data source can be used to search for and return existing k8s clusters.

## Example Usage

```hcl
data "ionoscloud_k8s_cluster" "k8s_cluster_example" {
  name     = "My_Cluster"
}
```

## Argument Reference

* `name` - (Optional) Name or part of the name of an existing cluster that you want to search for.
* `id` - (Optional) ID of the cluster you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - id of the cluster
* `name` - name of the cluster
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
* `k8s_version` - Kubernetes version 
* `node_pools` - list of the IDs of the node pools in this cluster
* `kube_config` - Kubernetes configuration
