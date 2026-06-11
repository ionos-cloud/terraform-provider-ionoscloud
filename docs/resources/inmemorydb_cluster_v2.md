---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_cluster_v2"
sidebar_current: "docs-resource-inmemorydb_cluster_v2"
description: |-
  Creates and manages an IONOS Cloud In-Memory DB v2 Cluster.
---

# ionoscloud_inmemorydb_cluster_v2

Manages an IONOS Cloud In-Memory DB v2 Cluster.

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
  description      = "Production In-Memory DB cluster"
  version          = "9.0"
  persistence_mode = "RDB"
  eviction_policy  = "allkeys-lru"
  logs_enabled     = true
  metrics_enabled  = true

  instances = {
    count = 1
    cores = 2
    ram   = 8
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

The following arguments are supported:

### Top-level

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | The location of the cluster. Changing this forces a new resource. Available: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`. |
| `name` | String | Yes | Cluster name (2–63 alphanumeric chars with dashes, underscores, dots). |
| `version` | String | Yes | The In-Memory DB version. Upgrades only (see `/versions` endpoint). |
| `persistence_mode` | String | Yes | Data persistence mode: `None`, `AOF`, `RDB`, `RDB_AOF`. |
| `eviction_policy` | String | Yes | Key eviction strategy: `noeviction`, `allkeys-lru`, `allkeys-lfu`, `allkeys-random`, `volatile-lru`, `volatile-lfu`, `volatile-random`, `volatile-ttl`. |
| `description` | String | No | Human-readable description. |
| `logs_enabled` | Boolean | No | Enable log collection for observability. |
| `metrics_enabled` | Boolean | No | Enable metrics collection for observability. |

### `instances` block (required)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `count` | Number | Yes | Number of instances (1–5). |
| `cores` | Number | Yes | CPU cores per instance (1–62). |
| `ram` | Number | Yes | RAM per instance in GB (4–240). RAM cannot be downgraded. Storage is automatically derived from RAM and persistence mode. |

### `connections` block (required)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `datacenter_id` | String | Yes | The Virtual Data Center ID to connect to. |
| `lan_id` | String | Yes | The numeric LAN ID within the data center. |
| `primary_instance_address` | String | Yes | Primary instance IP in CIDR notation. |

### `snapshot` block (required)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | Object Storage location for snapshots. |
| `retention_days` | Number | No | Days to retain snapshots (1–365, default 7). |
| `snapshot_hours` | List(Number) | Yes | UTC hours for scheduled snapshots (0–23). At least one hour must be specified. |

### `maintenance_window` block (required)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `time` | String | Yes | Start time in UTC (`HH:MM:SS`). |
| `day_of_the_week` | String | Yes | Day of the week: `Sunday`–`Saturday`. |

### `credentials` block (required)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `username` | String | Yes | Username (2–16 alphanumeric + underscore). |
| `password` | Block | Yes | Pre-hashed password. |

#### `credentials.password` block

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `algorithm` | String | No | Hash algorithm (`SHA-256`). |
| `hash` | String | Yes | **Sensitive.** Hex-encoded SHA-256 hash (64 lowercase hex chars). |

### `restore_from_snapshot` block (optional, write-only)

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `source_snapshot_id` | String | No | UUID of the snapshot to restore from. Required for create-time restore; not used for in-place restore via update. |
| `recovery_target_datetime` | String | No | ISO 8601 timestamp to restore from the most recent snapshot at or before that time. Optional for create-time restore; required for in-place restore via update. |

> **Note:** This block is write-only and is not returned by the API after apply.

### `timeouts` block

Standard Terraform timeouts: `create`, `read`, `update`, `delete`.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

| Attribute | Description |
|-----------|-------------|
| `id` | The UUID of the cluster. |
| `dns_name` | The DNS name for connecting to the cluster's primary instance. |

## Import

InMemoryDB v2 clusters can be imported using `<location>:<cluster_id>`:

```bash
terraform import ionoscloud_inmemorydb_cluster_v2.example de/txl:e69b22a5-8fee-56b1-b6fb-4a07e4205ead
```
