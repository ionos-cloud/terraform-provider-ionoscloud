---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_vcpu_server"
sidebar_current: "docs-resource-vcpu-server"
description: |-
  Creates and manages IonosCloud VCPU Server objects.
---

# ionoscloud_vcpu_server

Manages a **VCPU Server** on IonosCloud.

## Example Usage

### VCPU Server

```hcl
data "ionoscloud_image" "example" {
    type                  = "HDD"
    image_alias           = "ubuntu:latest"
    location              = "us/las"
}

resource "ionoscloud_datacenter" "example" {
    name                  = "Datacenter Example"
    location              = "de/txl"
    description           = "Datacenter Description"
    sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    public                = true
    name                  = "Lan Example"
}

resource "ionoscloud_ipblock" "example" {
    location              = ionoscloud_datacenter.example.location
    size                  = 4
    name                  = "IP Block Example"
}

resource "ionoscloud_vcpu_server" "example" {
    name                  = "VCPU Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    availability_zone     = "ZONE_1"
    image_name            = data.ionoscloud_image.example.id
    image_password        = random_password.server_image_password.result
    volume {
        name              = "system"
        size              = 5
        disk_type         = "SSD Standard"
        user_data         = "foo"
        bus               = "VIRTIO"
        availability_zone = "ZONE_1"
    }
    nic {
        lan               = ionoscloud_lan.example.id
        name              = "system"
        dhcp              = true
        firewall_active   = true
        firewall_type     = "BIDIRECTIONAL"
        ips               = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
        firewall {
          protocol          = "TCP"
          name              = "SSH"
          port_range_start  = 22
          port_range_end    = 22
          source_mac        = "00:0a:95:9d:68:17"
          source_ip         = ionoscloud_ipblock.example.ips[2]
          target_ip         = ionoscloud_ipblock.example.ips[3]
          type              = "EGRESS" 
        }
    }
    label {
        key = "labelkey1"
        value = "labelvalue1"
    }
    label {
        key = "labelkey2"
        value = "labelvalue2"
    }
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
                  
```

## Argument reference

- `name` - (Required)[string] The name of the server.
- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `hostname` - (Optional)[string] The hostname of the  resource. Allowed characters are a-z, 0-9 and - (minus). Hostname should not start with minus and should not be longer than 63 characters.
- `cores` - (Optional)[integer] Number of server CPU cores.
- `ram` - (Optional)[integer] The amount of memory for the server in MB.
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. E.g: `AUTO`, `ZONE_1`, `ZONE_2`. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `volume` - (Required) See the [Volume](volume.md) section.
- `nic` - (Optional) See the [Nic](nic.md) section.
- `firewall` - (Optional) Allows to define firewall rules inline in the server. See the [Firewall](firewall.md) section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - ***DEPRECATED*** Please refer to [ionoscloud_server_boot_device_selection](server_boot_device_selection.md) (Optional)[string] The associated boot drive, if any. Must be the UUID of a bootable CDROM image that can be retrieved using the [ionoscloud_image](./../data-sources/image.md) data source.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `firewallrule_ids` - (Computed) The associated firewall rules.
- `ssh_keys` - (Optional)[list] Immutable List of absolute or relative paths to files containing public SSH key that will be injected into IonosCloud provided Linux images. Also accepts ssh keys directly. Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key. This field may only be set in creation requests. When reading, it always returns null. SSH keys are only supported if a public Linux image is used for the volume creation. Does not support `~` expansion to homedir in the given path.
- `image_password` - (Optional)[string] The password for the image.
- `label` - (Optional) A label can be seen as an object with only two required fields: `key` and `value`, both of the `string` type. Please check the example presented above to see how a `label` can be used in the plan. A server can have multiple labels.
- `inline_volume_ids` - (Computed) A list with the IDs for the volumes that are defined inside the server resource.

> **⚠ WARNING** 
> 
> ssh_keys field is immutable.

## Import

Resource VCPU Server can be imported using the `resource id` and the `datacenter id`, for example, passing only resource id and datacenter id means that the first nic found linked to the server will be attached to it.

```shell
terraform import ionoscloud_vcpu_server.myserver {datacenter uuid}/{server uuid}
```
Optionally, you can pass `primary_nic` and `firewallrule_id` so terraform will know to import also the first nic and firewall rule (if it exists on the server):
```shell
terraform import ionoscloud_vcpu_server.myserver {datacenter uuid}/{server uuid}/{primary nic id}/{firewall rule id}
```

## Notes

Please note that for any secondary volume, you need to set the **licence_type** property to **UNKNOWN**

⚠️ **Note:** Important for deleting an `firewall` rule from within a list of inline resources defined on the same nic. There is one limitation to removing one firewall rule
from the middle of the list of `firewall` rules. Terraform will actually modify the existing rules and delete the last one.
More details [here](https://github.com/hashicorp/terraform/issues/14275). There is a workaround described in the issue 
that involves moving the resources in the list prior to deletion.
`terraform state mv <resource-name>.<resource-id>[<i>] <resource-name>.<resource-id>[<j>]`