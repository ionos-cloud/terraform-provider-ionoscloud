---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_versions_v2"
sidebar_current: "docs-ionoscloud_pg_versions_v2"
description: |-
  Get information on DBaaS PostgreSQL v2 Versions
---

# ionoscloud_pg_versions_v2

The **DBaaS PostgreSQL v2 Versions data source** can be used to retrieve the list of available PostgreSQL versions.

## Example Usage

```hcl
data "ionoscloud_pg_versions_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The region in which to look up available versions. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

## Attributes Reference

The following attributes are returned by the datasource:

* `versions` - The list of available PostgreSQL versions. Each version has the following attributes:
  * `id` - The ID (UUID) of the PostgreSQL version.
  * `version` - The PostgreSQL version string.
  * `status` - The support status of the version (e.g. `BETA`, `SUPPORTED`, `RECOMMENDED`, `DEPRECATED`).
  * `comment` - Additional information about the version status.
  * `can_upgrade_to` - List of versions that this version can be upgraded to.
