---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_clusters_v2"
sidebar_current: "docs-datasource-inmemorydb_clusters_v2"
description: |-
  Lists IONOS Cloud In-Memory DB v2 Clusters.
---

# ionoscloud_inmemorydb_clusters_v2 (Data Source)

Lists IONOS Cloud In-Memory DB v2 Clusters in a given location, with optional name filter.

## Example Usage

```hcl
data "ionoscloud_inmemorydb_clusters_v2" "all" {
  location = "de/txl"
}

data "ionoscloud_inmemorydb_clusters_v2" "filtered" {
  location = "de/txl"
  name     = "my-cluster"
}
```

## Argument Reference

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `location` | String | Yes | The location to query. |
| `name` | String | No | Filter by name (case-insensitive contains match). |

## Attributes Reference

| Attribute | Description |
|-----------|-------------|
| `items` | List of cluster objects. Each item contains all cluster attributes (see `ionoscloud_inmemorydb_cluster_v2` resource for the full schema). |
