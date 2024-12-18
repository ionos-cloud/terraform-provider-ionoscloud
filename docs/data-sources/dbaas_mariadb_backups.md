---
subcategory: "Database as a Service - MariaDB"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_mariadb_backups"
sidebar_current: "docs-ionoscloud_mariadb_backups"
description: |-
  Get information on DBaaS MariaDB Backups
---

# ionoscloud_mariadb_backups

The **DBaaS MariaDB Backups data source** can be used to search for and return existing DBaaS MariaDB Backups for a specific cluster.

## Example Usage

### Get all backups for a specific cluster
```hcl
data "ionoscloud_mariadb_backups" "example" { 
    cluster_id = "cluster_id"
    location   = "de/txl"
}
```

### Get a specific backup
```hcl
data "ionoscloud_mariadb_backups" "example" {
    backup_id = "backup_id"
    location   = "de/txl"
}
```

## Argument Reference

* `cluster_id` - (Optional)[string] The unique ID of the cluster.
* `backup_id` - (Optional)[string] The unique ID of the backup.
* `location`- (Optional)[string] The location of the cluster. Different service endpoints are used based on location, possible options are: "de/fra", "de/txl", "es/vit", "fr/par", "gb/lhr", "us/ewr", "us/las", "us/mci". If not set, the endpoint will be the one corresponding to "de/txl".

⚠️ **Note:** Either `cluster_id` or `backup_id` must be used, but not both at the same time.

> **⚠ WARNING:** `Location` attribute will become required in the future.


## Attributes Reference

The following attributes are returned by the datasource:

* `bakups` - List of backups.
    * `cluster_id` - The unique ID of the cluster that was backed up.
    * `earliest_recovery_target_time` - The oldest available timestamp to which you can restore.
    * `size` - Size of all base backups in Mebibytes (MiB). This is at least the sum of all base backup sizes.
    * `base_backups` - The list of backups for the specified cluster
      * `size` - The size of the backup in Mebibytes (MiB). This is the size of the binary backup file that was stored
      * `created` - The ISO 8601 creation timestamp