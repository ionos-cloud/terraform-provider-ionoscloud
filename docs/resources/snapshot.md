---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: snapshot"
sidebar_current: "docs-resource-snapshot"
description: |-
  Creates and manages snapshot objects.
---

# ionoscloud\_snapshot

Manages **Snapshots** on IonosCloud.

## Example Usage

```hcl
data "ionoscloud_image" "example" {
    type                  = "HDD"
    image_alias           = "ubuntu:latest"
    location              = "us/las"
}

resource "ionoscloud_datacenter" "example" {
    name                  = "Datacenter Example"
    location              = "us/las"
    description           = "Datacenter Description"
    sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    public                = true
    name                  = "Lan Example"
}

resource "ionoscloud_ipblock" "example" {
    location              = ionoscloud_datacenter.example.location
    size                  = 4
    name                  = "IP Block Example"
}

resource "ionoscloud_server" "example" {
    name                  = "Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    availability_zone     = "ZONE_1"
    cpu_family            = "AMD_OPTERON"
    image_name            = data.ionoscloud_image.example.id
    image_password        = random_password.server_image_password.result
    type                  = "ENTERPRISE"
    volume {
        name              = "system"
        size              = 5
        disk_type         = "SSD Standard"
        user_data         = "foo"
        bus               = "VIRTIO"
        availability_zone = "ZONE_1"
    }
    nic {
        lan               = ionoscloud_lan.example.id
        name              = "system"
        dhcp              = true
        firewall_active   = true
        firewall_type     = "BIDIRECTIONAL"
        ips               = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
    firewall {
        protocol          = "TCP"
        name              = "SSH"
        port_range_start  = 22
        port_range_end    = 22
        source_mac        = "00:0a:95:9d:68:17"
        source_ip         = ionoscloud_ipblock.example.ips[2]
        target_ip         = ionoscloud_ipblock.example.ips[3]
        type              = "EGRESS"
    }
  }
}

resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = ionoscloud_datacenter.example.id
  volume_id     = ionoscloud_server.example.boot_volume
  name          = "Snapshot Example"
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of the Virtual Data Center.
* `name` - (Required)[string] The name of the snapshot.
* `volume_id` - (Required)[string] The ID of the specific volume to take the snapshot from.

## Attribute reference

Beside the configurable arguments, the resource returns the following additional attributes:

* `description` - Human readable description
* `licence_type` - OS type of this Snapshot
* `location` - Location of that image/snapshot
* `size` - The size of the image in GB
* `sec_auth_protection` - Boolean value representing if the snapshot requires extra protection e.g. two factor protection
* `cpu_hot_plug` -  Is capable of CPU hot plug (no reboot required)
* `cpu_hot_unplug` -  Is capable of CPU hot unplug (no reboot required)
* `ram_hot_plug` -  Is capable of memory hot plug (no reboot required)
* `ram_hot_unplug` -  Is capable of memory hot unplug (no reboot required)
* `nic_hot_plug` -  Is capable of nic hot plug (no reboot required)
* `nic_hot_unplug` -  Is capable of nic hot unplug (no reboot required)
* `disc_virtio_hot_plug` -  Is capable of Virt-IO drive hot plug (no reboot required)
* `disc_virtio_hot_unplug` -  Is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.
* `disc_scsi_hot_plug` -  Is capable of SCSI drive hot plug (no reboot required)
* `disc_scsi_hot_unplug` -  Is capable of SCSI drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.


## Import

Resource Snapshot can be imported using the `snapshot id`, e.g.

```shell
terraform import ionoscloud_snapshot.mysnapshot {snapshot uuid}
```