---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_snapshots_v2"
sidebar_current: "docs-datasource-inmemorydb_snapshots_v2"
description: |-
  Lists IONOS Cloud In-Memory DB v2 Snapshots.
---

# ionoscloud_inmemorydb_snapshots_v2 (Data Source)

Lists IONOS Cloud In-Memory DB v2 Snapshots, with an optional cluster ID filter.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_snapshots_v2" "all" {
  location = "de/txl"
}

data "ionoscloud_inmemorydb_snapshots_v2" "for_cluster" {
  location          = "de/txl"
  filter_cluster_id = "e69b22a5-8fee-56b1-b6fb-4a07e4205ead"
}
```

## Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | The location to query. |
| `filter_cluster_id` | String | No | Filter snapshots by cluster UUID. |

## Attributes Reference

Each item in `items` contains:

| Attribute | Description |
|-----------|-------------|
| `id` | The UUID of the snapshot. |
| `location` | The location of the snapshot. |
| `cluster_id` | The ID of the cluster this snapshot belongs to. |
| `datacenter_id` | The ID of the data center where the snapshot was created. |
| `earliest_recovery_target_time` | The earliest time for which a snapshot is available to restore from (RFC3339). |
| `latest_recovery_target_time` | The most recent time for which a snapshot is available to restore from (RFC3339). Empty if available up to the current time. |
| `snapshot_location` | The Object Storage location where the snapshot is stored. |
| `cluster_version` | The In-Memory DB version of the cluster at the time of the snapshot. |
| `snapshot_size` | The size of the snapshot in GB. |
| `required_size_for_restore` | The minimum storage size in GB required to restore from this snapshot. |
