---
layout: "ionoscloud"
page_title: "IonosCloud: snapshot"
sidebar_current: "docs-datasource-snapshot"
description: |-
  Get information on a IonosCloud Snapshots
---

# ionoscloud\_snapshot

The snapshots data source can be used to search for and return an existing snapshot which can then be used to provision a server.

## Example Usage

```hcl
data "ionoscloud_snapshot" "snapshot_example" {
  name     = "my snapshot"
  size     = "2"
  location = "location_id"
}
```

## Argument Reference

 * `id` - (Optional) Uuid of an existing snapshot that you want to search for.
 * `name` - (Required) Name of an existing snapshot that you want to search for.
 * `location` - (Optional) Id of the existing snapshot's location.
 * `size` - (Optional) The size of the snapshot to look for.


## Attributes Reference

The following attributes are returned by the datasource:

* `id` - UUID of the snapshot
* `name` - The name of the snapshot.
* `description` - Human readable description
* `licence_type` - OS type of this Snapshot
* `location` - Location of that image/snapshot
* `size` - The size of the image in GB
* `sec_auth_protection` - Boolean value representing if the snapshot requires extra protection e.g. two factor protection
* `cpu_hot_plug` -  Is capable of CPU hot plug (no reboot required)
* `cpu_hot_unplug` -  Is capable of CPU hot unplug (no reboot required)
* `ram_hot_plug` -  Is capable of memory hot plug (no reboot required)
* `ram_hot_unplug` -  Is capable of memory hot unplug (no reboot required)
* `nic_hot_plug` -  Is capable of nic hot plug (no reboot required)
* `nic_hot_unplug` -  Is capable of nic hot unplug (no reboot required)
* `disc_virtio_hot_plug` -  Is capable of Virt-IO drive hot plug (no reboot required)
* `disc_virtio_hot_unplug` -  Is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.
* `disc_scsi_hot_plug` -  Is capable of SCSI drive hot plug (no reboot required)
* `disc_scsi_hot_unplug` -  Is capable of SCSI drive hot unplug (no reboot required). This works only for non-Windows virtual Machines.
