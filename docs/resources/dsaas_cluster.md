---
subcategory: "Data Stack as a Service"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dsaas_cluster"
sidebar_current: "docs-resource-dsaas_cluster"
description: |-
Creates and manages DSaaS Cluster objects.
---

# ionoscloud\_pg_cluster

Manages a **DSaaS Cluster**.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name        = "Datacenter_Example"
  location    = "de/txl"
  description = "Datacenter for testing DSaaS Cluster"
}

resource "ionoscloud_dsaas_cluster" "example" {
  datacenter_id   		=  ionoscloud_datacenter.example.id
  name 					= "DSaaS_Cluster_Example"
  maintenance_window {
    day_of_the_week  	= "Sunday"
    time				= "09:00:00"
  }
  data_platform_version	= "1.1.0"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The UUID of the virtual data center (VDC) the cluster is provisioned.
* `name` - (Required)[string] The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
* `data_platform_version` - (Optional)[int] The version of the DataPlatform.
* `maintenance_window` - (Optional)[string] Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  * `time` - (Required)[string] Time at which the maintenance should start. Must conform to the 'HH:MM:SS' 24-hour format. This pattern matches the "HH:MM:SS 24-hour format with leading 0" format. For more information take a look at [this link](https://stackoverflow.com/questions/7536755/regular-expression-for-matching-hhmm-time-format).
  * `day_of_the_week` - (Required)[string] Must be set with one the values `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` or `Sunday`.

## Import

Resource DSaaS Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import ionoscloud_dsaas_cluster.mycluser {cluster uuid}
```