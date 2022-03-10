---
subcategory: "Managed Kubernetes"
layout: "ionoscloud"
page_title: "IonosCloud: k8s_node_pool"
sidebar_current: "docs-resource-k8s-node-pool"
description: |-
  Creates and manages IonosCloud Kubernetes Node Pools.
---

# ionoscloud_k8s_node_pool

Manages a **Managed Kubernetes Node Pool**, part of a managed Kubernetes cluster on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                  = "Datacenter Example"
  location              = "us/las"
  description           = "datacenter description"
  sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = false
  name                  = "Lan Example"
}

resource "ionoscloud_ipblock" "example" {
  location              = "us/las"
  size                  = 3
  name                  = "IP Block Example"
}

resource "ionoscloud_k8s_cluster" "example" {
  name                  = "k8sClusterExample"
  k8s_version           = "1.20.10"
  public                = true
  maintenance_window {
    day_of_the_week     = "Sunday"
    time                = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets { 
     name               = <your_s3_bucket>
  }
}

resource "ionoscloud_k8s_node_pool" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  k8s_cluster_id        = ionoscloud_k8s_cluster.example.id
  name                  = "k8sNodePoolExample"
  k8s_version           = ionoscloud_k8s_cluster.example.k8s_version
  maintenance_window {
    day_of_the_week     = "Monday"
    time                = "09:00:00Z"
  } 
  auto_scaling {
    min_node_count      = 1
    max_node_count      = 1
  }
  cpu_family            = "INTEL_XEON"
  availability_zone     = "AUTO"
  storage_type          = "SSD"
  node_count            = 1
  cores_count           = 2
  ram_size              = 2048
  storage_size          = 40
  public_ips            = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
  lans {
    id                  = ionoscloud_lan.example.id
    dhcp                = true
	routes {
       network          = "1.2.3.5/24"
       gateway_ip       = "10.1.5.17"
     }
   }  
  labels                = {
    lab1                = "value1"
    lab2                = "value2"
  }
  annotations           = {
    ann1                = "value1"
    ann2                = "value2"
  }
}

```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Kubernetes Cluster. *This attribute is immutable*.
- `k8s_version` - (Optional)[string] The desired Kubernetes Version. for supported values, please check the API documentation. The provider will ignore changes of patch level.
- `auto_scaling` - (Optional)[string] Wether the Node Pool should autoscale. For more details, please check the API documentation
  - `min_node_count` - (Optional)[int] The minimum number of worker nodes the node pool can scale down to. Should be less than max_node_count
  - `max_node_count` - (Optional)[int] The maximum number of worker nodes that the node pool can scale to. Should be greater than min_node_count
- `lans` - (Optional)[list] A list of numeric LAN id's you want this node pool to be part of. For more details, please check the API documentation, as well as the example above
  - `id` - (Required)[int] The LAN ID of an existing LAN at the related datacenter
  - `dhcp` - (Optional)[boolean] Indicates if the Kubernetes Node Pool LAN will reserve an IP using DHCP. Default value is `true`
  - `routes` - (Optional) An array of additional LANs attached to worker nodes
    - `network` - (Required)[string] IPv4 or IPv6 CIDR to be routed via the interface
    - `gateway_ip` - (Required)[string] IPv4 or IPv6 Gateway IP for the route
- `maintenance_window` - (Optional) See the **maintenance_window** section in the example above
  - `time` - (Required)[string] A clock time in the day when maintenance is allowed
  - `day_of_the_week` - (Required)[string] Day of the week when maintenance is allowed
- `datacenter_id` - (Required)[string] A Datacenter's UUID
- `k8s_cluster_id`- (Required)[string] A k8s cluster's UUID
- `cpu_family` - (Required)[string] The desired CPU Family - See the API documentation for more information. *This attribute is immutable*.
- `availability_zone` - (Required)[string] - The desired Compute availability zone - See the API documentation for more information. *This attribute is immutable*.
- `storage_type` -(Required)[string] - The desired storage type - SSD/HDD. *This attribute is immutable*.
- `node_count` -(Required)[int] - The desired number of nodes in the node pool
- `cores_count` -(Required)[int] - The CPU cores count for each node of the node pool. *This attribute is immutable*.
- `ram_size` -(Required)[int] - The desired amount of RAM, in MB. *This attribute is immutable*.
- `storage_size` -(Required)[int] - The desired amount of storage for each node, in GB. *This attribute is immutable*.
- `public_ips` - (Optional)[list] A list of public IPs associated with the node pool; must have at least `node_count + 1` elements  
- `labels` - (Optional)[map] A key/value map of labels
- `annotations` - (Optional)[map] A key/value map of annotations
- `gateway_ip` - (Optional)[string] Public IP address for the gateway performing source NAT for the node pool's nodes belonging to a private cluster. Required only if the node pool belongs to a private cluster.

## Import

A Kubernetes Node Pool resource can be imported using its Kubernetes cluster's uuid as well as its own UUID, both of which you can retreive from the cloud API: `resource id`, e.g.:

```shell
terraform import ionoscloud_k8s_node_pool.demo {k8s_cluster_uuid}/{k8s_nodepool_id}
```

This can be helpful when you want to import kubernetes node pools which you have already created manually or using other means, outside of terraform, towards the goal of managing them via Terraform

> :warning: **If you are upgrading from v5.x.x to v6.x.x**: You have to modify you plan for lans to match the new structure, by putting the ids from the old slice in lans.id fields. This is not backwards compatible.
