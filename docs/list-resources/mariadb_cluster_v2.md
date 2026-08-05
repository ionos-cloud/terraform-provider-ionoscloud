---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: mariadb_cluster_v2"
description: |-
  Lists IONOS CLOUD MariaDB V2 Clusters.
---

# List Resource: ionoscloud_mariadb_cluster_v2

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [MariaDB V2 Clusters](https://docs.ionos.com/cloud/databases/mariadb) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `.tfquery.hcl` file, separate from your main Terraform configuration.

### List all clusters

```hcl
list "ionoscloud_mariadb_cluster_v2" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter clusters by location

```hcl
list "ionoscloud_mariadb_cluster_v2" "de_txl" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "location"
      field_value = "de/txl"
    }]
  }
}
```

### Filter clusters by name and location

```hcl
list "ionoscloud_mariadb_cluster_v2" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "my-cluster" },
      { field_name = "location", field_value = "de/txl" },
    ]
  }
}
```

### Generate resource configuration from existing clusters

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_mariadb_cluster_v2` resource blocks for all existing clusters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_mariadb_cluster_v2` resource block for each discovered cluster into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

Supported `location` values: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

> **Performance note:** When no `location` filter is set, the provider queries every regional endpoint in sequence. Adding a `location` filter reduces the query to a single endpoint call.

> **Note:** Unlike the [`ionoscloud_mariadb_clusters_v2`](../data-sources/mariadb_clusters_v2.md) data source, the `name` filter here is an **exact match**, not a partial match.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                          |
|------------|----------------------------------------------------------------------|
| `id`       | The UUID of the cluster.                                             |
| `location` | The regional endpoint the cluster was fetched from (e.g. `de/txl`). |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_mariadb_cluster_v2` resource schema:

- `id` - The UUID of the cluster.
- `name` - The name of the cluster.
- `description` - A human-readable description for the cluster.
- `version` - The MariaDB version (e.g. `11.4`).
- `dns_name` - The DNS name used to access the cluster.
- `location` - The regional location of the cluster (e.g. `de/txl`).
- `logs_enabled` - Whether log collection is enabled.
- `metrics_enabled` - Whether metrics collection is enabled.
- `instances` - Compute and storage configuration block:
  - `count` - The total number of instances in the cluster.
  - `cores` - The number of CPU cores per instance.
  - `ram` - The amount of memory per instance in gigabytes (GB).
  - `storage_size` - The amount of storage per instance in gigabytes (GB).
- `connections` - Connection information block:
  - `datacenter_id` - UUID of the connected datacenter.
  - `lan_id` - ID of the connected LAN.
  - `primary_instance_address` - IP address of the primary instance in CIDR notation.
- `backup` - Backup configuration block:
  - `location` - The Object Storage location where the backup will be created.
  - `retention_days` - Number of days cluster backups are retained.
- `maintenance_window` - Maintenance window block:
  - `time` - Start time in UTC (HH:MM:SS).
  - `day_of_the_week` - Day of the week (e.g. `Sunday`).
- `credentials` - Credentials block:
  - `username` - The username of the initial MariaDB user.
  - `database` - The name of the initial database.

> **Note:** `credentials.password` is not available in the list resource — the API never returns the password. Only `credentials.username` and `credentials.database` are populated.
