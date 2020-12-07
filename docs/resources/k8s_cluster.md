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
- `k8s_version` - (Optional)[string] The desired Kubernetes Version. For supported values, please check the API documentation.
- `maintenance_window` - (Optional) See the **maintenance_window** section in the example above

## Import

A Kubernetes Cluster resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_k8s_cluster.demo {k8s_cluster uuid}
```

This can be helpful when you want to import kubernetes clusters which you have already created manually or using other means, outside of terraform.

## Important Notes

- Please note that every `ionoscloud_datacenter` resource you plan to add kubernetes node pools for the cluster to needs to also be specified as a dependency of the Kubernetes cluster by using the `depends_on` meta-property (For more details, please see https://www.terraform.io/docs/configuration/resources.html#resource-dependencies). This will ensure that resources are destroyed in the right order. In case you do not do this, you might encounter problems when deleting the Virtual Datacenter. In return, this will give you the ability to keep the data in your PersistentVolumeClaims across NodePools created in the same Virtual Datacenter `ionoscloud_datacenter`