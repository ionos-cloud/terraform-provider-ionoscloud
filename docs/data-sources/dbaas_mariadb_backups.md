---
subcategory: "Database as a Service - MariaDB"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_mariadb_backups"
sidebar_current: "docs-ionoscloud_mariadb_backups"
description: |-
  Get information on DBaaS MariaDB Backups
---

# ionoscloud\_mariadb_backups

The **DBaaS MariaDB Backups data source** can be used to search for and return existing DBaaS MariaDB Backups for a specific cluster.

## Example Usage

### Get all backups for a specific cluster
```hcl
data "ionoscloud_mariadb_backups" "example" {
	cluster_id = <cluster_id>
}
```

### Get a specific backup
```hcl
data "ionoscloud_mariadb_backups" "example" {
	backup_id = <backup_id>
}
```

## Argument Reference

* `cluster_id` - (Optional)[string] The unique ID of the cluster.
* `backup_id` - (Optional)[string] The unique ID of the backup.

⚠️ **Note:** Either `cluster_id` or `backup_id` must be used, but not both at the same time.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_backups` - List of backups.
    * `id` - The unique ID of the backup.
    * `size` - The size of the backup in Mebibytes (MiB). This is the size of the binary backup file that was stored.
    * `created` - The ISO 8601 creation timestamp