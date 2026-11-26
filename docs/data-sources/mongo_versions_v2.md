---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_versions_v2"
sidebar_current: "docs-datasource-mongo_versions_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Versions.
---

# ionoscloud_mongo_versions_v2

The `ionoscloud_mongo_versions_v2` data source can be used to retrieve the supported MongoDB V2 versions in a given location, along with their support status and valid upgrade targets.

## Example Usage

```hcl
data "ionoscloud_mongo_versions_v2" "example" {
  location = "gb/bhx"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `de/fra`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/mci`.

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of supported versions. Each item exposes:
  * `id` - The UUID of the version.
  * `version` - The MongoDB version string.
  * `status` - The support status of the version (one of `BETA`, `SUPPORTED`, `RECOMMENDED`, `DEPRECATED`).
  * `comment` - Additional information about the version status.
  * `creation_freeze_date` - The date (RFC 3339) after which new clusters can no longer be created with this version. Empty if not set.
  * `deprecation_date` - The date (RFC 3339) on which this version becomes deprecated. Empty if not set.
  * `can_upgrade_to` - List of versions this version can be upgraded to. Only versions listed here are valid upgrade targets when changing the `version` field. Downgrades are not supported.
