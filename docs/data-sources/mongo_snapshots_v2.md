---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_snapshots_v2"
sidebar_current: "docs-datasource-mongo_snapshots_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Snapshots.
---

# ionoscloud_mongo_snapshots_v2

The `ionoscloud_mongo_snapshots_v2` data source can be used to retrieve information about MongoDB V2 snapshots in a given location, with an optional filter by cluster ID.

## Example Usage

```hcl
data "ionoscloud_mongo_snapshots_v2" "all" {
  location = "gb/bhx"
}

data "ionoscloud_mongo_snapshots_v2" "by_cluster" {
  location   = "gb/bhx"
  cluster_id = "example-cluster-id"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `gb/bhx`, `us/ewr`.
* `cluster_id` - (Optional)[string] Filter snapshots by the cluster they belong to.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of snapshots. Each item exposes:
  * `id` - The UUID of the snapshot.
  * `cluster_id` - The ID of the cluster this snapshot belongs to.
  * `cluster_name` - The name of the cluster this snapshot belongs to.
  * `mongo_version` - The MongoDB version of the cluster at snapshot time.
  * `location` - The Object Storage location where the snapshot is stored.
  * `total_size` - The total size of the snapshot in Mebibytes (MiB).
  * `required_storage_size` - The minimum storage size required to restore this snapshot, in Mebibytes (MiB).
  * `earliest_recovery_target_time` - The earliest point in time to which the cluster can be restored from this snapshot.
  * `latest_recovery_target_time` - The latest point in time to which the cluster can be restored. Null if the snapshot can be restored up to the current time.
