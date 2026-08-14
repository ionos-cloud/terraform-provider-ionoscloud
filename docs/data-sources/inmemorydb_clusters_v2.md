---
subcategory: "Database as a Service - InMemoryDB V2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_inmemorydb_clusters_v2"
sidebar_current: "docs-datasource-inmemorydb_clusters_v2"
description: |-
  Lists IONOS CLOUD InMemoryDB V2 Clusters.
---

# ionoscloud_inmemorydb_clusters_v2

The `ionoscloud_inmemorydb_clusters_v2` data source can be used to retrieve information about existing InMemoryDB V2 clusters in a given location, with an optional name filter.

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

* `location` - (Required)[string] The location to query.
* `name` - (Optional)[string] Filter by name (case-insensitive contains match).

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of clusters. Each item exposes the same attributes as the [ionoscloud_inmemorydb_cluster_v2](inmemorydb_cluster_v2.md) data source. Note that `credentials.password` is not available — only `credentials.username` is returned by the API.
