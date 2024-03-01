---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_k8s_clusters"
sidebar_current: "docs-ionoscloud-datasource-clusters"
description: |-
  Retrieves a list of Ionos Cloud Kubernetes Clusters
---

# ionoscloud\_servers

The **k8s_clusters data source** can be used to search for and return existing kubernetes clusters based on filters used.

## Example Usage

### By Name
```hcl
data ionoscloud_k8s_clusters example {
  filter {
    name = "name"
    value = "k8sClusterExample"
  }
}
```

### By Name and k8s version Family
```hcl
data ionoscloud_k8s_clusters example2 {
  filter {
    name = "name"
    value = "k8sClusterExample"
  }
  filter {
    name = "k8s_version"
    value = "1.27"
  }
}
```


### Retrieve private clusters only, by Name and Cluster State
```hcl
data ionoscloud_servers example {
  filter{
    name = "name"
    value = "k8sClusterExample"
  }
  filter {
    name = "state"
    value = "ACTIVE"
  }
  filter {
    name = "public"
    value = "false"
  }
}
```

## Argument Reference

* `filter` -  (Optional) One or more property name - value pairs to be used in filtering the cluster list by the specified attributes. You can use most of the top level fields from the  [k8s_cluster](../data-sources/k8s_cluster.md) resource **except** those containing other nested structures such as `maintenance_window` or `config`.

**NOTE:** Filtering uses partial matching for all types of values. Searching for a cluster using `name:testCluster` will find all clusters who have the `testCluster` substring in their name. This also applies to values for properties that would normally be boolean or numerical.

## Attributes Reference

The following attributes are returned by the datasource:

* `clusters` - list of Kubernetes clusters that match the provided filters. The elements of this list are structurally identical to the `k8s_cluster` datasource, which is limited to retrieving only 1 cluster in a single query.
* `entries` - indicates the number of clusters found and added to the list after the query has been performed with the specified filters.
For a full reference of all the attributes returned, check out [documentation](../data-sources/k8s_cluster.md)