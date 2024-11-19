---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_cube_server"
sidebar_current: "docs-resource-cube_server"
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
    name            = "Basic Cube XS"
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

### With IPv6 Enabled

```hcl
data "ionoscloud_template" "example" {
  name            = "Basic Cube XS"
}
resource "ionoscloud_datacenter" "example" {
	name            = "Datacenter Example"
	location        = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "de/txl"
  size = 4
  name = "webserver_ipblock"
}
resource "ionoscloud_lan" "example" {
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "Lan Example"
  ipv6_cidr_block = cidrsubnet(ionoscloud_datacenter.example.ipv6_cidr_block,8,10)
}
resource "ionoscloud_cube_server" "example" {
  name              = "Server Example"
  availability_zone = "AUTO"
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
    ips             = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1]]
    
    dhcpv6          = false
    ipv6_cidr_block = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5)
    ipv6_ips        = [ 
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),1),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),2),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),3)
                      ]

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
- `hostname` - (Optional)(Computed) The hostname of the resource. Allowed characters are a-z, 0-9 and - (minus). Hostname should not start with minus and should not be longer than 63 characters. If no value provided explicitly, it will be populated with the name of the server
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `vm_state` - (Optional)[string] Sets the power state of the cube server. E.g: `RUNNING` or `SUSPENDED`.
- `volume` - (Required) See the [Volume](volume.md) section.
- `nic` - (Required) See the [Nic](nic.md) section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - ***DEPRECATED*** Please refer to [ionoscloud_server_boot_device_selection](server_boot_device_selection.md) (Optional)[string] The associated boot drive, if any. Must be the UUID of a bootable CDROM image that can be retrieved using the [ionoscloud_image](../data-sources/image.md) data source.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `ssh_key_path` - (Optional)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
- `image_password` - (Optional)[string] Required if `ssh_key_path` is not provided.
- `security_groups_ids` - (Optional) The list of Security Group IDs for the resource.

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
