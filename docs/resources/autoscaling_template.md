---
layout: "ionoscloud"
page_title: "IonosCloud: autoscaling_template"
sidebar_current: "docs-resource-autoscaling_template"
description: |-
Creates and manages IonosCloud Autoscaling Template.
---

# ionoscloud_autoscaling_template

Manages an Autoscaling Template on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_autoscaling_template" "autoscaling_template" {
	availability_zone    = "AUTO"
    cores				 = 2
	cpu_family           = "INTEL_SKYLAKE"
	location			 = "de/txl"
    name                 = "autoscaling_template"
    nics    {
		lan              = ionoscloud_lan.autoscaling_template.id
        name             = "test_autoscaling_template"
    }
    ram                  = 1024
	volumes  {
    	image            = "e309f108-b48d-11eb-b9b3-d2869b2d44d9"
		image_password   = "test12345678"
        name             = "test_autoscaling_template"
		size             = 50
    	type             = "HDD"
	}
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional)[string] The name of the Autoscaling Template.
- `availability_zone` - (Required)[string] Zone where the VMs created using this Template.
- `cores` - (Required)[int] The total number of cores for the VMs.
- `cpu_family` - (Required)[string] CPU family for the VMs created using this Template. If null, the VM will be created with the default CPU family from the assigned location.
- `location` - (Required)[string] Location of the Template.
- `nics` - (Optional)[list] List of NICs associated with this Template.
    - `lan` - (Required)[string] Lan Id for this template Nic.
    - `name` - (Required)[string] Name for this template Nic.
- `ram` - (Required)[int] The amount of memory for the VMs in MB, e.g. 2048. Size must be specified in multiples of 256 MB with a minimum of 256 MB; however, if you set ramHotPlug to TRUE then you must use a minimum of 1024 MB. If you set the RAM size more than 240GB, then ramHotPlug will be set to FALSE and can not be set to TRUE unless RAM size not set to less than 240GB.
- `volumes` - (Optional)[list] List of volumes associated with this Template. Only a single volume is currently supported.
    - `image` - (Required)[string] Image installed on the volume. Only UUID of the image is supported currently.
    - `image_password` - (Optional)[string] Image password for this template volume.
    - `name` - (Optional)[string] Name of the template volume.
    - `size` - (Optional)[int] User-defined size for this template volume in GB.
    - `ssh_keys` - (Optional)[list] Ssh keys that has access to the volume.
    - `type` - (Required)[string] Storage Type for this template volume (SSD or HDD).
    - `user_data` - (Optional)[string] user-data (Cloud Init) for this template volume.