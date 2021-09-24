---
layout: "ionoscloud"
page_title: "IonosCloud: k8s_cluster"
sidebar_current: "docs-resource-k8s-cluster"
description: |-
  Creates and manages IonosCloud Kubernetes Clusters.
---

# ionoscloud_k8s_cluster

Manages a Managed Kubernetes cluster on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_k8s_cluster" "example" {
  name        = "example"
  k8s_version = "1.18.5"
  maintenance_window {
    day_of_the_week = "Monday"
    time            = "09:30:00Z"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the Kubernetes Cluster.
- `k8s_version` - (Optional)[string] The desired Kubernetes Version. For supported values, please check the API documentation. The provider will ignore changes of patch level.
- `maintenance_window` - (Optional) See the **maintenance_window** section in the example above
  - `time` - (Required)[string] A clock time in the day when maintenance is allowed
  - `day_of_the_week` - (Required)[string] Day of the week when maintenance is allowed
- `available_upgrade_versions` - (Computed) List of available versions for upgrading the cluster
- `viable_node_pool_versions` - (Computed) List of versions that may be used for node pools under this cluster
- `api_subnet_allow_list` - (Optional) Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6.
- `s3_buckets` - (Optional) List of S3 bucket configured for K8s usage. For now it contains only an S3 bucket used to store K8s API audit logs.

## Import

A Kubernetes Cluster resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_k8s_cluster.demo {k8s_cluster uuid}
```

This can be helpful when you want to import kubernetes clusters which you have already created manually or using other means, outside of terraform.

## Important Notes

- Please note that every `ionoscloud_datacenter` resource you plan to add kubernetes node pools for the cluster to needs to also be specified as a dependency of the Kubernetes cluster by using the `depends_on` meta-property (For more details, please see https://www.terraform.io/docs/configuration/resources.html#resource-dependencies). This will ensure that resources are destroyed in the right order. In case you do not do this, you might encounter problems when deleting the Virtual Datacenter. In return, this will give you the ability to keep the data in your PersistentVolumeClaims across NodePools created in the same Virtual Datacenter `ionoscloud_datacenter`