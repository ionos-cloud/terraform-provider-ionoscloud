---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_versions_v2"
sidebar_current: "docs-datasource-inmemorydb_versions_v2"
description: |-
  Lists supported IONOS Cloud In-Memory DB v2 Versions.
---

# ionoscloud_inmemorydb_versions_v2 (Data Source)

Lists all supported IONOS Cloud In-Memory DB v2 versions in a given location.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_versions_v2" "example" {
  location = "de/txl"
}

output "recommended_version" {
  value = [
    for v in data.ionoscloud_inmemorydb_versions_v2.example.items :
    v.version if v.status == "RECOMMENDED"
  ]
}
```

## Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | The location to query. |

## Attributes Reference

| Attribute | Description |
|-----------|-------------|
| `items` | List of version objects, each with `version`, `status`, `comment`, and `can_upgrade_to`. |
