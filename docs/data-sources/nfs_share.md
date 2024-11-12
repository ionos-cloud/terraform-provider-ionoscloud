---
subcategory: "Network File Storage"
layout: "ionoscloud"
page_title: "IonosCloud: nfs_share"
sidebar_current: "docs-datasource-nfs_share"
description: |-
  Get information on Network File Storage (NFS) Share objects
---

# ionoscloud_nfs_share

Returns information about shares of Network File Storage (NFS) on IonosCloud.

## Example Usage

## By ID
```hcl
data "ionoscloud_nfs_share" "example" {
  location = <location>
  cluster_id = <cluster-id>
  id = <share-id>
}
```

## By Name
```hcl
data "ionoscloud_nfs_share" "example" {
  location = <location>
  cluster_id = <cluster-id>
  name = <partial-name>
  partial_match = true
}

output "share_test_output" {
    value = format("share %s quota %sMiB path '%s'",
        data.ionoscloud_nfs_share.example.name,
        data.ionoscloud_nfs_share.example.quota,
        data.ionoscloud_nfs_share.example.nfs_path,
    )
}
```

## Argument Reference

- `location` - (Optional) The location where the Network File Storage share is located.
- `cluster_id` - (Required) The ID of the Network File Storage cluster.
- `name` - (Optional) Name of the Network File Storage share.
- `id` - (Optional) ID of the Network File Storage share.
- `partial_match` - (Optional) Whether partial matching is allowed or not when using the name filter. Defaults to `false`.

## Attributes Reference

- `id` - The ID of the Network File Storage share.
- `name` - The name of the Network File Storage share.
- `location` - The location where the Network File Storage share is located.
- `cluster_id` - The ID of the Network File Storage cluster.
- `nfs_path` - Path to the NFS export. The NFS path is the path to the directory being exported.
- `quota` - The quota in MiB for the export. The quota can restrict the amount of data that can be stored within the export. The quota can be disabled using `0`.
- `gid` - The group ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.
- `uid` - The user ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.
- `client_groups` - The groups of clients are the systems connecting to the Network File Storage cluster. Each client group supports the following:
    - `description` - Optional description for the clients groups.
    - `ip_networks` - The allowed host or network to which the export is being shared. The IP address can be either IPv4 or IPv6 and has to be given with CIDR notation.
    - `hosts` - A singular host allowed to connect to the share. The host can be specified as IP address and can be either IPv4 or IPv6.
    - `nfs` - The NFS configuration for the client group. Each NFS configuration supports the following:
        - `squash` - The squash mode for the export. The squash mode can be: none - No squash mode. no mapping, root-anonymous - Map root user to anonymous uid, all-anonymous - Map all users to anonymous uid.