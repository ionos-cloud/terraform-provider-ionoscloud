---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dataplatform_cluster"
sidebar_current: "docs-dataplatform_cluster"
description: |-
  Get information on a Dataplatform Cluster.
---

# ionoscloud_dataplatform_cluster

⚠️ **Note:** Data Platform is currently in the Early Access (EA) phase.
We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

The **Dataplatform Cluster Data Source** can be used to search for and return an existing Dataplatform Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_dataplatform_cluster" "example" {
  id	= <cluster_id>
}
```

### By Name

```hcl
data "ionoscloud_dataplatform_cluster" "example" {
  name	= "Dataplatform_Cluster_Example"
}
```

### By Name with Partial Match

```hcl
data "ionoscloud_dataplatform_cluster" "example" {
  name	= "_Example"
  partial_match = true
}
```

## Argument Reference

* `id` - (Optional) ID of the cluster you want to search for.
* `name` - (Optional) Name of an existing cluster that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the cluster.
* `datacenter_id` - The UUID of the virtual data center (VDC) in which the cluster is provisioned.
* `name` - The name of your cluster.
* `data_platform_version` - The version of the Data Platform.
* `maintenance_window` - Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  * `time` - Time at which the maintenance should start. 
  * `day_of_the_week`
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
  * key - is the user name
  * value - is the token
* `server` - cluster server (same as `config[0].clusters[0].cluster.server` but provided as an attribute for ease of use)
* `ca_crt` - base64 decoded cluster certificate authority data (provided as an attribute for direct use)

**NOTE**: The whole `config` node is marked as **sensitive**.

## Example of accessing a Dataplatform Cluster using the user's token

```
resource "ionoscloud_dataplatform_cluster" "example" {
  datacenter_id   		=  ionoscloud_datacenter.example.id
  name 					= "Dataplatform_Cluster_Example"
  maintenance_window {
    day_of_the_week  	= "Sunday"
    time				= "09:00:00"
  }
  data_platform_version	= "22.11"
}

data "ionoscloud_dataplatform_cluster" "example" {
  name = "Dataplatform_Cluster_Example"
}

provider "kubernetes" {
  host = data.ionoscloud_dataplatform_cluster.example.server
  token =  data.ionoscloud_dataplatform_cluster.example.user_tokens["cluster-admin"]
}
```

## Example of accessing a kubernetes cluster using the token from the config

```
resource "ionoscloud_dataplatform_cluster" "example" {
  datacenter_id   		=  ionoscloud_datacenter.example.id
  name 					= "Dataplatform_Cluster_Example"
  maintenance_window {
    day_of_the_week  	= "Sunday"
    time				= "09:00:00"
  }
  data_platform_version	= "22.11"
}

data "ionoscloud_dataplatform_cluster" "example" {
  name = "Dataplatform_Cluster_Example"
}

provider "kubernetes" {
  host = data.ionoscloud_dataplatform_cluster.example.config[0].clusters[0].cluster.server
  token =  data.ionoscloud_dataplatform_cluster.example.config[0].users[0].user.token
}
```


## Example of dumping the kube_config raw data into a yaml file

**NOTE**: Dumping `kube_config` data into files poses a security risk.

**NOTE**: Using `sensitive_content` for `local_file` does not show the data written to the file during the plan phase.

```
data "ionoscloud_dataplatform_cluster" "example" {
  name = "Dataplatform_Cluster_Example"
}

resource "local_file" "kubeconfig" {
    sensitive_content     = yamlencode(jsondecode(data.ionoscloud_dataplatform_cluster.example.kube_config))
    filename              = "kubeconfig.yaml"
}

```