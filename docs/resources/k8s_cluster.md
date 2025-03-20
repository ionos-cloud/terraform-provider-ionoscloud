---
subcategory: "Managed Kubernetes"
layout: "ionoscloud"
page_title: "IonosCloud: k8s_cluster"
sidebar_current: "docs-resource-k8s-cluster"
description: |-
  Creates and manages IonosCloud Kubernetes Clusters.
---

# ionoscloud_k8s_cluster

Manages a [Managed Kubernetes Cluster](https://docs.ionos.com/cloud/containers/managed-kubernetes/overview) on IonosCloud.

## Example Usage

### Public cluster

```hcl
resource "ionoscloud_k8s_cluster" "example" {
  name                  = "k8sClusterExample"
  k8s_version           = "1.31.2"
  maintenance_window {
    day_of_the_week     = "Sunday"
    time                = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets { 
     name               = "globally_unique_bucket_name"
  }
}
```

### Private Cluster

```hcl
resource "ionoscloud_datacenter" "testdatacenter" {
  name                    = "example"
  location                = "de/fra"
  description             = "Test datacenter"
}

resource "ionoscloud_ipblock" "k8sip" {
  location = "de/fra"
  size = 1
  name = "IP Block Private K8s"
}

resource "ionoscloud_k8s_cluster" "example" {
  name                  = "k8sClusterExample"
  k8s_version           = "1.31.2"
  maintenance_window {
    day_of_the_week     = "Sunday"
    time                = "09:00:00Z"
  }
  api_subnet_allow_list = ["1.2.3.4/32"]
  s3_buckets {
     name               = "globally_unique_bucket_name"
  }
  location = "de/fra"
  nat_gateway_ip = ionoscloud_ipblock.k8sip.ips[0]
  node_subnet = "192.168.0.0/16"
  public = false
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Kubernetes Cluster.
- `k8s_version` - (Optional)[string] The desired Kubernetes Version. For supported values, please check the API documentation. Downgrades are not supported. The provider will ignore downgrades of patch level.
- `maintenance_window` - (Optional) A maintenance window comprise of a day of the week and a time for maintenance to be allowed
  - `time` - (Required)[string] A clock time in the day when maintenance is allowed
  - `day_of_the_week` - (Required)[string] Day of the week when maintenance is allowed
- `viable_node_pool_versions` - (Computed)[list] List of versions that may be used for node pools under this cluster
- `api_subnet_allow_list` - (Optional)[list] Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6.
- `s3_buckets` - (Optional)[list] List of IONOS Object Storage buckets configured for K8s usage. For now it contains only an IONOS Object Storage bucket used to store K8s API audit logs.
- `public` - (Optional)[boolean] Indicates if the cluster is public or private. This attribute is immutable.
- `nat_gateway_ip` - (Optional)[string] The NAT gateway IP of the cluster if the cluster is private. This attribute is immutable. Must be a reserved IP in the same location as the cluster's location. This attribute is mandatory if the cluster is private.
- `node_subnet` - (Optional)[string] The node subnet of the cluster, if the cluster is private. This attribute is optional and immutable. Must be a valid CIDR notation for an IPv4 network prefix of 16 bits length.
- `location` - (Optional)[string] This attribute is mandatory if the cluster is private. The location must be enabled for your contract, or you must have a data center at that location. This property is not adjustable.
- `allow_replace` - (Optional)[bool] When set to true, allows the update of immutable fields by first destroying and then re-creating the cluster.

⚠️ **_Warning: `allow_replace` - lets you update immutable fields, but it first destroys and then re-creates the cluster in order to do it. Set the field to true only if you know what you are doing._**

## Import

A Kubernetes Cluster resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_k8s_cluster.demo k8s_cluster uuid
```

This can be helpful when you want to import kubernetes clusters which you have already created manually or using other means, outside of terraform.

⚠️ **_Warning: **During a maintenance window, k8s can update your `k8s_version` if the old one reaches end of life. This upgrade will not be shown in the plan, as we prevent
terraform from doing a downgrade, as downgrading `k8s_version` is not supported._**