---
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-volume"
description: |-
  Creates and manages IonosCloud Volume objects.
---

# ionoscloud\_volume

Manages a volume on IonosCloud.

## Example Usage

A primary volume will be created with the server. If there is a need for additional volumes, this resource handles it.

```hcl
resource "ionoscloud_volume" "example" {
  datacenter_id = "${ionoscloud_datacenter.example.id}"
  server_id     = "${ionoscloud_server.example.id}"
  image_name    = "${var.ubuntu}"
  size          = 5
  disk_type     = "HDD"
  ssh_key_path  = "${var.private_key_path}"
  bus           = "VIRTIO"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `server_id` - (Required)[string] The ID of a server.
* `disk_type` - (Required)[string] The volume type: HDD or SSD.
* `bus` - (Optional)[Boolean] The bus type of the volume: VIRTIO or IDE.
* `size` -  (Required)[integer] The size of the volume in GB.
* `ssh_key_path` -  (Optional)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
* `sshkey` - (Computed) The associated public SSH key.
* `image_password` - (Optional)[string] Required if `sshkey_path` is not provided.
* `image_name` - (Optional)[string] The image or snapshot UUID. May also be an image alias. It is required if `licence_type` is not provided.
* `licence_type` - (Optional)[string] Required if `image_name` is not provided.
* `name` - (Optional)[string] The name of the volume.
* `availability_zone` - (Optional)[string] The storage availability zone assigned to the volume: AUTO, ZONE_1, ZONE_2, or ZONE_3.
* `user_data` - (Optional)[string] The cloud-init configuration for the volume as base64 encoded string. The property is immutable and is only allowed to be set on a new volume creation. This option will work only with cloud-init compatible images.
* `backup_unit_id`- (Optional)[string] he uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property.
* `device_number` - (Computed) The LUN ID of the storage volume. Null for volumes not mounted to any VM
* `cpu_hot_plug` - (Computed)[string] Is capable of CPU hot plug (no reboot required)
* `ram_hot_plug` - (Computed)[string] Is capable of memory hot plug (no reboot required)
* `nic_hot_plug` - (Computed)[string] Is capable of nic hot plug (no reboot required)
* `nic_hot_unplug` - (Computed)[string] Is capable of nic hot unplug (no reboot required)
* `disc_virtio_hot_plug` - (Computed)[string] Is capable of Virt-IO drive hot plug (no reboot required)
* `disc_virtio_hot_unplug` - (Computed)[string] Is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.

## Import

Resource Volume can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_volume.myvolume {datacenter uuid}/{server uuid}/{volume uuid}
```
