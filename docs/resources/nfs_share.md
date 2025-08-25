---
subcategory: "Network File Storage"
layout: "ionoscloud"
page_title: "IonosCloud: nfs_share"
sidebar_current: "docs-resource-nfs_share"
description: |-
  Creates and manages Network File Storage (NFS) Share objects on IonosCloud.
---

# ionoscloud_nfs_cluster

Creates and manages Network File Storage (NFS) Share objects on IonosCloud.

## Example Usage

```hcl
# Basic example

resource "ionoscloud_datacenter" "nfs_dc" {
  name                = "NFS Datacenter"
  location            = "de/txl"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "nfs_lan" {
  datacenter_id = ionoscloud_datacenter.nfs_dc.id
  public        = false
  name          = "Lan for NFS"
}

resource "ionoscloud_nfs_cluster" "example" {
  name = "test"
  location = "de/txl"
  size = 2

  nfs {
    min_version = "4.2"
  }

  connections {
    datacenter_id = ionoscloud_datacenter.nfs_dc.id
    ip_address    = "192.168.100.10/24"
    lan           = ionoscloud_lan.nfs_lan.id
  }
}

resource "ionoscloud_nfs_share" "example" {
  location = "de/txl"
  cluster_id = ionoscloud_nfs_cluster.example.id

  name = "example-share"
  quota = 512
  gid = 512
  uid = 512

  client_groups {
    description = "Client Group 1"
    ip_networks = ["10.234.50.0/24"]
    hosts = ["10.234.62.123"]
    nfs {
      squash = "all-anonymous"
    }
  }
}
```

## Argument Reference

## Argument Reference

The following arguments are supported:

- `location` - (Optional) The location of the Network File Storage Cluster. If this is not set and if no value is provided for the `IONOS_API_URL` env var, the default `location` will be: `de/fra`. Other available locations are: `de/fra/2`, `de/txl`, `fr/par`, `gb/lhr`, `es/vit`, `us/las`, `us/ewr`, `us/mci`.
- `cluster_id` - (Required) The ID of the Network File Storage Cluster.
- `name` - (Required) The directory being exported.
- `quota` - (Optional) The quota in MiB for the export. The quota can restrict the amount of data that can be stored within the export. The quota can be disabled using `0`. Default is `0`.
- `gid` - (Optional) The group ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.
- `uid` - (Optional) The user ID that will own the exported directory. If not set, **anonymous** (`512`) will be used.
- `client_groups` - (Required) The groups of clients are the systems connecting to the Network File Storage cluster. Each group includes:
  - `description` - (Optional) Optional description for the clients groups.
  - `ip_networks` - (Required) The allowed host or network to which the export is being shared. The IP address can be either IPv4 or IPv6 and has to be given with CIDR notation.
  - `hosts` - (Required) A singular host allowed to connect to the share. The host can be specified as IP address and can be either IPv4 or IPv6.
  - `nfs` - (Required) NFS specific configurations. Each configuration includes:
    - `squash` - (Required) The squash mode for the export. The squash mode can be:
      - `none` - No squash mode. no mapping,
      - `root-anonymous` - Map root user to anonymous uid,
      - `all-anonymous` - Map all users to anonymous uid. 

## Import

A Network File Storage Share resource can be imported using its `location`, `cluster_id` and `resource id`:

```shell
terraform import ionoscloud_nfs_share.name location:cluster_id:resource_id
```
