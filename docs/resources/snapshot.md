---
layout: "ionoscloud"
page_title: "IonosCloud: snapshot"
sidebar_current: "docs-resource-snapshot"
description: |-
  Creates and manages snapshot objects.
---

# ionoscloud\_snapshot

Manages snapshots on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_snapshot" "test_snapshot" {
  datacenter_id = "datacenterId"
  volume_id = "volumeId"
  name = "my snapshot"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of the Virtual Data Center.
* `name` - (Required)[string] The name of the snapshot.
* `volume_id` - (Required)[string] The ID of the specific volume to take the snapshot from.

## Attribute reference

Beside the configurable arguments, the resource returns the following additional attributes:

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


## Import

Resource Snapshot can be imported using the `snapshot id`, e.g.

```shell
terraform import ionoscloud_snapshot.mysnapshot {snapshot uuid}
```