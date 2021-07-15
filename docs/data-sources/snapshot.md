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

 * `name` - (Required) Name of an existing snapshot that you want to search for.
 * `location` - (Optional) Id of the existing snapshot's location.
 * `size` - (Optional) The size of the snapshot to look for.

## Attributes Reference

The following attributes are returned by the datasource:

 * `id` - UUID of the snapshot
