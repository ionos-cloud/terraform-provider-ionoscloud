---
subcategory: "Database as a Service - InMemoryDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_replicaset"
sidebar_current: "docs-datasource-inmemorydb_replicaset"
description: |-
  Gets information about an existing InMemoryDB Replica Set.
---

# ionoscloud_inmemorydb_replicaset

The `ionoscloud_inmemorydb_replicaset` data source can be used to retrieve information about an existing InMemoryDB Replica Set.

## Example Usage

### By id
```hcl
data "ionoscloud_inmemorydb_replicaset" "example" {
  id = "example-id"
  location = "es/vit"
}
```

### By display_name
```hcl
data "ionoscloud_inmemorydb_replicaset" "example" {
  display_name = "example-id"
  location = "us/las"
}
```

## Argument Reference

* `id` - (Optional) The ID of the InMemoryDB Replica Set.
* `display_name` - (Optional) The display name of the InMemoryDB Replica Set.
* `location` - (Optional) The location of the InMemoryDB Replica Set.

> **Note:** Either `id` or `display_name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `dns_name` - [string] The DNS name pointing to your replica set. Will be used to connect to the active/standalone instance.
* `connections` - [object] The network connection for your replica set. Only one connection is allowed. It includes:
  * `cidr` - [string] The IP and subnet for your instance. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24.
  * `datacenter_id` - [string] The datacenter to connect your instance to.
  * `lan_id` - [string] The numeric LAN ID to connect your instance to.
* `credentials` - [object] Credentials for the InMemoryDB replicaset, only one type of password can be used since they are mutually exclusive. It includes:
  * `username` - [string] The username for the initial InMemoryDB user. Some system usernames are restricted (e.g. 'admin', 'standby').
* `eviction_policy` - [string] The eviction policy for the replica set, possible values are:
  * `noeviction` - No eviction policy is used. InMemoryDB will never remove any data. If the memory limit is reached, an error will be returned on write operations.
  * `allkeys-lru` - The least recently used keys will be removed first.
  * `allkeys-lfu` - The least frequently used keys will be removed first.
  * `allkeys-random` - Random keys will be removed.
  * `volatile-lru` - The least recently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-lfu` - The least frequently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-random` - Random keys will be removed, but only among keys with the `expire` field set to `true`.
  * `volatile-ttl` - The key with the nearest time to live will be removed first, but only among keys with the `expire` field set to `true`.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur. It includes:
  * `time` - [string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - [string] The name of the week day.
* `persistence_mode` - [string] Specifies How and If data is persisted, possible values are:
  * `None` - Data is inMemory only and will not be persisted. Useful for cache only applications.
  * `AOF` - (Append Only File) AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the InMemoryDB protocol itself.
  * `RDB` - RDB persistence performs snapshots of the current in memory state.
  * `RDB_AOF` - Both, RDB and AOF persistence are enabled.
* `version` - [string] The InMemoryDB version of your replica set.
* `replicas` - [int] The total number of replicas in the replica set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is > 1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance.
* `resources` - [object] The resources of the individual replicas. It includes:
  * `cores` - [int] The number of CPU cores per instance.
  * `ram` - [int] The amount of memory per instance in gigabytes (GB).
  * `storage` - [int] The size of the storage in GB. The size is derived from the amount of RAM and the persistence mode and is not configurable.
