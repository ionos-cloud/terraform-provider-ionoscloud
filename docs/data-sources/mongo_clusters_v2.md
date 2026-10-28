---
subcategory: "Database as a Service - MongoDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mongo_clusters_v2"
sidebar_current: "docs-datasource-mongo_clusters_v2"
description: |-
  Lists IONOS CLOUD MongoDB V2 Clusters.
---

# ionoscloud_mongo_clusters_v2

The `ionoscloud_mongo_clusters_v2` data source can be used to retrieve information about existing MongoDB V2 clusters in a given location, with an optional name and edition filter.

## Example Usage

```hcl
data "ionoscloud_mongo_clusters_v2" "all" {
  location = "gb/bhx"
}

data "ionoscloud_mongo_clusters_v2" "filtered" {
  location = "gb/bhx"
  name     = "my-cluster"
  edition  = "ENTERPRISE"
}
```

## Argument Reference

* `location` - (Required)[string] The location to query. Available locations: `gb/bhx`, `us/ewr`.
* `name` - (Optional)[string] Filter clusters by name (partial match).
* `edition` - (Optional)[string] Filter clusters by edition (`ENTERPRISE` or `BUSINESS`).

## Attributes Reference

The following attributes are returned by the datasource:

* `items` - List of clusters. Each item exposes the same attributes as the [ionoscloud_mongo_cluster_v2](mongo_cluster_v2.md) data source. Note that `credentials.password` is never returned by the API.
