---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_backups"
sidebar_current: "docs-ionoscloud_pg_backups"
description: |-
  Get information on DbaaS PgSql Backups
---

# ionoscloud\_pg_backups

The **DbaaS Postgres Backups data source** can be used to search for and return existing DbaaS Postgres Backups for a specific Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_pg_backups" "example" {
	cluster_id = <cluster_id>
}
```

## Argument Reference

* `cluster_id` - (Required) The unique ID of the cluster.

`cluster_id` must be provided. If it is not provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_id` - Id of the cluster.
* `cluster_backups` - List of backups.
    * `id` - The unique ID of the resource.
    * `cluster_id` - The unique ID of the cluster
    * `size` - The size of all base backups including the wal size in MB.
    * `location` - The S3 location where the backups will be stored.
    * `version` - The PostgreSQL version this backup was created from.
    * `is_active` - Whether a cluster currently backs up data to this backup.
    * `earliest_recovery_target_time` - The oldest available timestamp to which you can restore.
    * `type`
    * `metadata` - Metadata of the resource.
        * `created_date` - The ISO 8601 creation timestamp.
        * `created_by`
        * `created_by_user_id`
        * `last_modified_date` - The ISO 8601 modified timestamp.
        * `last_modified_by`
        * `last_modified_by_user_id`

**NOTE:** If the `earliestRecoveryTargetTime` is missing in your backup, the cluster from where you want to restore wasn't able to do a base backup. This can happen, when you e.g. quickly delete a newly created cluster, since the base backup will be triggered up to a minute after the cluster is available.