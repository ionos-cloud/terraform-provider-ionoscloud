---
subcategory: "Database as a Service - RedisDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_redis_snapshot"
sidebar_current: "docs-datasource-redis_snapshot"
description: |-
  Gets information about an existing RedisDB Snapshot.
---

# ionoscloud_redis_snapshot

The `ionoscloud_redis_snapshot` data source can be used to retrieve information about an existing RedisDB Snapshot.

## Example Usage

```hcl
data "ionoscloud_redis_snapshot" "example" {
  id = "snapshot-id"
  location = "de/txl"
}
```

## Argument Reference

* `id` - (Required) The ID of the RedisDB Snapshot.
* `location` - (Required) The location of the RedisDB Snapshot.

## Attributes Reference

* `metadata` - Metadata of the snapshot.
  * `created_date` - The ISO 8601 creation timestamp.
  * `datacenter_id` - The ID of the datacenter in which the snapshot is located.
  * `last_modified_date` - The ISO 8601 modified timestamp.
  * `replica_set_id` - The ID of the replica set from which the snapshot was created.
  * `snapshot_time` - The time at which the snapshot was taken.