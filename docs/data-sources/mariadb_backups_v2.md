---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_backups_v2"
sidebar_current: "docs-datasource-mariadb_backups_v2"
description: |-
  Lists IONOS CLOUD MariaDB V2 Backups.
---

# ionoscloud_mariadb_backups_v2

The `ionoscloud_mariadb_backups_v2` data source can be used to retrieve information about existing MariaDB V2 backups, with an optional cluster ID filter.

## Example Usage

```hcl
data "ionoscloud_mariadb_backups_v2" "all" {
  location = "de/txl"
}

data "ionoscloud_mariadb_backups_v2" "for_cluster" {
  location   = "de/txl"
  cluster_id = "example-id"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `cluster_id` - (Optional)[string] Filter backups by the cluster they belong to.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of backups. Each item includes:
  * `id` - The ID (UUID) of the backup.
  * `location` - The Object Storage location where the backup is stored.
  * `cluster_id` - The ID of the cluster this backup belongs to.
  * `cluster_name` - The name of the cluster this backup belongs to.
  * `mariadb_cluster_version` - The MariaDB version of the cluster at backup time.
  * `earliest_recovery_target_time` - The earliest point in time to which the cluster can be restored from this backup (RFC3339).
  * `latest_recovery_target_time` - The latest point in time to which the cluster can be restored (RFC3339). Empty if the backup can be restored up to the current time.
