---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_server"
sidebar_current: "docs-ionoscloud-datasource-server"
description: |-
Get information on a Ionos Cloud Servers
---

# ionoscloud\_server

The server data source can be used to search for and return existing servers.

## Example Usage

```hcl
data "ionoscloud_server" "server_example" {
  name     = "My Server"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing server that you want to search for.
* `id` - (Optional) ID of the server you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the server
* `name` - The name of the server
* `datacenter_id`- The ID of the Virtual Data Center
* `cores` - Number of server CPU cores
* `cpu_family`-  CPU architecture on which server gets provisioned
* `ram` - The amount of memory for the server in MB
* `availability_zone` - The availability zone in which the server should exist
* `vm_state` - Status of the virtual Machine
* `boot_cdrom`
* `boot_volume`
* `cdroms` - list of
  * `id` - Id of the attached cdrom
  * `name` - The name of the attached cdrom
  * `description` - Description of cdrom
  * `location` - Location of that image/snapshot
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
  * `image_type` - Type of image
  * `image_aliases` - List of image aliases mapped for this Image
  * `public` - Indicates if the image is part of the public repository or not
* `volumes` - list of
  * `id` - Id of the attached volume
  * `name` - Name of the attached volume
  * `type` - Hardware type of the volume.
  * `size` - The size of the volume in GB
  * `availability_zone` - The availability zone in which the volume should exist
  * `image` - Image or snapshot ID to be used as template for this volume
  * `image_alias` - List of image aliases mapped for this Image
  * `image_password` - Initial password to be set for installed OS
  * `ssh_keys` - Public SSH keys are set on the image as authorized keys for appropriate SSH login to the instance using the corresponding private key
  * `bus` - The bus type of the volume
  * `licence_type` - OS type of this volume
  * `cpu_hot_plug` - Is capable of CPU hot plug (no reboot required)
  * `ram_hot_plug` - Is capable of memory hot plug (no reboot required)
  * `nic_hot_plug` - Is capable of nic hot plug (no reboot required)
  * `nic_hot_unplug` - Is capable of nic hot unplug (no reboot required)
  * `disc_virtio_hot_plug` - Is capable of Virt-IO drive hot plug (no reboot required)
  * `disc_virtio_hot_unplug` - Is capable of Virt-IO drive hot unplug (no reboot required)
  * `device_number` - The Logical Unit Number of the storage volume
* `nics` - list of
  * `id` - Id of the attached nic
  * `name` - Name of the attached nid
  * `mac` - The MAC address of the NIC
  * `ips` - Collection of IP addresses assigned to a nic
  * `dhcp` - Indicates if the nic will reserve an IP using DHCP
  * `lan` - The LAN ID the NIC will sit on
  * `firewall_active` - Activate or deactivate the firewall
  * `nat` - Indicates if NAT is enabled on this NIC. This is now deprecated.
  * `firewall_rules` - list of
    * `id` - Id of the firewall rule
    * `name` - Name of the firewall rule
    * `protocol` - he protocol for the rule
    * `source_mac` - Only traffic originating from the respective MAC address is allowed
    * `source_ip` - Only traffic originating from the respective IPv4 address is allowed. Value null allows all source IPs
    * `target_ip` - In case the target NIC has multiple IP addresses, only traffic directed to the respective IP address of the NIC is allowed
    * `icmp_code` - Defines the allowed code (from 0 to 254) if protocol ICMP is chosen
    * `icmp_type` - Defines the allowed type (from 0 to 254) if the protocol ICMP is chosen
    * `port_range_start` - Defines the start range of the allowed port (from 1 to 65534) if protocol TCP or UDP is chosen
    * `port_range_end` - Defines the end range of the allowed port (from 1 to 65534) if the protocol TCP or UDP is chosen
