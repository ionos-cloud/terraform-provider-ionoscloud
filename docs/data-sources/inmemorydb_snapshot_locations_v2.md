---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_snapshot_locations_v2"
sidebar_current: "docs-datasource-inmemorydb_snapshot_locations_v2"
description: |-
  Lists IONOS Cloud In-Memory DB v2 Snapshot Locations.
---

# ionoscloud_inmemorydb_snapshot_locations_v2 (Data Source)

Lists all IONOS Cloud In-Memory DB v2 Snapshot Locations available for a given API endpoint.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_snapshot_locations_v2" "example" {
  location = "de/txl"
}

output "available_snapshot_locations" {
  value = data.ionoscloud_inmemorydb_snapshot_locations_v2.example.items[*].snapshot_region
}
```

## Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | The InMemoryDB API endpoint location to query. |

## Attributes Reference

| Attribute | Description |
|-----------|-------------|
| `items` | List of snapshot location objects, each with `id` and `snapshot_region`. |
