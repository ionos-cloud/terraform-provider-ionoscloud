---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-server"
description: |-
  Creates and manages IonosCloud Server objects.
---

# ionoscloud_server

Manages a **Server** on IonosCloud.

## Example Usage

This resource will create an operational server. After this section completes, the provisioner can be called.

### ENTERPRISE Server

```hcl
data "ionoscloud_image" "example" {
    type                  = "HDD"
    cloud_init            = "V1"
    location              = "us/las"
}

resource "ionoscloud_datacenter" "example" {
    name                  = "Datacenter Example"
    location              = "us/las"
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

resource "ionoscloud_server" "example" {
    name                  = "Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    availability_zone     = "ZONE_1"
    cpu_family            = "AMD_OPTERON"
    image_name            = data.ionoscloud_image.example.id
    image_password        = random_password.server_image_password.result
    type                  = "ENTERPRISE"
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

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  type              = "CUBE"
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

- `template_uuid` - (Optional)[string] The UUID of the template for creating a CUBE server; the available templates for CUBE servers can be found on the templates resource
- `name` - (Required)[string] The name of the server.
- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `cores` - (Optional)[integer] Number of server CPU cores.
- `ram` - (Optional)[integer] The amount of memory for the server in MB.
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. E.g: `AUTO`, `ZONE_1`, `ZONE_2`. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `cpu_family` - (Optional)[string] CPU architecture on which server gets provisioned; not all CPU architectures are available in all datacenter regions; available CPU architectures can be retrieved from the datacenter resource. E.g.: "AMD_OPTERON", "INTEL_XEON" or "INTEL_SKYLAKE".
- `volume` - (Required) See the [Volume](volume.md) section.
- `nic` - (Required) See the [Nic](nic.md) section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - (Optional)[string] The associated boot drive, if any.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `ssh_key_path` - (Optional)[list] List of absolute paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images.  Also accepts ssh keys directly. Required for IonosCloud Linux images. Required if `image_password` is not provided. This property is immutable.
- `ssh_keys` - (Optional)[list] Immutable List of absolute paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images.  Also accepts ssh keys directly. Required for IonosCloud Linux images. Required if `image_password` is not provided. This property is immutable.
- `image_password` - (Optional)[string] Required if `ssh_key_path` is not provided.
- `type` - (Optional)[string] Server usages: [ENTERPRISE](https://docs.ionos.com/cloud/compute-engine/virtual-servers/virtual-servers) or [CUBE](https://docs.ionos.com/cloud/compute-engine/virtual-servers/cloud-cubes). This property is immutable.
- `label` - (Optional) A label can be seen as an object with only two required fields: `key` and `value`, both of the `string` type. Please check the example presented above to see how a `label` can be used in the plan. A server can have multiple labels.
- `inline_volume_ids` - (Computed) A list with the IDs for the volumes that are defined inside the server resource.

> **⚠ WARNING** 
> 
> Image_name under volume level is deprecated, please use image_name under server level
> ssh_key_path and ssh_keys fields are immutable.


> **⚠ WARNING**
> 
> If you want to create a **CUBE** server, you have to provide the `template_uuid`. In this case you can not set `cores`, `ram` and `volume.size` arguments, these being mutually exclusive with `template_uuid`.
> 
> In all the other cases (**ENTERPRISE** servers) you have to provide values for `cores`, `ram` and `volume size`.


## Import

Resource Server can be imported using the `resource id` and the `datacenter id`, e.g.. Passing only resource id and datacenter id means that the first nic found linked to the server will be attached to it.

```shell
terraform import ionoscloud_server.myserver {datacenter uuid}/{server uuid}
```
Optionally, you can pass `primary_nic` and `firewallrule_id` so terraform will know to import also the first nic and firewall rule (if it exists on the server):
```shell
terraform import ionoscloud_server.myserver {datacenter uuid}/{server uuid}/{primary nic id}/{firewall rule id}
```

## Notes

Please note that for any secondary volume, you need to set the **licence_type** property to **UNKNOWN**
