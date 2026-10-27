---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_cluster_v2"
sidebar_current: "docs-resource-mongo_cluster_v2"
description: |-
  Creates and manages an IONOS CLOUD MongoDB V2 Cluster.
---

# ionoscloud_mongo_cluster_v2

Manages a DBaaS MongoDB V2 Cluster.

A cluster has one of two editions, selected with the `edition` attribute:

* `ENTERPRISE` — compute and storage are set directly via the `instances` object; supports `sharding` and `bi_connector_enabled`.
* `BUSINESS` — compute and storage come from a predefined template via `instance_template_id`, and the number of instances via `instance_count`.

## Example Usage — ENTERPRISE edition

```hcl
resource "ionoscloud_datacenter" "example" {
  name     = "example"
  location = "gb/bhx"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "example"
}

resource "random_password" "mongo" {
  length  = 16
  special = true
}

resource "ionoscloud_mongo_cluster_v2" "example" {
  location    = "gb/bhx"
  name        = "my-mongo-cluster"
  description = "MongoDB cluster"
  edition     = "ENTERPRISE"
  version     = "7.0"

  logs_enabled    = true
  metrics_enabled = true

  instances = {
    count        = 1
    cores        = 1
    ram          = 8
    storage_size = 10
  }

  connections = {
    datacenter_id      = ionoscloud_datacenter.example.id
    lan_id             = ionoscloud_lan.example.id
    instance_addresses = ["192.168.1.100/24"]
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 7
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "dbadmin"
    password = random_password.mongo.result
    database = "my_database"
  }
}
```

## Example Usage — BUSINESS edition

```hcl
data "ionoscloud_mongo_templates_v2" "business" {
  location = "gb/bhx"
  edition  = "BUSINESS"
}

resource "random_password" "mongo" {
  length  = 16
  special = true
}

resource "ionoscloud_mongo_cluster_v2" "example" {
  location             = "gb/bhx"
  name                 = "my-business-cluster"
  edition              = "BUSINESS"
  version              = "7.0"
  instance_template_id = data.ionoscloud_mongo_templates_v2.business.items[0].id
  instance_count       = 1

  connections = {
    datacenter_id      = ionoscloud_datacenter.example.id
    lan_id             = ionoscloud_lan.example.id
    instance_addresses = ["192.168.1.100/24"]
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 7
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "dbadmin"
    password = random_password.mongo.result
    database = "my_database"
  }
}
```

## Argument Reference

