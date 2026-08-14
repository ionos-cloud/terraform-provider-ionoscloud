---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: inmemorydb_cluster_v2"
description: |-
  Lists IONOS CLOUD InMemoryDB V2 Clusters.
---

# List Resource: ionoscloud_inmemorydb_cluster_v2

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [InMemoryDB V2 Clusters](https://docs.ionos.com/cloud/databases/in-memory-db) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `.tfquery.hcl` file, separate from your main Terraform configuration.

### List all clusters

```hcl
list "ionoscloud_inmemorydb_cluster_v2" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter clusters by location

```hcl
list "ionoscloud_inmemorydb_cluster_v2" "de_txl" {
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
list "ionoscloud_inmemorydb_cluster_v2" "prod" {
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

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_inmemorydb_cluster_v2` resource blocks for all existing clusters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_inmemorydb_cluster_v2` resource block for each discovered cluster into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

Supported `location` values: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

> **Performance note:** When no `location` filter is set, the provider queries every regional endpoint in sequence. Adding a `location` filter reduces the query to a single endpoint call.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                          |
|------------|----------------------------------------------------------------------|
| `id`       | The UUID of the cluster.                                             |
| `location` | The regional endpoint the cluster was fetched from (e.g. `de/txl`). |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_inmemorydb_cluster_v2` resource schema:

- `id` - The UUID of the cluster.
- `name` - The name of the cluster.
- `description` - A human-readable description for the cluster.
- `version` - The InMemoryDB version (e.g. `9.0`).
- `dns_name` - The DNS name used to connect to the cluster's primary instance.
- `location` - The regional location of the cluster (e.g. `de/txl`).
- `persistence_mode` - The data persistence mode (`None`, `AOF`, `RDB`, `RDB_AOF`).
- `eviction_policy` - The key eviction strategy.
- `logs_enabled` - Whether log collection is enabled.
- `metrics_enabled` - Whether metrics collection is enabled.
- `instances` - Instance sizing block:
  - `count` - Number of instances.
  - `cores` - CPU cores per instance.
  - `ram` - RAM per instance in GB.
- `connections` - Network connection block:
  - `datacenter_id` - UUID of the connected datacenter.
  - `lan_id` - ID of the connected LAN.
  - `primary_instance_address` - IP address of the primary instance in CIDR notation.
- `maintenance_window` - Maintenance window block:
  - `time` - Start time in UTC (HH:MM:SS).
  - `day_of_the_week` - Day of the week (e.g. `Sunday`).
- `snapshot` - Snapshot configuration block:
  - `location` - Object Storage location for snapshots.
  - `retention_days` - Number of days snapshots are retained.
  - `snapshot_hours` - UTC hours at which snapshots are taken.
- `credentials` - Credentials block:
  - `username` - The username for the InMemoryDB user.

> **Note:** `credentials.password` is not available in the list resource — the API never returns the password hash. Only `credentials.username` is populated.
