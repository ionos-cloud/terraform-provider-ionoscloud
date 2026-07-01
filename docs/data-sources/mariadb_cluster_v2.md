---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_cluster_v2"
sidebar_current: "docs-datasource-mariadb_cluster_v2"
description: |-
  Reads an IONOS CLOUD MariaDB V2 Cluster by ID or name.
---

# ionoscloud_mariadb_cluster_v2

The `ionoscloud_mariadb_cluster_v2` data source can be used to retrieve information about an existing MariaDB V2 cluster.

## Example Usage

### By id
```hcl
data "ionoscloud_mariadb_cluster_v2" "by_id" {
  id       = "example-id"
  location = "de/txl"
}
```

### By name
```hcl
data "ionoscloud_mariadb_cluster_v2" "by_name" {
  name     = "my-mariadb-cluster"
  location = "de/txl"
}
```

## Argument Reference

* `id` - (Optional) The UUID of the cluster.
* `name` - (Optional) The cluster name (exact match).
* `location` - (Required)[string] The location to query. Available locations: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

> **Note:** Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The UUID of the cluster.
* `name` - The name of the cluster.
* `location` - The location of the cluster.
* `description` - Human-readable description for the cluster.
* `version` - The MariaDB version for the cluster.
* `logs_enabled` - Whether log collection and reporting is enabled for this cluster's observability.
* `metrics_enabled` - Whether metrics collection and reporting is enabled for this cluster's observability.
* `dns_name` - The DNS name used to access the cluster.
* `instances` - Compute and storage configuration:
  * `count` - The total number of instances in the cluster.
  * `cores` - The number of CPU cores per instance.
  * `ram` - The amount of memory per instance in gigabytes (GB).
  * `storage_size` - The amount of storage per instance in gigabytes (GB).
* `connections` - Connection information of the MariaDB cluster:
  * `datacenter_id` - The ID of the Virtual Data Center the cluster is connected to.
  * `lan_id` - The numeric LAN ID the cluster is connected to.
  * `primary_instance_address` - The IP address and netmask of the cluster's primary instance, in CIDR notation.
* `backup` - Backup location and retention configuration:
  * `location` - The Object Storage location where the backup is stored.
  * `retention_days` - The number of days cluster backups are retained.
* `maintenance_window` - Maintenance window configuration:
  * `time` - Start of the maintenance window in UTC time.
  * `day_of_the_week` - The name of the week day.
* `credentials` - Credentials block:
  * `username` - The username of the initial MariaDB user.
  * `database` - The name of the initial database.

> **Note:** `credentials.password` is not available — the API never returns the password.
