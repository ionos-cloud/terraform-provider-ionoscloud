---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_volume"
sidebar_current: "docs-ionoscloud-datasource-volume"
description: |-
  Get information on a Ionos Cloud Volume
---

# ionoscloud\_volume

The volume data source can be used to search for and return existing volumes.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_volume" "example" {
  datacenter_id = "datacenter_id"
  id			= "volume_id"
}
```

### By Name
```hcl
data "ionoscloud_volume" "example" {
  datacenter_id = "datacenter_id"
  name			= "Volume Example"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing volume that you want to search for.
* `id` - (Optional) ID of the volume you want to search for.

Either `volume` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the volume.
* `name` - The name of the volume.
* `disk_type` - The volume type: HDD or SSD.
* `bus` - The bus type of the volume: VIRTIO or IDE.
* `size` -  The size of the volume in GB.
* `ssh_key_path` -  List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Also accepts ssh keys directly.
* `sshkey` - The associated public SSH key.
* `image_password` - Required if `sshkey_path` is not provided.
* `image` - The image or snapshot UUID.
* `licence_type` - The type of the licence.
* `availability_zone` - The storage availability zone assigned to the volume: AUTO, ZONE_1, ZONE_2, or ZONE_3. This property is immutable.
* `user_data` - The cloud-init configuration for the volume as base64 encoded string. The property is immutable and is only allowed to be set on a new volume creation. This option will work only with cloud-init compatible images.
* `backup_unit_id`- The uuid of the Backup Unit that user has access to. The property is immutable and is only allowed to be set on a new volume creation. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property.
* `device_number` - The LUN ID of the storage volume. Null for volumes not mounted to any VM
* `cpu_hot_plug` - Is capable of CPU hot plug (no reboot required)
* `ram_hot_plug` - Is capable of memory hot plug (no reboot required)
* `nic_hot_plug` - Is capable of nic hot plug (no reboot required)
* `nic_hot_unplug` - Is capable of nic hot unplug (no reboot required)
* `disc_virtio_hot_plug` - Is capable of Virt-IO drive hot plug (no reboot required)
* `disc_virtio_hot_unplug` - Is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.
* `boot_server` - The UUID of the attached server.