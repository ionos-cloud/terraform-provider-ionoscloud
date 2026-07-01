---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_backup_locations_v2"
sidebar_current: "docs-datasource-mariadb_backup_locations_v2"
description: |-
  Lists IONOS CLOUD MariaDB V2 Backup Locations.
---

# ionoscloud_mariadb_backup_locations_v2

The `ionoscloud_mariadb_backup_locations_v2` data source can be used to retrieve the supported MariaDB V2 backup locations for a given API endpoint.

## Example Usage

```hcl
data "ionoscloud_mariadb_backup_locations_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Requests are routed to the corresponding regional MariaDB endpoint. Available locations: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of backup location objects. Each item includes:
  * `id` - The ID (UUID) of the backup location.
  * `location` - The Object Storage location name.
