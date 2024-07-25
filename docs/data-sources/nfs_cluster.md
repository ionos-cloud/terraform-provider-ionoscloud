---
subcategory: "Network File Storage"
layout: "ionoscloud"
page_title: "IonosCloud: nfs_cluster"
sidebar_current: "docs-datasource-nfs_cluster"
description: |-
  Get information on Network File Storage (NFS) Cluster objects
---

# ionoscloud_nfs_cluster

Returns information about clusters of Network File Storage (NFS) on IonosCloud.

## Example Usage

## By ID
```hcl
data "ionoscloud_nfs_cluster" "example" {
  location = <location>
  id = <cluster-id>
}
```

## By Name
```hcl
data "ionoscloud_nfs_cluster" "example" {
  location = <location>
  name = <partial-name>
  partial_match = true
}
```

## Argument Reference

* `location` - (Required) The location where the Network File Storage cluster is located.
* `name` - (Optional) Name of the Network File Storage cluster.
* `id` - (Optional) ID of the Network File Storage cluster.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using the name filter. Defaults to `false`.

## Attributes Reference

The following attributes are returned by the datasource:

-`id` - The ID of the Network File Storage cluster.
- `name` - The name of the Network File Storage cluster.
- `location` - The location where the Network File Storage cluster is located.
- `size` - The size of the Network File Storage cluster in TiB. Note that the cluster size cannot be reduced after provisioning. This value determines the billing fees. Default is `2`. The minimum value is `2` and the maximum value is `42`.
- `nfs` - The NFS configuration for the Network File Storage cluster. Each NFS configuration supports the following:
    - `min_version` - The minimum supported version of the NFS cluster. Default is `4.2`
- `connections` - A list of connections for the Network File Storage cluster. You can specify only one connection. Each connection supports the following:
    - `datacenter_id` - The ID of the datacenter where the Network File Storage cluster is located.
    - `ip_address` - The IP address and prefix of the Network File Storage cluster. The IP address can be either IPv4 or IPv6. The IP address has to be given with CIDR notation.
    - `lan` - The Private LAN to which the Network File Storage cluster must be connected.
