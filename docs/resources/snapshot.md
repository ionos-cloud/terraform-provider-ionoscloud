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
