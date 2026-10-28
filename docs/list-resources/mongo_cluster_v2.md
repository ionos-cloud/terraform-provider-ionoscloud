---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: mongo_cluster_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Clusters.
---

# List Resource: ionoscloud_mongo_cluster_v2

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [MongoDB V2 Clusters](https://docs.ionos.com/cloud/databases/mongodb) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `.tfquery.hcl` file, separate from your main Terraform configuration.

### List all clusters

```hcl
list "ionoscloud_mongo_cluster_v2" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

### Filter clusters by location

```hcl
list "ionoscloud_mongo_cluster_v2" "gb_bhx" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "location"
      field_value = "gb/bhx"
    }]
  }
}
```

### Filter clusters by name and location

```hcl
list "ionoscloud_mongo_cluster_v2" "prod" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [
      { field_name = "name",     field_value = "my-cluster" },
      { field_name = "location", field_value = "gb/bhx" },
    ]
  }
}
```

### Generate resource configuration from existing clusters

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_mongo_cluster_v2` resource blocks for all existing clusters:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_mongo_cluster_v2` resource block for each discovered cluster into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `location`.
  - `field_value` - (Required) The exact value to match against.

Supported `location` values: `gb/bhx`, `us/ewr`.

> **Performance note:** When no `location` filter is set, the provider queries every regional endpoint in sequence. Adding a `location` filter reduces the query to a single endpoint call.

## Identity Attributes

Each result exposes the following identity attributes, usable for import:

| Attribute  | Description                                                          |
|------------|----------------------------------------------------------------------|
| `id`       | The UUID of the cluster.                                             |
| `location` | The regional endpoint the cluster was fetched from (e.g. `gb/bhx`). |

## Attributes Reference

Each result exposes the following attributes when `include_resource = true`, matching the `ionoscloud_mongo_cluster_v2` resource schema:

- `id` - The UUID of the cluster.
- `name` - The name of the cluster.
- `location` - The regional location of the cluster (e.g. `gb/bhx`).
- `edition` - The cluster edition (`ENTERPRISE` or `BUSINESS`).
- `version` - The MongoDB version for the cluster.
- `description` - A human-readable description for the cluster.
- `connection_uri` - The connection string used to access the cluster.
- `bi_connector_enabled` - Whether the MongoDB BI Connector is enabled. Only set for the `ENTERPRISE` edition.
- `logs_enabled` - Whether log collection and reporting is enabled.
- `metrics_enabled` - Whether metrics collection and reporting is enabled.
- `instance_count` - The number of instances in the cluster. Only set for the `BUSINESS` edition.
- `instance_template_id` - The ID of the instance template. Only set for the `BUSINESS` edition.
- `instances` - Instance sizing block. Only set for the `ENTERPRISE` edition:
  - `count` - The total number of instances in the cluster (one primary and n-1 secondary).
  - `cores` - The number of CPU cores per instance.
  - `ram` - The amount of memory per instance in GB.
  - `storage_size` - The amount of storage per instance in GB.
- `sharding` - Sharding block. Only set for the `ENTERPRISE` edition:
  - `enabled` - Whether sharding is enabled.
  - `shards` - The number of shards in the cluster.
- `connections` - Network connection block:
  - `datacenter_id` - UUID of the connected datacenter.
  - `lan_id` - ID of the connected LAN.
  - `instance_addresses` - IP addresses assigned to the cluster instances, in CIDR notation.
- `maintenance_window` - Maintenance window block:
  - `time` - Start time in UTC (HH:MM:SS).
  - `day_of_the_week` - Day of the week (e.g. `Sunday`).
- `snapshot` - Snapshot configuration block:
  - `location` - Object Storage location for snapshots.
  - `retention_days` - Number of days snapshots are retained.
- `credentials` - Credentials block:
  - `username` - The username of the initial MongoDB user.
  - `database` - The name of the initial database.

> **Note:** `credentials.password` is not available in the list resource — the API never returns the password. Only `credentials.username` and `credentials.database` are populated.
