---
subcategory: "Network File Storage"
layout: "ionoscloud"
page_title: "IonosCloud: nfs_cluster"
sidebar_current: "docs-resource-nfs_cluster"
description: |-
  Creates and manages Network File Storage (NFS) Cluster objects
---

# ionoscloud_nfs_cluster

Create clusters of Network File Storage (NFS) on IonosCloud.

## Example Usage

```hcl
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

data "ionoscloud_image" "HDD_image" {
  image_alias = "ubuntu:20.04"
  type        = "HDD"
  cloud_init  = "V1"
  location    = "de/txl"
}

resource "random_password" "password" {
  length  = 16
  special = false
}

# needed for the NIC - which provides the IP address for the NFS cluster.
resource "ionoscloud_server" "nfs_server" {
  name              = "Server for NFS"
  datacenter_id     = ionoscloud_datacenter.nfs_dc.id
  cores             = 1
  ram               = 2048
  availability_zone = "ZONE_1"
  cpu_family        = "INTEL_SKYLAKE"
  image_name        = data.ionoscloud_image.HDD_image.id
  image_password    = random_password.password.result
  volume {
    name      = "system"
    size      = 14
    disk_type = "SSD"
  }
  nic {
    name            = "NIC A"
    lan             = ionoscloud_lan.nfs_lan.id
    dhcp            = true
    firewall_active = true
  }
}

locals {
  nfs_server_cidr = format("%s/24", ionoscloud_server.nfs_server.nic[0].ips[0])
  nfs_cluster_ip = cidrhost(local.nfs_server_cidr, 10)
  nfs_cluster_cidr = format("%s/24", local.nfs_cluster_ip)
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
    ip_address    = local.nfs_cluster_cidr
    lan           = ionoscloud_lan.nfs_lan.id
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the Network File Storage cluster.
- `location` - (Required) The location where the Network File Storage cluster is located.
  - `de/fra` - Frankfurt
  - `de/txl` - Berlin
- `size` - (Required) The size of the Network File Storage cluster in TiB. Note that the cluster size cannot be reduced after provisioning. This value determines the billing fees. Default is `2`. The minimum value is `2` and the maximum value is `42`.
- `nfs` - (Optional) The NFS configuration for the Network File Storage cluster. Each NFS configuration supports the following:
    - `min_version` - (Optional) The minimum supported version of the NFS cluster. Supported values: `4.2`. Default is `4.2`.
- `connections` - (Required) A list of connections for the Network File Storage cluster. You can specify only one connection. Connections are **immutable**. Each connection supports the following:
    - `datacenter_id` - (Required) The ID of the datacenter where the Network File Storage cluster is located.
    - `ip_address` - (Required) The IP address and prefix of the Network File Storage cluster. The IP address can be either IPv4 or IPv6. The IP address has to be given with CIDR notation. 
    - `lan` - (Required) The Private LAN to which the Network File Storage cluster must be connected.

## Import

A Network File Storage Cluster resource can be imported using its `location` and `resource id`:

```shell
terraform import ionoscloud_nfs_cluster.name {location}:{uuid}
```
