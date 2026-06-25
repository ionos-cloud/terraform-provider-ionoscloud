---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_cluster_v2"
sidebar_current: "docs-datasource-inmemorydb_cluster_v2"
description: |-
  Reads an IONOS Cloud In-Memory DB v2 Cluster by ID or name.
---

# ionoscloud_inmemorydb_cluster_v2 (Data Source)

Reads an IONOS Cloud In-Memory DB v2 Cluster.

## Example Usage

```hcl
# Look up by ID
data "ionoscloud_inmemorydb_cluster_v2" "by_id" {
  id       = "e69b22a5-8fee-56b1-b6fb-4a07e4205ead"
  location = "de/txl"
}

# Look up by name
data "ionoscloud_inmemorydb_cluster_v2" "by_name" {
  name     = "my-inmemorydb-cluster"
  location = "de/txl"
}
```

## Argument Reference

Exactly one of `id` or `name` must be provided.

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | One of | The UUID of the cluster. |
| `name` | String | One of | The cluster name (exact match). |
| `location` | String | Yes | The location of the cluster. |

## Attributes Reference

| Attribute | Description |
|-----------|-------------|
| `id` | The UUID of the cluster. |
| `name` | The name of the cluster. |
| `location` | The location of the cluster. |
| `description` | Human-readable description for the cluster. |
| `version` | The In-Memory DB version. |
| `persistence_mode` | The data persistence mode. |
| `eviction_policy` | The key eviction strategy. |
| `logs_enabled` | Whether log collection is enabled. |
| `metrics_enabled` | Whether metrics collection is enabled. |
| `dns_name` | The DNS name for connecting to the cluster's primary instance. |
| `instances.count` | Number of instances. |
| `instances.cores` | CPU cores per instance. |
| `instances.ram` | RAM per instance in GB. |
| `connections.datacenter_id` | The Virtual Data Center ID. |
| `connections.lan_id` | The numeric LAN ID. |
| `connections.primary_instance_address` | The primary instance IP in CIDR notation. |
| `snapshot.location` | Object Storage location for snapshots. |
| `snapshot.retention_days` | Days snapshots are retained. |
| `snapshot.snapshot_hours` | UTC hours at which snapshots are taken. |
| `maintenance_window.time` | Maintenance window start time in UTC (HH:MM:SS). |
| `maintenance_window.day_of_the_week` | Maintenance window day of the week. |
| `credentials.username` | The username for the In-Memory DB user. |
