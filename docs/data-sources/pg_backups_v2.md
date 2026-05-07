---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD : ionoscloud_pg_backups_v2"
sidebar_current: "docs-ionoscloud_pg_backups_v2"
description: |-
  Get information on DBaaS PostgreSQL v2 Backups
---

# ionoscloud_pg_backups_v2

The **DBaaS PostgreSQL v2 Backups data source** can be used to list existing DBaaS PostgreSQL v2 Backups.
An optional cluster ID filter can be used to narrow down results to backups of a specific cluster.

## Example Usage

### List all backups in a location
```hcl
data "ionoscloud_pg_backups_v2" "example" {
  location = "de/txl"
}
```

### List backups for a specific cluster
```hcl
data "ionoscloud_pg_backups_v2" "example" {
  location   = "de/txl"
  cluster_id = "cluster_id"
}
```

## Argument Reference

* `location` - (Required)[string] The region in which to look up backups. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `cluster_id` - (Optional)[string] The ID (UUID) of the cluster to filter backups by.

## Attributes Reference

The following attributes are returned by the datasource:

* `backups` - The list of backups. Each backup has the following attributes:
  * `id` - The ID (UUID) of the backup.
  * `cluster_id` - The ID (UUID) of the cluster the backup belongs to.
  * `postgres_cluster_version` - The PostgreSQL version of the cluster when the backup was created.
  * `is_active` - Whether the backup is active.
  * `earliest_recovery_target_time` - The earliest point in time to which the cluster can be restored.
  * `latest_recovery_target_time` - The latest point in time to which the cluster can be restored. If the backup can be restored up to the current time, this field will be null.
  * `location` - The S3 location where the backup is stored.