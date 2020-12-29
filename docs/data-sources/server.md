---
layout: "ionoscloud"
page_title: "ProfitBricks : ionoscloud_server"
sidebar_current: "docs-ionoscloud-datasource-server"
description: |-
Get information on a Ionos Cloud Servers
---

# ionoscloud\_server

The lans data source can be used to search for and return existing servers.

## Example Usage

```hcl
data "ionoscloud_server" "server_example" {
  name     = "My Server"
}
```

## Argument Reference

* `name` - (Optional) Name or part of the name of an existing server that you want to search for.
* `id` - (Optional) ID of the server you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `datacenter_id`
* `cores`
* `ram`
* `availability_zone`
* `vm_state`
* `cpu_family`
* `boot_cdrom`
* `boot_volume`
* `cdroms` - list of
  * `id`
  * `name`
  * `description`
  * `location`
  * `size`
  * `cpu_hot_plug`
  * `cpu_hot_unplug`
  * `ram_hot_plug`
  * `ram_hot_unplug`
  * `nic_hot_plug`
  * `nic_hot_unplug`
  * `disc_virtio_hot_plug`
  * `disc_virtio_hot_unplug`
  * `disc_scsi_hot_plug`
  * `disc_scsi_hot_unplug`
  * `licence_type`
  * `image_type`
  * `image_aliases`
  * `public`
* `volumes` - list of
  * `id`
  * `name`
  * `type`
  * `size`
  * `availability_zone`
  * `image`
  * `image_alias`
  * `image_password`
  * `ssh_keys`
  * `bus`
  * `licence_type`
  * `cpu_hot_plug`
  * `cpu_hot_unplug`
  * `ram_hot_plug`
  * `ram_hot_unplug`
  * `nic_hot_plug`
  * `nic_hot_unplug`
  * `disc_virtio_hot_plug`
  * `disc_virtio_hot_unplug`
  * `disc_scsi_hot_plug`
  * `disc_scsi_hot_unplug`
  * `device_number`
* `nics` - list of
  * `id`
  * `name`
  * `mac`
  * `ips`
  * `dhcp`
  * `lan`
  * `firewall_active`
  * `nat`
  * `firewall_rules` - list of 
    * `id`
    * `name`
    * `protocol`
    * `source_mac`
    * `source_ip`
    * `target_ip`
    * `icmp_code`
    * `icmp_type`
    * `port_range_start`
    * `port_range_end`