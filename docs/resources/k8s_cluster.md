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
- `public` - The indicator if the cluster is public or private. Be aware that setting it to false is currently in beta phase. Default value is true
- `gateway_ip` - The IP address of the gateway used by the cluster. This is mandatory when `public` is set to `false` and should not be provided otherwise.
  - `time` - (Required)[string] A clock time in the day when maintenance is allowed
  - `day_of_the_week` - (Required)[string] Day of the week when maintenance is allowed
- `available_upgrade_versions` - (Computed) List of available versions for upgrading the cluster
- `viable_node_pool_versions` - (Computed) List of versions that may be used for node pools under this cluster

## Import

A Kubernetes Cluster resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_k8s_cluster.demo {k8s_cluster uuid}
```

This can be helpful when you want to import kubernetes clusters which you have already created manually or using other means, outside of terraform.

## Important Notes

- Please note that every `ionoscloud_datacenter` resource you plan to add kubernetes node pools for the cluster to needs to also be specified as a dependency of the Kubernetes cluster by using the `depends_on` meta-property (For more details, please see https://www.terraform.io/docs/configuration/resources.html#resource-dependencies). This will ensure that resources are destroyed in the right order. In case you do not do this, you might encounter problems when deleting the Virtual Datacenter. In return, this will give you the ability to keep the data in your PersistentVolumeClaims across NodePools created in the same Virtual Datacenter `ionoscloud_datacenter`