---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_pg_cluster_v2"
sidebar_current: "docs-ionoscloud_pg_cluster_v2"
description: |-
  Get information on a DBaaS PostgreSQL v2 Cluster
---

# ionoscloud_pg_cluster_v2

The **DBaaS PostgreSQL v2 Cluster data source** can be used to search for and return an existing DBaaS PostgreSQL v2 Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_pg_cluster_v2" "example" {
  id       = "cluster_id"
  location = "de/txl"
}
```

### By Name
```hcl
data "ionoscloud_pg_cluster_v2" "example" {
  name     = "cluster_name"
  location = "de/txl"
}
```

## Argument Reference

* `id` - (Optional)[string] ID of the cluster you want to search for.
* `name` - (Optional)[string] The name of an existing cluster that you want to search for.
* `location` - (Required)[string] The region in which to look up the cluster. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.

Either `id` or `name` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The ID (UUID) of the cluster.
* `name` - The name of the PostgreSQL cluster.
* `description` - Human-readable description of the cluster.
* `version` - The PostgreSQL version of the cluster.
* `dns_name` - The DNS name used to access the cluster.
* `location` - The location of the PostgreSQL cluster.
* `backup_location` - The S3 location where the backups are stored.
* `replication_mode` - Replication mode across the instances.
* `connection_pooler` - How database connections are managed and reused.
* `logs_enabled` - Whether the collection and reporting of logs is enabled for this cluster.
* `metrics_enabled` - Whether the collection and reporting of metrics is enabled for this cluster.
* `instances` - The instance configuration for the PostgreSQL cluster.
  * `count` - The total number of instances in the cluster (one primary and n-1 secondary).
  * `cores` - The number of CPU cores per instance.
  * `ram` - The amount of memory per instance in gigabytes (GB).
  * `storage_size` - The amount of storage per instance in gigabytes (GB).
* `connections` - Connection information of the PostgreSQL cluster.
  * `datacenter_id` - The datacenter the cluster is connected to.
  * `lan_id` - The numeric LAN ID the cluster is connected to.
  * `primary_instance_address` - The IP and netmask assigned to the cluster primary instance.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - Start of the maintenance window in UTC time.
  * `day_of_the_week` - The name of the week day.
