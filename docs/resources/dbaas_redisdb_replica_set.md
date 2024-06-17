---
subcategory: "Database as a Service - RedisDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_redis_replicaset"
sidebar_current: "docs-resource-redis_replicaset"
description: |-
  Creates and manages DBaaS RedisDB Replica Set objects.
---

# ionoscloud_redis_replicaset

Manages a **DBaaS RedisDB Replica Set**.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for DBaaS RedisDB replica sets"
}

resource "ionoscloud_lan"  "example" {
  datacenter_id           = ionoscloud_datacenter.example.id 
  public                  = false
  name                    = "example"
}

resource "ionoscloud_server" "example" {
  name                    = "example"
  datacenter_id           = ionoscloud_datacenter.example.id
  cores                   = 2
  ram                     = 2048
  availability_zone       = "ZONE_1"
  cpu_family              = "INTEL_SKYLAKE"
  image_name              = "debian-10-genericcloud-amd64-20240114-1626"
  image_password          = "password"
  volume {
    name                  = "example"
    size                  = 6
    disk_type             = "SSD Standard"
  }
  nic {
    lan                   = ionoscloud_lan.example.id
    name                  = "example"
    dhcp                  = true
  }
}

locals {
 prefix                   = format("%s/%s", ionoscloud_server.example.nic[0].ips[0], "24")
 database_ip              = cidrhost(local.prefix, 1)
 database_ip_cidr         = format("%s/%s", local.database_ip, "24")
}

resource "ionoscloud_redis_replicaset" "example" {
  location = ionoscloud_datacenter.example.location
  display_name = "ExampleReplicaSet"
  redis_version = "7.2"
  replicas = 4
  resources {
    cores = 1
    ram = 6
  }
  persistence_mode = "RDB"
  eviction_policy = "noeviction"
  connections   {
    datacenter_id         =  ionoscloud_datacenter.example.id
    lan_id                =  ionoscloud_lan.example.id
    cidr                  =  local.database_ip_cidr
  }
  maintenance_window {
    day_of_the_week       = "Monday"
    time                  = "10:00:00"
  }
  credentials {
    username = "myuser"
    plain_text_password = "testpassword"
  }
}
```

## Argument Reference
* `display_name` - (Required)[string] The human readable name of your replica set.
* `location` - (Required)[string] The location of your replica set. Updates to the value of the field force the replica set to be re-created.
* `redis_version` - (Required)[string] The RedisDB version of your replica set.
* `replicas` - (Required)[int] The total number of replicas in the replica set (one active and n-1 passive). In case of a standalone instance, the value is 1. In all other cases, the value is > 1. The replicas will not be available as read replicas, they are only standby for a failure of the active instance.
* `resources` - (Required)[object] The resources of the individual replicas.
  * `cores` - (Required)[int] The number of CPU cores per instance.
  * `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB).
  * `storage` - (Computed)[int] The size of the storage in GB. The size is derived from the amount of RAM and the persistence mode and is not configurable.
* `persistence_mode` - (Required)[string] Specifies How and If data is persisted, possible values are:
  * `None` - Data is inMemory only and will not be persisted. Useful for cache only applications.
  * `AOF` - (Append Only File) AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the Redis protocol itself.
  * `RDB` - (Redis Database) RDB persistence performs snapshots of the current in memory state.
  * `RDB_AOF` - Booth, RDB and AOF persistence are enabled.
* `eviction_policy` - (Required)[string] The eviction policy for the replica set, possible values are:
  * `noeviction` - No eviction policy is used. Redis will never remove any data. If the memory limit is reached, an error will be returned on write operations.
  * `allkeys-lru` - The least recently used keys will be removed first.
  * `allkeys-lfu` - The least frequently used keys will be removed first.
  * `allkeys-random` - Random keys will be removed.
  * `volatile-lru` - The least recently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-lfu` - The least frequently used keys will be removed first, but only among keys with the `expire` field set to `true`.
  * `volatile-random` - Random keys will be removed, but only among keys with the `expire` field set to `true`.
  * `volatile-ttl` - The key with the nearest time to live will be removed first, but only among keys with the `expire` field set to `true`.
* `connections` - (Required)[object] The network connection for your replica set. Only one connection is allowed.
  * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
  * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
  * `cidr` - (Required)[string] The IP and subnet for your instance. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24.
* `credentials` - (Required)[object] Credentials for the Redis replicaset, only one type of password can be used since they are mutually exclusive
  * `username` - (Required)[string] The username for the initial RedisDB user. Some system usernames are restricted (e.g. 'admin', 'standby').
  * `plain_text_password` - (Optional)[string] The password for a RedisDB user, this is a field that is marked as `Sensitive`.
  * `hashed_password` - (Optional)[object] The hashed password for a RedisDB user.
    * `algorithm` - (Required)[string] The value can be only: "SHA-256".
    * `hash` - (Required)[string] The hashed password.
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `initial_snapshot_id` - (Optional)[string] The ID of a snapshot to restore the replica set from. If set, the replica set will be created from the snapshot.
* `dns_name` - (Computed)[string] The DNS name pointing to your replica set. Will be used to connect to the active/standalone instance.

## Import

Resource DBaaS RedisDB Replica Set can be imported using the `replicaset_id` and the `location`, separated by `:`, e.g:

```shell
terraform import ionoscloud_redis_replicaset.example {location}:{replicaSet UUID}
```