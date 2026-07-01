---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_clusters_v2"
sidebar_current: "docs-datasource-mariadb_clusters_v2"
description: |-
  Lists IONOS CLOUD MariaDB V2 Clusters.
---

# ionoscloud_mariadb_clusters_v2

The `ionoscloud_mariadb_clusters_v2` data source can be used to retrieve information about existing MariaDB V2 clusters in a given location, with an optional name filter.

## Example Usage

```hcl
data "ionoscloud_mariadb_clusters_v2" "all" {
  location = "de/txl"
}

data "ionoscloud_mariadb_clusters_v2" "filtered" {
  location = "de/txl"
  name     = "my-cluster"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `de/fra`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `name` - (Optional)[string] Filter clusters by name (**partial match** — the value is passed directly to the API's name filter, so it matches any cluster name containing the given string, not only an exact match).

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of clusters. Each item exposes the same attributes as the [ionoscloud_mariadb_cluster_v2](mariadb_cluster_v2.md) data source. Note that `credentials.password` is not available — only `credentials.username` and `credentials.database` are returned by the API.