* `location` - (Required)[string] The location of the cluster. Changing this forces a new resource. Available: `gb/bhx`, `us/ewr`.
* `name` - (Required)[string] The name of the MongoDB cluster. Between 2 and 63 characters; must begin and end with an alphanumeric character, with dashes, underscores, dots and alphanumerics in between.
* `edition` - (Required)[string] The edition of the MongoDB cluster: `ENTERPRISE` or `BUSINESS`. Changing this forces a new resource.
* `version` - (Required)[string] The MongoDB version for the cluster, e.g. `6.0`, `7.0`, `8.0`. See the [`ionoscloud_mongo_versions_v2`](../data-sources/mongo_versions_v2.md) data source for the list of supported versions and their available upgrade paths.
* `description` - (Optional)[string] Human-readable description for the cluster.
* `logs_enabled` - (Optional)(Computed)[bool] Allows or disallows the collection and reporting of logs for this cluster's observability. If the observability service is not activated on the contract, this setting is accepted but has no effect. Defaults to false.
* `metrics_enabled` - (Optional)(Computed)[bool] Allows or disallows the collection and reporting of metrics for this cluster's observability. If the observability service is not activated on the contract, this setting is accepted but has no effect. Defaults to false.
* `instances` - (Optional)[object] Compute and storage configuration for each instance. Required for the `ENTERPRISE` edition; not valid for `BUSINESS`.
  * `count` - (Required)[int] The total number of instances in the cluster, one primary and n-1 secondary (between 1 and 9).
  * `cores` - (Required)[int] The number of CPU cores per instance (between 1 and 62).
  * `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB) (between 6 and 240).
  * `storage_size` - (Required)[int] The amount of storage per instance in gigabytes (GB) (between 10 and 4096).
* `instance_template_id` - (Optional)[string] The ID (UUID) of the instance template defining cores, RAM and storage (see the [`ionoscloud_mongo_templates_v2`](../data-sources/mongo_templates_v2.md) data source). Required for the `BUSINESS` edition; not valid for `ENTERPRISE`. Immutable — changing it forces a new resource.
* `instance_count` - (Optional)[int] The total number of instances in the cluster, one primary and n-1 secondary (between 1 and 9). Required for the `BUSINESS` edition; not valid for `ENTERPRISE`.
* `sharding` - (Optional)[object] Sharding properties. Can only be set for the `ENTERPRISE` edition.
  * `enabled` - (Required)[bool] Activates or deactivates sharding for this cluster.
  * `shards` - (Optional)(Computed)[int] The number of shards to create in the cluster (between 2 and 5). When sharding is enabled, the instances count signifies the number of replicas per shard.
* `bi_connector_enabled` - (Optional)(Computed)[bool] Enables or disables the MongoDB BI Connector. Can only be enabled for the `ENTERPRISE` edition. Defaults to false.
* `connections` - (Required)[object] Connection information of the MongoDB cluster.
  * `datacenter_id` - (Required)[string] The datacenter to connect your instance to. **Immutable** — changing it forces a new resource.
  * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to. **Immutable** — changing it forces a new resource.
  * `instance_addresses` - (Required)[list of string] A list of IP addresses (with netmask, in CIDR notation) assigned to the cluster instances. The number of addresses must match the number of instances in the cluster. The range `10.208.0.0/12` is not available.
* `snapshot` - (Required)[object] Snapshot storage and retention configuration for the cluster.
  * `location` - (Required)[string] The Object Storage location where the snapshots will be created. For added data safety, use a different location than the cluster.
  * `retention_days` - (Required)[int] Configures how many days cluster snapshots are retained before being automatically deleted (between 1 and 365).
* `maintenance_window` - (Required)[object] A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start of the maintenance window in UTC time (`HH:MM:SS`).
  * `day_of_the_week` - (Required)[string] The name of the week day: `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday` or `Sunday`.
* `credentials` - (Required)[object] Credentials for the initial database user to be created.
  * `username` - (Required)[string] The username of the initial MongoDB user. Between 2 and 16 characters, alphanumeric and underscores only. Some usernames are reserved for platform use.
  * `password` - (Required)[string] **Sensitive.** The password for the initial MongoDB user. Between 8 and 256 characters. Not returned by the API — will be null in state after `terraform import`.
  * `database` - (Required)[string] The name of the initial database to be created. Between 2 and 63 characters, alphanumeric and underscores only.
* `restore_from_snapshot` - (Optional)[object] Restores the cluster from a snapshot.
  * `source_snapshot_id` - (Optional)[string] UUID of the snapshot to restore from. Required when `restore_from_snapshot` is set during cluster creation; not valid for in-place restore during an update.
  * `recovery_target_datetime` - (Optional)[string] ISO 8601 timestamp causing the system to replay snapshots up to the specified time. Optional for create-time restore; required for in-place restore during an update.

> **Note:** `restore_from_snapshot` is not returned by the API. The values are stored in state as configured but will be null after `terraform import`.

* `timeouts` - (Optional) Standard Terraform timeouts: `create`, `update`, `delete`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID (UUID) of the cluster.
* `connection_uri` - The connection string used to access the cluster.

## Import

MongoDB V2 clusters can be imported using `<location>:<cluster_id>`:

```shell
terraform import ionoscloud_mongo_cluster_v2.example gb/bhx:example-id
```

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can be used with the `identity` attribute:

```hcl
import {
  to = ionoscloud_mongo_cluster_v2.example
  identity = {
    id       = "example-id"
    location = "gb/bhx"
  }
}

resource "ionoscloud_mongo_cluster_v2" "example" {
  ### Configuration omitted for brevity ###
}
```
