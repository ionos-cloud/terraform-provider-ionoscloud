---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: cube server"
sidebar_current: "docs-resource-server"
description: |-
  Creates and manages IonosCloud Cube Server objects.
---

# ionoscloud_cube_server

Manages a **Cube Server** on IonosCloud.

## Example Usage

This resource will create an operational server. After this section completes, the provisioner can be called.

### CUBE Server

```hcl
data "ionoscloud_template" "example" {
    name            = "CUBES XS"
}

resource "ionoscloud_datacenter" "example" {
	name            = "Datacenter Example"
	location        = "de/txl"
}

resource "ionoscloud_lan" "example" {
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "Lan Example"
}

resource "ionoscloud_cube_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  template_uuid     = data.ionoscloud_template.example.id
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name            = "Volume Example"
    licence_type    = "LINUX" 
    disk_type       = "DAS"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```

## Argument reference

- `template_uuid` - (Required)[string] The UUID of the template for creating a CUBE server; the available templates for CUBE servers can be found on the templates resource
- `name` - (Required)[string] The name of the server.
- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `cpu_family` - (Optional)[string] Sets the CPU type. "AMD_OPTERON", "INTEL_XEON" or "INTEL_SKYLAKE".
- `volume` - (Required) See the [Volume](volume.md) section.
- `nic` - (Required) See the [Nic](nic.md) section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - (Optional)[string] The associated boot drive, if any.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `ssh_key_path` - (Optional)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
- `image_password` - (Optional)[string] Required if `ssh_key_path` is not provided.

> **⚠ WARNING** 
> 
> Image_name under volume level is deprecated, please use image_name under server level


> **⚠ WARNING**
> 
> For creating a **CUBE** server, you can not set `volume.size` argument.
>

## Import

Resource Server can be imported using the `resource id` and the `datacenter id`, e.g.

```shell
terraform import ionoscloud_cube_server.myserver {datacenter uuid}/{server uuid}
```

## Notes

Please note that for any secondary volume, you need to set the **licence_type** property to **UNKNOWN**
