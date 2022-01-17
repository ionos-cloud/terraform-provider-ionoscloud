---
subcategory: "Managed Kubernetes"
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
* `api_subnet_allow_list` - access to the K8s API server is restricted to these CIDRs
* `s3_buckets` - list of S3 bucket configured for K8s usage
* `kube_config` - Kubernetes configuration
* `config` - structured kubernetes config consisting of a list with 1 item with the following fields:
  * api_version - Kubernetes API Version
  * kind - "Config"
  * current-context - string
  * clusters - list of
    * name - name of cluster
    * cluster - map of
      * certificate-authority-data - **base64 decoded** cluster CA data
      * server -  server address in the form `https://host:port`
  * contexts - list of
    * name - context name
    * context - map of
      * cluster - cluster name
      * user - cluster user
  * users - list of
    * name - user name
    * user - map of
      * token - user token used for authentication
* `user_tokens` - a convenience map to be search the token of a specific user
  - key - is the user name
  - value - is the token
* `server` - cluster server (same as `config[0].clusters[0].cluster.server` but provided as an attribute for ease of use)
* `ca_crt` - base64 decoded cluster certificate authority data (provided as an attribute for direct use)

**NOTE**: The whole `config` node is marked as **sensitive**.

## Example of accessing a kubernetes cluster using the user's token

```
resource "ionoscloud_k8s_cluster" "test" {
  name = "test_cluster"
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "03:58:25Z"
  }
}

data "ionoscloud_k8s_cluster" "test" {
  name = "test_cluster"
}

provider "kubernetes" {
  host = data.ionoscloud_k8s_cluster.test.server
  token =  data.ionoscloud_k8s_cluster.test.user_tokens["cluster-admin"]
}
```

## Example of accessing a kubernetes cluster using the token from the config

```
resource "ionoscloud_k8s_cluster" "test" {
  name = "test_cluster"
  maintenance_window {
    day_of_the_week = "Saturday"
    time            = "03:58:25Z"
  }
}

data "ionoscloud_k8s_cluster" "test" {
  name = "test_cluster"
}

provider "kubernetes" {
  host = data.ionoscloud_k8s_cluster.test.config[0].clusters[0].cluster.server
  token =  data.ionoscloud_k8s_cluster.test.config[0].users[0].user.token
}
```


## Example of dumping the kube_config raw data into a yaml file

**NOTE**: Dumping `kube_config` data into files poses a security risk.

```
data "ionoscloud_k8s_cluster" "k8s_cluster_example" {
  name     = "k8s-demo"
}

resource "null_resource" "getcfg" {
  provisioner "local-exec" {
    command = "echo \"${yamlencode(data.ionoscloud_k8s_cluster.k8s_cluster_example.kube_config)}\" > kubecfg.yaml"
  }
}
```
