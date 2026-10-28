---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_snapshot_locations_v2"
sidebar_current: "docs-datasource-mongo_snapshot_locations_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Snapshot Locations.
---

# ionoscloud_mongo_snapshot_locations_v2

The `ionoscloud_mongo_snapshot_locations_v2` data source can be used to retrieve the supported MongoDB V2 snapshot locations in a given region. These are the Object Storage locations that can be used for the `snapshot.location` attribute of a cluster.

## Example Usage

```hcl
data "ionoscloud_mongo_snapshot_locations_v2" "example" {
  location = "gb/bhx"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `gb/bhx`, `us/ewr`.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of snapshot locations. Each item exposes:
  * `id` - The UUID of the snapshot location.
  * `location` - The Object Storage location name.
