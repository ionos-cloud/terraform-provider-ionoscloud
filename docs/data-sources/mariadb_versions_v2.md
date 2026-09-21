---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_versions_v2"
sidebar_current: "docs-datasource-mariadb_versions_v2"
description: |-
  Lists supported IONOS CLOUD MariaDB V2 Versions.
---

# ionoscloud_mariadb_versions_v2

The `ionoscloud_mariadb_versions_v2` data source can be used to retrieve the supported MariaDB V2 versions in a given location.

## Example Usage

```hcl
data "ionoscloud_mariadb_versions_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `de/fra`, `de/fra/1`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of supported versions. Each item includes:
  * `id` - The ID (UUID) of the version.
  * `version` - The MariaDB version string (e.g. `11.4`).
  * `status` - The support status of the version.
  * `comment` - Additional human-readable information about the version lifecycle.
  * `can_upgrade_to` - List of versions that a cluster running this version can be upgraded to.
