---
subcategory: "Database as a Service - In-Memory DB v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_inmemorydb_cluster_v2"
sidebar_current: "docs-datasource-inmemorydb_cluster_v2"
description: |-
  Reads an IONOS Cloud In-Memory DB v2 Cluster by ID or name.
---

# ionoscloud_inmemorydb_cluster_v2 (Data Source)

Reads an IONOS Cloud In-Memory DB v2 Cluster.

## Example Usage

```hcl
# Look up by ID
data "ionoscloud_inmemorydb_cluster_v2" "by_id" {
  id       = "e69b22a5-8fee-56b1-b6fb-4a07e4205ead"
  location = "de/txl"
}

# Look up by name
data "ionoscloud_inmemorydb_cluster_v2" "by_name" {
  name     = "my-inmemorydb-cluster"
  location = "de/txl"
}
```

## Argument Reference

Exactly one of `id` or `name` must be provided.

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | String | One of | The UUID of the cluster. |
| `name` | String | One of | The cluster name (exact match). |
| `location` | String | Yes | The location of the cluster. |

## Attributes Reference

All resource attributes are exported. See the resource documentation for the full list.
