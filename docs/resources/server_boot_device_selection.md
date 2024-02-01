---
layout: "ionoscloud"
page_title: "IonosCloud: Server Boot Device Selection"
sidebar_current: "docs-resource-server-boot-device-selection"
description: |-
  Manages the selection of boot devices for IonosCloud Server objects.
---

# ionoscloud_server_boot_device_selection

Manages the selection of a boot device for IonosCloud Servers. 

## Example Usage

The boot device of a `ionoscloud_server`, `ionoscloud_vcpu_server` or `ionoscloud_cube_server` can be selected with this resource.
Deleting this resource will revert the boot device back to the default volume, which is the first inline volume created together with the server.
This resource also allows switching between a `volume` and a `ionoscloud_image` CDROM. Note that CDROM images are detached after they are no longer set as boot devices.

### Select an external volume
```hcl
resource "ionoscloud_server_boot_device_selection" "example"{
  datacenter_id  = ionoscloud_datacenter.example.id
  server_id      = ionoscloud_server.example.id
  boot_device_id = ionoscloud_volume.example.id
}

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  cores             = 2
  ram               = 2048
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name = "Inline Updated"
    size = 20
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "example" {
  server_id = ionoscloud_server.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  name = "External 1"
  size = 10
  disk_type = "HDD"
  availability_zone = "AUTO"
  image_name = "debian:latest"
  image_password = random_password.server_image_password.result
}
```

### Select an inline volume again
```hcl
resource "ionoscloud_server_boot_device_selection" "example"{
  datacenter_id  = ionoscloud_datacenter.example.id
  server_id      = ionoscloud_server.example.id
  boot_device_id = ionoscloud_server.example.inline_volume_ids[0]
}

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  cores             = 2
  ram               = 2048
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name = "Inline Updated"
    size = 20
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "example" {
  server_id = ionoscloud_server.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  name = "External 1"
  size = 10
  disk_type = "HDD"
  availability_zone = "AUTO"
  image_name = "debian:latest"
  image_password = random_password.server_image_password.result
}
```

### Select a CDROM image
```hcl
resource "ionoscloud_server_boot_device_selection" "example"{
  datacenter_id  = ionoscloud_datacenter.example.id
  server_id      = ionoscloud_server.example.inline_volume_ids[0]
  boot_device_id = data.ionoscloud_image.example.id
}

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  cores             = 2
  ram               = 2048
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name = "Inline Updated"
    size = 20
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "example" {
  server_id = ionoscloud_server.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  name = "External 1"
  size = 10
  disk_type = "HDD"
  availability_zone = "AUTO"
  image_name = "debian:latest"
  image_password = random_password.server_image_password.result
}

data "ionoscloud_image" "example" {
  name = "ubuntu-20.04"
  location = "de/txl"
  type = "CDROM"
}
```

### Perform a network boot
```hcl
resource "ionoscloud_server_boot_device_selection" "example"{
  datacenter_id  = ionoscloud_datacenter.example.id
  server_id      = ionoscloud_server.example.inline_volume_ids[0]
  # boot_device_id = data.ionoscloud_image.example.id   VM will boot in the PXE shell when boot_device_id is omitted
}

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  cores             = 2
  ram               = 2048
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name = "Inline Updated"
    size = 20
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "example" {
  server_id = ionoscloud_server.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  name = "External 1"
  size = 10
  disk_type = "HDD"
  availability_zone = "AUTO"
  image_name = "debian:latest"
  image_password = random_password.server_image_password.result
}

data "ionoscloud_image" "example" {
  name = "ubuntu-20.04"
  location = "de/txl"
  type = "CDROM"
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `server_id` - (Required)[string] The ID of a server.
- `boot_device_id` - (Optional)[string] The ID of a bootable device such as a volume or an image data source. If this field is omitted from the configuration, the VM will be restarted with no primary boot device, and it will enter the PXE shell for network booting. 
***Note***: If the network booting process started by the PXE shell fails, the VM will still boot into the image of the attached storage as a fallback. This behavior imitates the "Boot from Network" option from [DCD](https://dcd.ionos.com/).

