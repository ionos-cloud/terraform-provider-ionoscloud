---
subcategory: "Database as a Service - PostgreSQL"
layout: "ionoscloud"
page_title: "IONOS CLOUD: pg_cluster_v2"
description: |-
  Lists IONOS CLOUD PostgreSQL v2 Clusters.
---

# List Resource: ionoscloud_pg_cluster_v2

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [IONOS PostgreSQL v2 Clusters](https://docs.ionos.com/cloud/databases/postgresql) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `tfquery.hcl` file, separate from your main Terraform configuration.

### List all clusters

```hcl
list "ionoscloud_pg_cluster_v2" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter clusters by location

```hcl
list "ionoscloud_pg_cluster_v2" "de_txl" {
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
list "ionoscloud_pg_cluster_v2" "prod" {
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

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_pg_cluster_v2` resource blocks for all existing clusters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_pg_cluster_v2` resource block for each discovered cluster into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

Supported `location` values: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

> **Performance note:** When no `location` filter is set, the provider queries every regional endpoint in sequence. Adding a `location` filter reduces the query to a single endpoint call.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                         |
|------------|---------------------------------------------------------------------|
| `id`       | The UUID of the cluster.                                            |
| `location` | The regional endpoint the cluster was fetched from (e.g. `de/txl`). |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_pg_cluster_v2` resource schema:

- `id` - The UUID of the cluster.
- `name` - The name of the cluster.
- `description` - A short description of the cluster.
- `version` - The PostgreSQL version (e.g. `16`).
- `dns_name` - The DNS name of the cluster endpoint.
- `location` - The regional location of the cluster (e.g. `de/txl`).
- `replication_mode` - The replication mode (`SYNCHRONOUS`, etc.).
- `connection_pooler` - The connection pooler mode (`NONE`, `PGBOUNCER`).
- `logs_enabled` - Whether log forwarding is enabled.
- `metrics_enabled` - Whether metrics export is enabled.
- `instances` - Instance sizing block:
  - `count` - Number of replicas.
  - `cores` - vCPU count per instance.
  - `ram` - RAM in MB per instance.
  - `storage_size` - Storage in MB per instance.
- `connections` - Network connection block:
  - `datacenter_id` - UUID of the connected datacenter.
  - `lan_id` - ID of the connected LAN.
  - `primary_instance_address` - IP address of the primary instance.
- `maintenance_window` - Maintenance window block:
  - `time` - Time of the maintenance window (HH:MM:SS).
  - `day_of_the_week` - Day of the week (e.g. `Sunday`).
- `backup` - Backup configuration block:
  - `location` - Location where backups are stored.
  - `retention_days` - Number of days backups are retained.
- `credentials` - Database credentials block:
  - `username` - The database username.
  - `database` - The default database name.
  - `password` - Always null (write-only, not returned by the API).
  - `password_version` - Always null (local state only, not available during listing).
- `restore_from_backup` - Set when the cluster was created from a backup:
  - `source_backup_id` - UUID of the source backup.
  - `recovery_target_datetime` - RFC3339 point-in-time recovery target.

> **Note:** `credentials.password` is write-only and never returned by the API (it will be null). `credentials.password_version` is not available via the list resource as it is only tracked in local Terraform state.
