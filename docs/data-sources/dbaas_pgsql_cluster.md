---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_cluster"
sidebar_current: "docs-ionoscloud_pg_cluster"
description: |-
  Get information on a DbaaS PgSql Cluster
---

# ionoscloud\_pg_cluster

The **DbaaS Postgres Cluster data source** can be used to search for and return an existing DbaaS Postgres Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID 
```hcl
data "ionoscloud_pg_cluster" "example" {
  id	= <cluster_id>
}
```

### By Name

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

* `postgres_version` - The PostgreSQL version of your cluster.
* `instances` - The total number of instances in the cluster (one master and n-1 standbys)
* `cores` - The number of CPU cores per replica.
* `ram` - The amount of memory per instance in megabytes. 
* `storage_size` - The amount of storage per instance in MB.
* `storage_type` - The storage type used in your cluster. 
* `connections` - Details about the network connection for your cluster.
  * `datacenter_id` - The datacenter to connect your cluster to.
  * `lan_id` - The LAN to connect your cluster to.
  * `cidr` - The IP and subnet for the database. 
* `location` - The physical location where the cluster will be created. This will be where all of your instances live. 
* `display_name` - The friendly name of your cluster.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur
  * `time` 
  * `day_of_the_week` 
* `credentials` - Credentials for the database user to be created.
  * `username` - The username for the initial postgres user. 
* `synchronization_mode` - Represents different modes of replication. 
* `from_backup` - The unique ID of the backup you want to restore.
  * `backup_id` - The PostgreSQL version of your cluster.
  * `recovery_target_time` - If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. 
  