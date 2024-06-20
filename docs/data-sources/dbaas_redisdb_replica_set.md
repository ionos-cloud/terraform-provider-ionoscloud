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

* `dns_name` - The DNS name of the RedisDB Replica Set.
* `connections` - The network connection for your Replica Set. Only one connection is allowed. It includes:
    * `cidr` - The IP and subnet for your Replica Set.
    * `datacenter_id` - The datacenter to connect your Replica Set to.
    * `lan_id` - The numeric LAN ID to connect your Replica Set to.
* `credentials` - The credentials for your Replica Set. It includes:
    * `username` - The username for your Replica Set.
* `eviction_policy` - The eviction policy of the RedisDB Replica Set.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur. It includes:
    * `time` - Start of the maintenance window in UTC time.
    * `day_of_the_week` - The name of the week day.
* `persistence_mode` - The persistence mode of the RedisDB Replica Set.
* `redis_version` - The version of Redis used in the Replica Set.
* `replicas` - The number of replicas in the Replica Set.
* `resources` - The resources allocated to the Replica Set. It includes:
    * `cores` - The number of CPU cores per instance.
    * `ram` - The amount of memory per instance in gigabytes (GB).
    * `storage` - The amount of storage per instance in gigabytes (GB).