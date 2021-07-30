---
layout: "ionoscloud"
page_title: "IonosCloud : autoscaling_template"
sidebar_current: "docs-datasource-autoscaling_template"
description: |-
Get information on a IonosCloud Autoscaling Template
---

# ionoscloud\_autoscaling_template

The autoscaling template data source can be used to search for and return an existing Autoscaling Template. You can provide a string for the name or id parameters which will be compared with provisioned Autoscaling Templates. If a single match is found, it will be returned. If your search results in multiple matches, an error will be generated. When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_autoscaling_template" "autoscaling_template" {
  name			= "test_ds"
}
```

## Argument Reference

* `id` - (Optional) Id of an existing Autoscaling Template that you want to search for.
* `name` - (Optional) Name of an existing Autoscaling Template that you want to search for.

Either `name` or `id` must be provided. If none or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:


* `id` - UUID of the Autoscaling Template.
* `name` - The name of the Autoscaling Template.
* `availability_zone` - Zone where the VMs created using this Template.
* `cores` - The total number of cores for the VMs.
* `cpu_family` - CPU family for the VMs created using this Template. If null, the VM will be created with the default CPU family from the assigned location.
* `location` - Location of the Template.
* `nics` - List of NICs associated with this Template.
    * `lan` - Lan Id for this template Nic.
    * `name` - Name for this template Nic.
* `ram` - The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.
* `volumes` - List of volumes associated with this Template. Only a single volume is currently supported.
    * `image` - Image installed on the volume. Only UUID of the image is supported currently.
    * `image_password` - Image password for this template volume.
    * `name` - Name of the template volume.
    * `size` - User-defined size for this template volume in GB.
    * `ssh_keys` - Ssh keys that has access to the volume.
    * `type` - Storage Type for this template volume (SSD or HDD).
    * `user_data` - user-data (Cloud Init) for this template volume.