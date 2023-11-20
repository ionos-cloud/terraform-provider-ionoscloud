---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-volume"
description: |-
  Creates and manages IonosCloud Volume objects.
---

# ionoscloud\_volume

Manages a **Volume** on IonosCloud.

## Example Usage

A primary volume will be created with the server. If there is a need for additional volumes, this resource handles it. Any of the additional volumes can be used as a boot volume.

```hcl
data "ionoscloud_image" "example" {
    type                  = "HDD"
    cloud_init            = "V1"
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
        is_boot_volume    = false
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

resource "ionoscloud_volume" "example" {
  datacenter_id           = ionoscloud_datacenter.example.id
  server_id               = ionoscloud_server.example.id
  name                    = "Volume Example"
  availability_zone       = "ZONE_1"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  image_name              = data.ionoscloud_image.example.id
  image_password          = random_password.volume_image_password.result
  user_data               = "foo"
  is_boot_volume          = true
}

resource "ionoscloud_volume" "example" {
  datacenter_id           = ionoscloud_datacenter.example.id
  server_id               = ionoscloud_server.example.id
  name                    = "Another Volume Example"
  availability_zone       = "ZONE_1"
  size                    = 5
  disk_type               = "SSD Standard"
  bus                     = "VIRTIO"
  licence_type            = "OTHER"
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}

resource "random_password" "volume_image_password" {
  length           = 16
  special          = false
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `server_id` - (Required)[string] The ID of a server.
* `disk_type` - (Required)[string] The volume type: HDD or SSD. This property is immutable.
* `bus` - (Optional)[Boolean] The bus type of the volume: VIRTIO or IDE.
* `size` -  (Required)[integer] The size of the volume in GB.
* `ssh_key_path` -  (Optional)[list] List of absolute paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Also accepts ssh keys directly. Required for IonosCloud Linux images. Required if `image_password` is not provided. This property is immutable.
* `ssh_keys` -  (Optional)[list] List of absolute paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Also accepts ssh keys directly. Required for IonosCloud Linux images. Required if `image_password` is not provided. This property is immutable.
* `sshkey` - (Computed) The associated public SSH key.
* `image_password` - (Optional)[string] Required if `sshkey_path` is not provided.
* `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
* `image` - (Computed) The image or snapshot UUID.
* `image_alias` - (Computed) The image alias.
* `licence_type` - (Optional)[string] Required if `image_name` is not provided.
* `name` - (Optional)[string] The name of the volume.
* `availability_zone` - (Optional)[string] The storage availability zone assigned to the volume: AUTO, ZONE_1, ZONE_2, or ZONE_3. This property is immutable
* `user_data` - (Optional)[string] The cloud-init configuration for the volume as base64 encoded string. The property is immutable and is only allowed to be set on a new volume creation. This option will work only with cloud-init compatible images.
* `backup_unit_id`- (Optional)[string] The uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property.
* `device_number`- (Computed) The Logical Unit Number of the storage volume. Null for volumes not mounted to any VM.
* `pci_slot`- (Computed) The PCI slot number of the storage volume. Null for volumes not mounted to any VM.
* `cpu_hot_plug` - (Computed)[string] Is capable of CPU hot plug (no reboot required)
* `ram_hot_plug` - (Computed)[string] Is capable of memory hot plug (no reboot required)
* `nic_hot_plug` - (Computed)[string] Is capable of nic hot plug (no reboot required)
* `nic_hot_unplug` - (Computed)[string] Is capable of nic hot unplug (no reboot required)
* `disc_virtio_hot_plug` - (Computed)[string] Is capable of Virt-IO drive hot plug (no reboot required)
* `disc_virtio_hot_unplug` - (Computed)[string] Is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.
* `boot_server` - (Computed)[string] The UUID of the attached server.
* `is_boot_volume` - (Computed)(Optional)[boolean] The volume can be set as the primary boot device of the server to which it is attached. If the property is omitted, the inline volume will be set as primary boot device, by default. Setting this property while a different volume is already the primary boot device will result in the other volume being unset, and the current volume becoming the primary boot device. There will always be one boot volume for the server.
> **⚠ WARNING**
>
> ssh_key_path and ssh_keys fields are immutable.
> If you want to create a **CUBE** server, the type of the inline volume must be set to **DAS**. In this case, you can not set the `size` argument since it is taken from the `template_uuid` you set in the server.


## Import

Resource Volume can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_volume.myvolume {datacenter uuid}/{server uuid}/{volume uuid}
```
