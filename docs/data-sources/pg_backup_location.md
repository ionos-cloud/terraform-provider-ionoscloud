---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD : ionoscloud_pg_backup_location_v2"
sidebar_current: "docs-ionoscloud_pg_backup_location_v2"
description: |-
  Get information on DBaaS PostgreSQL v2 Backup Locations
---

# ionoscloud_pg_backup_location_v2

The **DBaaS PostgreSQL v2 Backup Locations data source** can be used to retrieve the list of available backup locations.
Use this data source to find valid values for the `backup_location` attribute of the `ionoscloud_pg_cluster_v2` resource.

## Example Usage

```hcl
data "ionoscloud_pg_backup_location_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The region in which to look up backup locations. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

## Attributes Reference

The following attributes are returned by the datasource:

* `backup_locations` - The list of available backup locations. Each backup location has the following attributes:
  * `id` - The ID (UUID) of the backup location.
  * `location` - The S3 location identifier (e.g. `de`, `eu-central-2`).
