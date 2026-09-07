---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_versions_v2"
sidebar_current: "docs-datasource-inmemorydb_versions_v2"
description: |-
  Lists supported IONOS CLOUD InMemoryDB V2 Versions.
---

# ionoscloud_inmemorydb_versions_v2

The `ionoscloud_inmemorydb_versions_v2` data source can be used to retrieve the supported InMemoryDB V2 versions in a given location.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_versions_v2" "example" {
  location = "de/txl"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of versions. Each item includes:
  * `id` - The UUID of the version.
  * `version` - The version string (e.g. `9.0`).
  * `status` - The support status of the version (e.g. `SUPPORTED`, `RECOMMENDED`).
  * `comment` - Additional human-readable information about the version lifecycle.
  * `can_upgrade_to` - List of versions that a cluster running this version can be upgraded to.
