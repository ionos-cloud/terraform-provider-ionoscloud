---
subcategory: "Database as a Service - RedisDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_redis_replicaset"
sidebar_current: "docs-datasource-redis_replicaset"
description: |-
  Gets information about an existing RedisDB Replica Set.
---

# ionoscloud_redis_replicaset

The `ionoscloud_redis_replicaset` data source can be used to retrieve information about an existing RedisDB Replica Set.

## Example Usage

### By id
```hcl
data "ionoscloud_redis_replicaset" "example" {
  id = "example-id"
  location = "es/vit"
}
```

### By display_name
```hcl
data "ionoscloud_redis_replicaset" "example" {
  display_name = "example-id"
  location = "us/las"
}
```

* `id` - (Optional) The ID of the RedisDB Replica Set.
* `display_name` - (Optional) The display name of the RedisDB Replica Set.
* `location` - (Required) The location of the RedisDB Replica Set.

> **Note:** Either `id` or `display_name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `dns_name` - (Computed)[string] The DNS name pointing to your replica set. Will be used to connect to the active/standalone instance.
* `connections` - (Required)[object] The network connection for your replica set. Only one connection is allowed. It includes:
  * `cidr` - (Required)[string] The IP and subnet for your instance. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24.
  * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
  * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
* `credentials` - (Required)[object] Credentials for the Redis replicaset, only one type of password can be used since they are mutually exclusive. It includes:
  * `username` - (Required)[string] The username for the initial RedisDB user. Some system usernames are restricted (e.g. 'admin', 'standby').
* `eviction_policy` - (Required)[string] The eviction policy for the replica set, possible values are:
  * `noeviction` - No eviction policy is used. Redis will never remove any data. If the memory limit is reached, an error will be returned on write operations.
  * `allkeys-lru` - The least recently used keys will be removed first.
  * `allkeys-lfu` - The least frequently used keys will be removed first.
  * `allkeys-random` - Random keys will be removed.
  * `volatile-lru` - The least recently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-lfu` - The least frequently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-random` - Random keys will be removed, but only among keys with the `expire` field set to `true`.
  * `volatile-ttl` - The key with the nearest time to live will be removed first, but only among keys with the `expire` field set to `true`.
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur. It includes:
  * `time` - (Required)[string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `persistence_mode` - (Required)[string] Specifies How and If data is persisted, possible values are:
  * `None` - Data is inMemory only and will not be persisted. Useful for cache only applications.
  * `AOF` - (Append Only File) AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the Redis protocol itself.
  * `RDB` - (Redis Database) RDB persistence performs snapshots of the current in memory state.
  * `RDB_AOF` - Both, RDB and AOF persistence are enabled.
* `redis_version` - (Required)[string] The RedisDB version of your replica set.
* `replicas` - (Required)[int] The total number of replicas in the replica set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is > 1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance.
* `resources` - (Required)[object] The resources of the individual replicas. It includes:
  * `cores` - (Required)[int] The number of CPU cores per instance.
  * `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB).
  * `storage` - (Computed)[int] The size of the storage in GB. The size is derived from the amount of RAM and the persistence mode and is not configurable.
