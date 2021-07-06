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

* `name` - (Optional) Name or an existing cluster that you want to search for.
* `id` - (Optional) ID of the cluster you want to search for.
k
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
* `maintenance_window`
* `available_upgrade_versions` - list of available versions for upgrading the cluster
* `viable_node_pool_versions` - list of versions that may be used for node pools under this cluster
* `node_pools` - list of the IDs of the node pools in this cluster
* `kube_config` - Kubernetes configuration
* `public` - the indicator if the cluster is public or private
* `gateway_ip` - the IP address of the gateway used by the cluster
* `api_subnet_allow_list` - access to the K8s API server is restricted to these CIDRs
* `s3_buckets` - list of S3 bucket configured for K8s usage