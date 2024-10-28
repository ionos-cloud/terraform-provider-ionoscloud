---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dataplatform_node_pool"
sidebar_current: "docs-resource-dataplatform_node_pool"
description: |-
  Creates and manages Dataplatform Node Pool objects.
---

# ionoscloud_dataplatform_node_pool

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name        = "Datacenter_Example"
  location    = "de/txl"
  description = "Datacenter for testing Dataplatform Cluster"
}

resource "ionoscloud_dataplatform_cluster" "example" {
  datacenter_id   		=  ionoscloud_datacenter.example.id
  name 					= "Dataplatform_Cluster_Example"
  maintenance_window {
    day_of_the_week  	= "Sunday"
    time				= "09:00:00"
  }
  version	= "23.7"
}

resource "ionoscloud_dataplatform_node_pool" "example" {
  cluster_id        = ionoscloud_dataplatform_cluster.example.id
  name              = "Dataplatform_Node_Pool_Example"
  node_count        = 1
  cpu_family        = "INTEL_SKYLAKE"
  cores_count       = 1
  ram_size          = 2048
  availability_zone = "AUTO"
  storage_type      = "HDD"
  storage_size      = 10
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:00:00"
  }
  labels 			= {
    foo   			= "bar"
    color 			= "green"
  }
  annotations 		= {
    ann1 			= "value1"
    ann2 			= "value2"
  }
}
```

## Argument reference

* `cluster_id` - (Required)[string] The UUID of an existing Dataplatform cluster.
* `name` - (Required)[string] The name of your node pool. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]). It can contain dashes (-), underscores (_), dots (.), and alphanumerics in-between.
* `node_count` - (Required)[int] The number of nodes that make up the node pool. Must be set with a minimum value of 1.
* `cpu_family` - (Optional)[string] A valid CPU family name or `AUTO` if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource. The default value is `AUTO`.
* `cores_count` - (Optional)[int] The number of CPU cores per node. Must be set with a minimum value of 1. The default value is `4`.
* `ram_size` - (Optional)[int] The RAM size for one node in MB. Must be set in multiples of `1024`MB, with a minimum size is of `2048`MB. The default value is `4096`.
* `availability_zone` - (Optional)[string] The availability zone of the virtual datacenter region where the node pool resources should be provisioned. Must be set with one of the values `AUTO`, `ZONE_1` or `ZONE_2`. The default value is `AUTO`.
* `storage_type` - (Optional)[int] The type of hardware for the volume. Must be set with one of the values `HDD` or `SSD`. The default value is `SSD`.
* `storage_size` - (Optional)[int] The size of the volume in GB. The size must be greater than `10`GB. The default value is `20`.
* `maintenance_window` - (Optional) Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
    * `time` - (Required)[string] Time at which the maintenance should start. Must conform to the 'HH:MM:SS' 24-hour format. This pattern matches the "HH:MM:SS 24-hour format with leading 0" format. For more information take a look at [this link](https://stackoverflow.com/questions/7536755/regular-expression-for-matching-hhmm-time-format).
    * `day_of_the_week` - (Required)[string] Must be set with one the values `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` or `Sunday`.
* `labels` - (Optional)[map] Key-value pairs attached to the node pool resource as [Kubernetes labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
* `annotations` - (Optional)[map] Key-value pairs attached to node pool resource as [Kubernetes annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/).

## Import

A Dataplatform Node Pool resource can be imported using its cluster's UUID as well as its own UUID, e.g.:

```shell
terraform import ionoscloud_dataplatform_node_pool.mynodepool {dataplatform_cluster_uuid}/{dataplatform_nodepool_id}
```
