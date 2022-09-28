---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: image"
sidebar_current: "docs-datasource-image"
description: |-
  Get information on a IonosCloud Images
---

# ionoscloud\_image

The **Image data source** can be used to search for and return an existing image which can then be used to provision a server.  

## Example Usage

```hcl
data "ionoscloud_image" "example" {
  name                  = "debian-10-genericcloud"
  partial_match         = true
  type                  = "HDD"
  cloud_init            = "V1"
  location              = "us/las"
}
```

## Argument Reference

 * `name` - (Required) Name of an existing image that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
 * `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.
 * `location` - (Optional) Id of the existing image's location.
 * `type` - (Optional) The image type, HDD or CD-ROM.
 * `cloud_init` - (Optional) Cloud init compatibility ("NONE" or "V1")

## Attributes Reference

 * `id` - UUID of the image
 * `name` - name of the image
 * `description` - description of the image
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
 * `license_type` - OS type of this Image
 * `public` - Indicates if the image is part of the public repository or not
 * `image_aliases` - List of image aliases mapped for this Image
 * `cloud_init` - Cloud init compatibility
 * `type` - This indicates the type of image
