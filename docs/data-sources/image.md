---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: image"
sidebar_current: "docs-datasource-image"
description: |-
  Get information on a IonosCloud Image
---

# ionoscloud\_image

The **Image data source** can be used to search for and return an existing image which can then be used to provision a server.  
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned. 
When this happens, please refine your search string so that it is specific enough to return only one result. In case multiple matches are found, enable debug(`TF_LOG=debug`) to show the name and location of the images.
## Example Usage

```hcl
data "ionoscloud_image" "cdrom" {
  image_alias = "ubuntu:latest_iso"
  type        = "CDROM"
  location    = "de/txl"
  cloud_init  = "NONE"
}
```
Finds an image with alias `ubuntu:latest_iso`, in location `de/txl`, that does not support `cloud_init` and is of type `CDROM`.
## Example Usage

```hcl
data "ionoscloud_image" "example" {
  image_alias        = "ubuntu:latest"
  location           = "de/txl"
}
```

Finds an image with alias `ubuntu:latest` in location `de/txl`. Uses exact matching on both fields.
## Example Usage

```hcl
data "ionoscloud_image" "example" {
    type                  = "HDD"
    cloud_init            = "V1"
    image_alias           = "ubuntu:latest"
    location              = "us/ewr"
}
```
Finds an image named `ubuntu-20.04.6` in location `de/txl`. Uses exact matching.

## Argument Reference

 * `name` - (Required) Name of an existing image that you want to search for. It will return an exact match if one exists, otherwise it will retrieve partial matches.
 * `location` - (Optional) Id of the existing image's location. Exact match. Possible values: `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `us/ewr`, `us/las`
 * `type` - (Optional) The image type, HDD or CD-ROM. Exact match.
 * `cloud_init` - (Optional) Cloud init compatibility ("NONE" or "V1"). Exact match.
 * `image_alias` - (Optional) Image alias of the image you are searching for. Exact match. E.g. =`centos:latest`, `ubuntu:latest`
 * `version` - (Optional) The version of the image that you want to search for.

If both "name" and "version" are provided the plugin will concatenate the two strings in this format [name]-[version].
The resulting string will be used to search for an exact match. An error will be thrown if one is not found.

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
 * `licence_type` - OS type of this Image
 * `public` - Indicates if the image is part of the public repository or not
 * `image_aliases` - List of image aliases mapped for this Image
 * `cloud_init` - Cloud init compatibility
 * `type` - This indicates the type of image
 * `location` - Location of that image/snapshot.
