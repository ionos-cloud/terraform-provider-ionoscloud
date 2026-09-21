---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_snapshot_locations_v2"
sidebar_current: "docs-datasource-inmemorydb_snapshot_locations_v2"
description: |-
  Lists IONOS CLOUD InMemoryDB V2 Snapshot Locations.
---

# ionoscloud_inmemorydb_snapshot_locations_v2

The `ionoscloud_inmemorydb_snapshot_locations_v2` data source can be used to retrieve available InMemoryDB V2 snapshot locations for a given API endpoint.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_snapshot_locations_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Requests are routed to the corresponding regional InMemoryDB endpoint.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of snapshot location objects. Each item includes:
  * `id` - The ID of the snapshot location.
  * `snapshot_region` - The snapshot region identifier.
