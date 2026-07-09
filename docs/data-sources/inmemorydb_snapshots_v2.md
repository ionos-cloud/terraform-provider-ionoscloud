---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_snapshots_v2"
sidebar_current: "docs-datasource-inmemorydb_snapshots_v2"
description: |-
  Lists IONOS CLOUD InMemoryDB V2 Snapshots.
---

# ionoscloud_inmemorydb_snapshots_v2

The `ionoscloud_inmemorydb_snapshots_v2` data source can be used to retrieve information about existing InMemoryDB V2 snapshots, with an optional cluster ID filter.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_snapshots_v2" "all" {
  location = "de/txl"
}

data "ionoscloud_inmemorydb_snapshots_v2" "for_cluster" {
  location   = "de/txl"
  cluster_id = "example-id"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query.
* `cluster_id` - (Optional)[string] Filter snapshots by cluster UUID.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of snapshots. Each item includes:
  * `id` - The UUID of the snapshot.
  * `location` - The location of the snapshot.
  * `cluster_id` - The ID of the cluster this snapshot belongs to.
  * `cluster_name` - The name of the cluster this snapshot belongs to.
  * `datacenter_id` - The ID of the data center where the snapshot was created.
  * `earliest_recovery_target_time` - The earliest time for which a snapshot is available to restore from (RFC3339).
  * `latest_recovery_target_time` - The most recent time for which a snapshot is available to restore from (RFC3339). Empty if available up to the current time.
  * `snapshot_location` - The Object Storage location where the snapshot is stored.
  * `cluster_version` - The InMemoryDB version of the cluster at the time of the snapshot.
  * `snapshot_size` - The size of the snapshot in GB.
  * `required_size_for_restore` - The minimum storage size in GB required to restore from this snapshot.
