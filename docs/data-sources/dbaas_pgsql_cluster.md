---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_cluster"
sidebar_current: "docs-ionoscloud_pg_cluster"
description: |-
Get information on a DbaaS PgSql Cluster
---

# ionoscloud\_pg_cluster

The DbaaS Postgres Cluster data source can be used to search for and return an existing DbaaS Postgres Cluster.

## Example Usage

```hcl
data "ionoscloud_pg_cluster" "example" {
  display_name	= "PostgreSQL_cluster"
}
```

## Argument Reference

* `display_name` - (Optional) Display name or an existing cluster that you want to search for.
* `id` - (Optional) ID of the cluster you want to search for.

Either `display_name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the cluster
* `postgres_version` - The PostgreSQL version of your cluster.
* `replicas` - The number of replicas in your cluster.
* `cpu_core_count` - The number of CPU cores per replica.
* `ram_size` - The amount of memory per replica.
* `storage_size` -The amount of storage per replica.
* `storage_type` - The storage type used in your cluster.
* `vdc_connections` - The VDC to connect to your cluster.
    * `vdc_id` 
    * `lan_id` 
    * `ip_address` - The IP and subnet for the database.
* `location` - The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)
* `display_name` - The friendly name of your cluster.
* `backup_enabled` - Deprecated: backup is always enabled. Enables automatic backups of your cluster.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur
    * `time` 
    * `weekday`
* `credentials` - Credentials for the database user to be created.
    * `username` - The username for the initial postgres user.
    * `password` 
* `synchronization_mode` - Represents different modes of replication. Can have one of the following values: asynchronous, synchronous, strictly_synchronous
