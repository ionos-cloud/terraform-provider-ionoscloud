---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_gpu_server"
sidebar_current: "docs-ionoscloud-datasource-gpu_server"
description: |-
  Get information on a Ionos Cloud GPU Server
---

# ionoscloud_gpu_server

The [GPU Server data source](https://docs.ionos.com/cloud/compute-services/compute-engine/cloud-gpu-vm) can be used to search for and return existing GPU servers. 
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_gpu_server" "example" {
   datacenter_id = "datacenter_id"
   id			 = "server_id"
}
```

### By Name
```hcl
data "ionoscloud_gpu_server" "example" {
   datacenter_id = "datacenter_id"
   name			 = "GPU Server Example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing server that you want to search for.
* `id` - (Optional) ID of the server you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `template_uuid` - The UUID of the template for creating a GPU server; the available templates for GPU servers can be found on the templates resource
* `id` - The id of that resource
* `name` - The name of that resource
* `hostname` - The hostname of the server
* `type` - Server usages: GPU
* `vm_state`- Status of the virtual Machine
* `datacenter_id` - The id of the datacenter
* `availability_zone` - The availability zone in which the server should exist
* `vm_state` - Status of the virtual Machine
* `boot_cdrom`
* `ram`
* `cores`
* `boot_volume`
* `boot_image`
* `token`
* `cdroms` - list of
  * `id` - Id of the attached cdrom
  * `name` - The name of the attached cdrom
  * `description` - Description of cdrom
  * `location` - Location of that image/snapshot
  * `size` - The size of the image in GB
  * `cpu_hot_plug` - Is capable of CPU hot plug (no reboot required)
  * `cpu_hot_unplug` - Is capable of CPU hot unplug (no reboot required)
  * `ram_hot_plug` - Is capable of memory hot plug (no reboot required)
  * `ram_hot_unplug` - Is capable of memory hot unplug (no reboot required)
  * `nic_hot_plug` - Is capable of nic hot plug (no reboot required)
  * `nic_hot_unplug` - Is capable of nic hot unplug (no reboot required)
  * `disc_virtio_hot_plug` - Is capable of Virt-IO drive hot plug (no reboot required)
  * `disc_virtio_hot_unplug` - Is capable of Virt-IO drive hot unplug (no reboot required)
  * `disc_scsi_hot_plug` - Is capable of SCSI drive hot plug (no reboot required)
  * `disc_scsi_hot_unplug` - Is capable of SCSI drive hot unplug (no reboot required)
  * `licence_type` - OS type of this Image
  * `image_type` - Type of image
  * `image_aliases` - List of image aliases mapped for this Image
  * `public` - Indicates if the image is part of the public repository or not
  * `image_aliases` - List of image aliases mapped for this Image
  * `cloud_init` - Cloud init compatibility
* `volumes` - list of
  * `id` - Id of the attached volume
  * `name` - Name of the attached volume
  * `type` - Hardware type of the volume.
  * `availability_zone` - The availability zone in which the volume should exist
  * `image` - Image or snapshot ID to be used as template for this volume
  * `image_password` - Initial password to be set for installed OS
  * `ssh_keys` - Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key
  * `bus` - The bus type of the volume
  * `licence_type` - OS type of this volume
  * `cpu_hot_plug` - Is capable of CPU hot plug (no reboot required)
  * `ram_hot_plug` - Is capable of memory hot plug (no reboot required)
  * `nic_hot_plug` - Is capable of nic hot plug (no reboot required)
  * `nic_hot_unplug` - Is capable of nic hot unplug (no reboot required)
  * `disc_virtio_hot_plug` - Is capable of Virt-IO drive hot plug (no reboot required)
  * `disc_virtio_hot_unplug` - Is capable of Virt-IO drive hot unplug (no reboot required)
  * `device_number` - The Logical Unit Number of the storage volume
  * `pci_slot` - The PCI slot number of the storage volume
