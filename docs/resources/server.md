---
layout: "ionoscloud"
page_title: "IonosCloud: server"
sidebar_current: "docs-resource-server"
description: |-
  Creates and manages IonosCloud Server objects.
---

# ionoscloud_server

Manages a Server on IonosCloud.

## Example Usage

This resource will create an operational server. After this section completes, the provisioner can be called.

```hcl
resource "ionoscloud_server" "example" {
  name              = "server"
  datacenter_id     = "${ionoscloud_datacenter.example.id}"
  cores             = 1
  ram               = 1024
  availability_zone = "ZONE_1"
  cpu_family        = "AMD_OPTERON"
  image_password    = "test1234"
  ssh_key_path      = "${var.private_key_path}"
  image_name        = "${var.ubuntu}"

  volume {
    name           = "new"
    size           = 5
    disk_type      = "SSD"
  }

  nic {
    lan             = "${ionoscloud_lan.example.id}"
    dhcp            = true
    ips             = ["${ionoscloud_ipblock.example.ips[0]}", "${ionoscloud_ipblock.example.ips[1]}"]
    firewall_active = true
  }
}
```

##Argument reference

- `template_uuid` - (Optional)[string] The UUID of the template for creating a CUBE server; the available templates for CUBE servers can be found on the templates resource
- `name` - (Required)[string] The name of the server.
- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `cores` - (Required)[integer] Number of server CPU cores.
- `ram` - (Required)[integer] The amount of memory for the server in MB.
- `image_name` - (Optional)[string] The name, ID or alias of the image. May also be a snapshot ID. It is required if `licence_type` is not provided. Attribute is immutable.
- `availability_zone` - (Optional)[string] The availability zone in which the server should exist. This property is immutable.
- `licence_type` - (Optional)[string] Sets the OS type of the server.
- `cpu_family` - (Optional)[string] Sets the CPU type. "AMD_OPTERON", "INTEL_XEON" or "INTEL_SKYLAKE".
- `volume` - (Required) See the Volume section.
- `nic` - (Required) See the NIC section.
- `boot_volume` - (Computed) The associated boot volume.
- `boot_cdrom` - (Optional)[string] The associated boot drive, if any.
- `boot_image` - (Optional)[string] The image or snapshot UUID / name. May also be an image alias. It is required if `licence_type` is not provided.
- `primary_nic` - (Computed) The associated NIC.
- `primary_ip` - (Computed) The associated IP address.
- `firewallrule_id` - (Computed) The associated firewall rule.
- `ssh_key_path` - (Optional)[list] List of paths to files containing a public SSH key that will be injected into IonosCloud provided Linux images. Required for IonosCloud Linux images. Required if `image_password` is not provided.
- `image_password` - (Optional)[string] Required if `ssh_key_path` is not provided.
- `type` - (Optional)[string] server usages: ENTERPRISE or CUBE. This property is immutable.

*note: image_name under volume level is deprecated, please use image_name under server level*

## Import

Resource Server can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_server.myserver {datacenter uuid}/{server uuid}/{primary_nic uuid}
# or
terraform import ionoscloud_server.myserver {datacenter uuid}/{server uuid}/{primary_nic uuid}/{firewall uuid}
```

## Notes

Please note that for any secondary volume, you need to set the **licence_type** property to **UNKNOWN**
