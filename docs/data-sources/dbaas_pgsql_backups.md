---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_backups"
sidebar_current: "docs-ionoscloud_pg_backups"
description: |-
Get information on DbaaS PgSql Backups
---

# ionoscloud\_pg_backups

The DbaaS Postgres Backups data source can be used to search for and return existing DbaaS Postgres Backups for a specific Cluster.

## Example Usage

```hcl
data "ionoscloud_pg_backups" "test_ds_dbaas_backups" {
	cluster_id = ionoscloud_pg_cluster.test_dbaas_cluster.id
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
    * `display_name` - The friendly name of your cluster.
    * `type`
    * `metadata` - Metadata of the resource.
        * `created_date` - The ISO 8601 creation timestamp.
        * `created_by`
        * `created_by_user_id`
        * `last_modified_date` - The ISO 8601 modified timestamp.
        * `last_modified_by`
        * `last_modified_by_user_id`