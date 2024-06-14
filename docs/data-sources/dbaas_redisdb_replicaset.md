---
subcategory: "Database as a Service - RedisDB"
layout: "ionoscloud"
page_title: "IonosCloud: dbaas_redisdb_replicaset"
sidebar_current: "docs-ionoscloud_dbaas_redisdb_replicaset"
description: |-
  Gets information about an existing RedisDB Replicaset.
---

# ionoscloud_dbaas_redisdb_replicaset

The `ionoscloud_dbaas_redisdb_replicaset` data source can be used to retrieve information about an existing RedisDB Replicaset.

## Example Usage

### By id
```hcl
data "ionoscloud_dbaas_redisdb_replicaset" "example" {
  id = "example-id"
  location = "es/vit"
}
```

### By display_name
```hcl
data "ionoscloud_dbaas_redisdb_replicaset" "example" {
  display_name = "example-id"
  location = "us/las"
}
```

* `id` - (Optional) The ID of the RedisDB Replicaset.
* `display_name` - (Optional) The display name of the RedisDB Replicaset.
* `location` - (Required) The location of the RedisDB Replicaset.

> **Note:** Either `id` or `display_name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `dns_name` - The DNS name of the RedisDB Replicaset.
* `connections` - The network connection for your Replicaset. Only one connection is allowed. It includes:
    * `cidr` - The IP and subnet for your Replicaset.
    * `datacenter_id` - The datacenter to connect your Replicaset to.
    * `lan_id` - The numeric LAN ID to connect your Replicaset to.
* `credentials` - The credentials for your Replicaset. It includes:
    * `username` - The username for your Replicaset.
* `eviction_policy` - The eviction policy of the RedisDB Replicaset.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur. It includes:
    * `time` - Start of the maintenance window in UTC time.
    * `day_of_the_week` - The name of the week day.
* `persistence_mode` - The persistence mode of the RedisDB Replicaset.
* `redis_version` - The version of Redis used in the Replicaset.
* `replicas` - The number of replicas in the Replicaset.
* `resources` - The resources allocated to the Replicaset. It includes:
    * `cores` - The number of CPU cores per instance.
    * `ram` - The amount of memory per instance in gigabytes (GB).
    * `storage` - The amount of storage per instance in gigabytes (GB).