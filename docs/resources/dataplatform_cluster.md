---
subcategory: "Dataplatform"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_dataplatform_cluster"
sidebar_current: "docs-resource-dataplatform_cluster"
description: |-
  Creates and manages Dataplatform Cluster objects.
---

# ionoscloud_dataplatform_cluster

Manages a **Dataplatform Cluster**.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name        = "Datacenter_Example"
  location    = "de/txl"
  description = "Datacenter for testing Dataplatform Cluster"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "LAN_Example"
}

resource "ionoscloud_dataplatform_cluster" "example" {
  datacenter_id   		=  ionoscloud_datacenter.example.id
  name 					= "Dataplatform_Cluster_Example"
  maintenance_window {
    day_of_the_week  	= "Sunday"
    time				= "09:00:00"
  }
  version	= "23.11"
  lans {
    lan_id = ionoscloud_lan.example.id
    dhcp = false
    routes {
      network = "182.168.42.1/24"
      gateway = "192.168.42.1"
    }
  }
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The UUID of the virtual data center (VDC) the cluster is provisioned.
* `name` - (Required)[string] The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]). It can contain dashes (-), underscores (_), dots (.), and alphanumerics in-between.
* `version` - (Optional)[int] The version of the Data Platform.
* `maintenance_window` - (Optional) Starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  * `time` - (Required)[string] Time at which the maintenance should start. Must conform to the 'HH:MM:SS' 24-hour format. This pattern matches the "HH:MM:SS 24-hour format with leading 0" format. For more information take a look at [this link](https://stackoverflow.com/questions/7536755/regular-expression-for-matching-hhmm-time-format).
  * `day_of_the_week` - (Required)[string] Must be set with one the values `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` or `Sunday`.
* `lans` - (Optional)[list] A list of LANs you want this node pool to be part of.
  * `lan_id` - (Required)[string] The LAN ID of an existing LAN at the related data center.
  * `dhcp` - (Optional)[bool] Indicates if the Kubernetes node pool LAN will reserve an IP using DHCP. The default value is 'true'.
  * `routes` - (Optional)[list] An array of additional LANs attached to worker nodes.
    * `gateway` - (Required)[string] IPv4 or IPv6 gateway IP for the route.
    * `network` - (Required)[string] IPv4 or IPv6 CIDR to be routed via the interface.

## Import

Resource Dataplatform Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import ionoscloud_dataplatform_cluster.mycluser {cluster uuid}
```
