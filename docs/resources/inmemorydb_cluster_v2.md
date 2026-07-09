---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_cluster_v2"
sidebar_current: "docs-resource-inmemorydb_cluster_v2"
description: |-
  Creates and manages an IONOS CLOUD InMemoryDB V2 Cluster.
---

# ionoscloud_inmemorydb_cluster_v2

Manages a DBaaS InMemoryDB V2 Cluster.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name     = "example"
  location = "de/txl"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "example"
}

resource "ionoscloud_inmemorydb_cluster_v2" "example" {
  location         = "de/txl"
  name             = "my-inmemorydb-cluster"
  description      = "InMemoryDB cluster"
  version          = "9.0"
  persistence_mode = "RDB"
  eviction_policy  = "allkeys-lru"
  logs_enabled     = true
  metrics_enabled  = true

  instances = {
    count = 1
    cores = 1
    ram   = 4
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.example.id
    lan_id                   = ionoscloud_lan.example.id
    primary_instance_address = "192.168.2.101/24"
  }

  snapshot = {
    location       = "eu-central-3"
    retention_days = 7
    snapshot_hours = [0, 6, 12, 18]
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "cacheuser"
    password = {
      algorithm = "SHA-256"
      hash      = "5e884898da28047151d0e56f8dc6292773603d0d6aabbdd62a11ef721d1542d8"
    }
  }
}
```

## Argument Reference

* `location` - (Required)[string] The location of the cluster. Changing this forces a new resource. Available: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `name` - (Required)[string] Cluster name (2–63 alphanumeric chars with dashes, underscores, dots).
* `version` - (Required)[string] The InMemoryDB version. Upgrades only (see `/versions` endpoint).
* `eviction_policy` - (Required)[string] Key eviction strategy: `noeviction`, `allkeys-lru`, `allkeys-lfu`, `allkeys-random`, `volatile-lru`, `volatile-lfu`, `volatile-random`, `volatile-ttl`.
* `persistence_mode` - (Optional)(Computed)[string] Data persistence mode: `None`, `AOF`, `RDB`, `RDB_AOF`. If not set, the API default is used.
* `description` - (Optional)[string] Human-readable description.
* `logs_enabled` - (Optional)(Computed)[bool] Enable log collection for observability. If not set, the API default is used.
* `metrics_enabled` - (Optional)(Computed)[bool] Enable metrics collection for observability. If not set, the API default is used.
* `instances` - (Required)[object] The instance sizing configuration.
  * `count` - (Required)[int] Number of instances (1–5).
  * `cores` - (Required)[int] CPU cores per instance (1–62).
  * `ram` - (Required)[int] RAM per instance in GB (4–240). The API does not support RAM downgrade. Storage is automatically derived from RAM and persistence mode.
* `connections` - (Required)[object] The network connection for your cluster. Only one connection is allowed.
  * `datacenter_id` - (Required)[string] The Virtual Data Center ID to connect to.
  * `lan_id` - (Required)[string] The numeric LAN ID within the data center.
  * `primary_instance_address` - (Required)[string] Primary instance IP in CIDR notation.
* `snapshot` - (Required)[object] Snapshot configuration.
  * `location` - (Required)[string] Object Storage location for snapshots. Changing this forces the re-creation of the cluster.
  * `retention_days` - (Required)[int] Days to retain snapshots (1–365).
  * `snapshot_hours` - (Required)[list of int] UTC hours for scheduled snapshots (0–23). At least one hour must be specified.
* `maintenance_window` - (Required)[object] A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start time in UTC (`HH:MM:SS`).
  * `day_of_the_week` - (Required)[string] Day of the week: `Sunday`–`Saturday`.
* `credentials` - (Required)[object] Credentials for the InMemoryDB cluster user.
  * `username` - (Required)[string] Username (2–16 alphanumeric + underscore).
  * `password` - (Required)[object] Pre-hashed password. Not returned by the API — will be null in state after `terraform import`.
    * `algorithm` - (Required)[string] Hash algorithm (`SHA-256`).
    * `hash` - (Required)[string] **Sensitive.** Hex-encoded SHA-256 hash (64 lowercase hex chars).
* `restore_from_snapshot` - (Optional)[object] Restore configuration from a snapshot.
  * `source_snapshot_id` - (Optional)[string] UUID of the snapshot to restore from. Must be provided when the block is used during cluster creation. Not applicable for in-place restore via update.
  * `recovery_target_datetime` - (Optional)[string] ISO 8601 timestamp to restore from the most recent snapshot at or before that time. Optional for create-time restore; must be provided for in-place restore via update.

> **Note:** `restore_from_snapshot` is not returned by the API. The values are stored in state as configured but will be null after `terraform import`.

* `timeouts` - (Optional) Standard Terraform timeouts: `create`, `update`, `delete`.

## Attributes Reference

The following attributes are exported:

* `id` - The UUID of the cluster.
* `dns_name` - The DNS name for connecting to the cluster's primary instance.

## Import

InMemoryDB V2 clusters can be imported using `<location>:<cluster_id>`:

```shell
terraform import ionoscloud_inmemorydb_cluster_v2.example de/txl:example-id
```

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can be used with the `identity` attribute:

```hcl
import {
  to = ionoscloud_inmemorydb_cluster_v2.example
  identity = {
    id       = "example-id"
    location = "de/txl"
  }
}

resource "ionoscloud_inmemorydb_cluster_v2" "example" {
  ### Configuration omitted for brevity ###
}
```
