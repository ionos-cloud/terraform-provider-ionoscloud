---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_cluster_v2"
sidebar_current: "docs-datasource-mongo_cluster_v2"
description: |-
  Reads an IONOS CLOUD MongoDB V2 Cluster by ID or name.
---

# ionoscloud_mongo_cluster_v2

The `ionoscloud_mongo_cluster_v2` data source can be used to retrieve information about an existing MongoDB V2 cluster.

## Example Usage

### By id
```hcl
data "ionoscloud_mongo_cluster_v2" "by_id" {
  id       = "example-id"
  location = "gb/bhx"
}
```

### By name
```hcl
data "ionoscloud_mongo_cluster_v2" "by_name" {
  name     = "my-mongodb-cluster"
  location = "gb/bhx"
}
```

## Argument Reference

* `id` - (Optional)[string] The UUID of the cluster.
* `name` - (Optional)[string] The cluster name (exact match).
* `location` - (Required)[string] The location of the cluster. Available locations: `gb/bhx`, `us/ewr`.

> **Note:** Exactly one of `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the cluster.
* `name` - The name of the cluster.
* `location` - The location of the cluster.
* `edition` - The cluster edition (`ENTERPRISE` or `BUSINESS`).
* `version` - The MongoDB version for the cluster.
* `description` - Human-readable description for the cluster.
* `connection_uri` - The connection string used to access the cluster.
* `bi_connector_enabled` - Whether the MongoDB BI Connector is enabled. Only set for the `ENTERPRISE` edition.
* `logs_enabled` - Whether log collection and reporting is enabled for this cluster's observability.
* `metrics_enabled` - Whether metrics collection and reporting is enabled for this cluster's observability.
* `instance_count` - The number of instances in the cluster. Only set for the `BUSINESS` edition.
* `instance_template_id` - The ID of the instance template. Only set for the `BUSINESS` edition.
* `instances` - Compute and storage configuration for each instance. Only set for the `ENTERPRISE` edition:
  * `count` - The total number of instances in the cluster (one primary and n-1 secondary).
  * `cores` - The number of CPU cores per instance.
  * `ram` - The amount of memory per instance in gigabytes (GB).
  * `storage_size` - The amount of storage per instance in gigabytes (GB).
* `sharding` - Sharding configuration. Only set for the `ENTERPRISE` edition:
  * `enabled` - Whether sharding is enabled.
  * `shards` - The number of shards in the cluster.
* `connections` - Connection information of the MongoDB cluster:
  * `datacenter_id` - The ID of the Virtual Data Center the cluster is connected to.
  * `lan_id` - The numeric LAN ID the cluster is connected to.
  * `instance_addresses` - The IP addresses (with netmask, in CIDR notation) assigned to the cluster instances.
* `snapshot` - Snapshot location and retention configuration:
  * `location` - The Object Storage location where the snapshot is stored.
  * `retention_days` - The number of days cluster snapshots are retained.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur:
  * `time` - Start of the maintenance window in UTC time.
  * `day_of_the_week` - The name of the week day.
* `credentials` - Credentials for the initial database user:
  * `username` - The username of the initial MongoDB user.
  * `database` - The name of the initial database.

> **Note:** `credentials.password` is never returned by the API.
