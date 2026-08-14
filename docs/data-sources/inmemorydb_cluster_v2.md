---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_cluster_v2"
sidebar_current: "docs-datasource-inmemorydb_cluster_v2"
description: |-
  Reads an IONOS CLOUD InMemoryDB V2 Cluster by ID or name.
---

# ionoscloud_inmemorydb_cluster_v2

The `ionoscloud_inmemorydb_cluster_v2` data source can be used to retrieve information about an existing InMemoryDB V2 cluster.

## Example Usage

### By id
```hcl
data "ionoscloud_inmemorydb_cluster_v2" "by_id" {
  id       = "example-id"
  location = "de/txl"
}
```

### By name
```hcl
data "ionoscloud_inmemorydb_cluster_v2" "by_name" {
  name     = "my-inmemorydb-cluster"
  location = "de/txl"
}
```

## Argument Reference

* `id` - (Optional) The UUID of the cluster.
* `name` - (Optional) The cluster name (exact match).
* `location` - (Required)[string] The location of the cluster.

> **Note:** Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the cluster.
* `name` - The name of the cluster.
* `location` - The location of the cluster.
* `description` - Human-readable description for the cluster.
* `version` - The InMemoryDB version.
* `persistence_mode` - The data persistence mode.
* `eviction_policy` - The key eviction strategy.
* `logs_enabled` - Whether log collection is enabled.
* `metrics_enabled` - Whether metrics collection is enabled.
* `dns_name` - The DNS name for connecting to the cluster's primary instance.
* `instances` - Instance sizing configuration:
  * `count` - Number of instances.
  * `cores` - CPU cores per instance.
  * `ram` - RAM per instance in GB.
* `connections` - Network connection configuration:
  * `datacenter_id` - The Virtual Data Center ID.
  * `lan_id` - The numeric LAN ID.
  * `primary_instance_address` - The primary instance IP in CIDR notation.
* `snapshot` - Snapshot configuration:
  * `location` - Object Storage location for snapshots.
  * `retention_days` - Days snapshots are retained.
  * `snapshot_hours` - UTC hours at which snapshots are taken.
* `maintenance_window` - Maintenance window configuration:
  * `time` - Maintenance window start time in UTC (HH:MM:SS).
  * `day_of_the_week` - Maintenance window day of the week.
* `credentials` - Credentials block:
  * `username` - The username for the InMemoryDB user.
