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
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned. 
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_image" "example" {
  name        = "ubuntu"
  type        = "CDROM"
  version     = "18.04.3-live-server-amd64.iso"
  location    = "de/fkb"
  cloud_init  = "NONE"
}
```

## Argument Reference

 * `name` - (Required) Name of an existing image that you want to search for.
 * `version` - (Optional) Version of the image (see details below).
 * `location` - (Optional) Id of the existing image's location.
 * `type` - (Optional) The image type, HDD or CD-ROM.
 * `cloud_init` - (Optional) Cloud init compatibility ("NONE" or "V1")

If both "name" and "version" are provided the plugin will concatenate the two strings in this format [name]-[version].

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
 * `location` - Location of that image/snapshot.
