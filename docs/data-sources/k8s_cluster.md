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

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - id of the cluster
* `name` - name of the cluster
* `k8s_version` - Kubernetes version
* `maintenance_window` - A maintenance window comprise of a day of the week and a time for maintenance to be allowed
  * `time` - A clock time in the day when maintenance is allowed
  * `day_of_the_week` - Day of the week when maintenance is allowed
* `available_upgrade_versions` - A list of available versions for upgrading the cluster
* `viable_node_pool_versions` - A list of versions that may be used for node pools under this cluster
* `public` - The indicator if the cluster is public or private. Be aware that setting it to false is currently in beta phase
* `gateway_ip` - The IP address of the gateway used by the cluster. This is mandatory when `public` is set to `false` and should not be provided otherwise.
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
* `node_pools` - list of the IDs of the node pools in this cluster
* `kube_config` - Kubernetes configuration
* `api_subnet_allow_list` - access to the K8s API server is restricted to these CIDRs
* `s3_buckets` - list of S3 bucket configured for K8s usage