---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_gpu_server"
sidebar_current: "docs-resource-gpu_server"
description: |-
  Creates and manages IonosCloud GPU Server objects.
---

# ionoscloud_gpu_server

A GPU Server is a Virtual Machine (VM) provisioned from a GPU-enabled template.

Check out the [docs page](https://docs.ionos.com/cloud/compute-services/compute-engine/cloud-gpu-vm)

## Example Usage

This resource will create an operational server. After this section completes, the provisioner can be called.

### GPU Server

```hcl
resource "ionoscloud_datacenter" "example" {
  name     = "Datacenter Example"
  location = "de/fra/2"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "de/fra"
  size     = 1
  name     = "webserver_ipblock"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = true
  name          = "Lan Example"
}

resource "random_password" "server_image_password" {
  length  = 16
  special = false
}

resource "ionoscloud_gpu_server" "example" {
  name              = "GPU Server Example"
  hostname          = "gpu-server-example"
  datacenter_id     = ionoscloud_datacenter.example.id
  availability_zone = "AUTO"

  template_uuid  = "6913ed82-a143-4c15-89ac-08fb375a97c5"
  image_name     = "ubuntu:latest"
  image_password = random_password.server_image_password.result

  vm_state = "RUNNING"

  volume {
    name                = "system"
    licence_type        = "LINUX"
    disk_type           = "SSD Premium"
    bus                 = "VIRTIO"
    availability_zone   = "AUTO"
    expose_serial       = true
    require_legacy_bios = false
  }

  nic {
    lan             = ionoscloud_lan.example.id
    name            = "system"
    dhcp            = true
    firewall_active = true
    firewall_type   = "INGRESS"
    ips             = [ionoscloud_ipblock.webserver_ipblock.ips[0]]

    firewall {
      protocol         = "TCP"
      name             = "SSH"
      port_range_start = 22
      port_range_end   = 22
      type             = "INGRESS"
    }
  }
}
```

## Argument reference

- `template_uuid` - (Required)[string] The UUID of the template used for creating a GPU server.
- `name` - (Required)[string] The name of the server.
- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `hostname` - (Optional)(Computed) The hostname of the resource. Allowed characters are a-z, 0-9 and - (minus). Hostname should not start with minus and should not be longer than 63 characters. If no value provided explicitly, it will be populated with the name of the server.
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `vm_state` - (Optional)[string] Sets the power state of the GPU server. E.g: `RUNNING` or `SUSPENDED`.
- `volume` - (Required) See the [Volume](volume.md) section.
- `nic` - (Optional) See the [Nic](nic.md) section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - ***DEPRECATED*** Please refer to [ionoscloud_server_boot_device_selection](server_boot_device_selection.md) (Optional)[string] The associated boot drive, if any. Must be the UUID of a bootable CDROM image that can be retrieved using the [ionoscloud_image](../data-sources/image.md) data source.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `ssh_key_path` - (Optional)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
- `image_password` - (Optional)[string] Required if `ssh_key_path` is not provided.
- `security_groups_ids` - (Optional) The list of Security Group IDs for the resource.
- `allow_replace` - (Optional)[bool] When set to true, allows the update of immutable fields by first destroying and then re-creating the server.

⚠️ **_Warning: `allow_replace` - lets you update immutable fields, but it first destroys and then re-creates the server in order to do it. This field should be used with care, understanding the risks._**

> **⚠ WARNING**
>
> Image_name under volume level is deprecated, please use image_name under server level

## Import

Resource GPU Server can be imported using the `resource id` and the `datacenter id`, e.g.

```shell
terraform import ionoscloud_gpu_server.myserver datacenter uuid/server uuid
```
