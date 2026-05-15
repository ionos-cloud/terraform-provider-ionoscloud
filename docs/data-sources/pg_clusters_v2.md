---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD : ionoscloud_pg_clusters_v2"
sidebar_current: "docs-ionoscloud_pg_clusters_v2"
description: |-
  List DBaaS PostgreSQL v2 Clusters
---

# ionoscloud_pg_clusters_v2

The **DBaaS PostgreSQL v2 Clusters data source** can be used to list existing DBaaS PostgreSQL v2 Clusters in a given location.
An optional name filter can be used to narrow down results.

## Example Usage

### List all clusters in a location
```hcl
data "ionoscloud_pg_clusters_v2" "example" {
  location = "de/txl"
}
```

### Filter clusters by name
```hcl
data "ionoscloud_pg_clusters_v2" "example" {
  location = "de/txl"
  name     = "cluster_name"
}
```

## Argument Reference

* `location` - (Required)[string] The region in which to look up clusters. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `name` - (Optional)[string] Filters clusters by name. Matches cluster names that contain the provided string.

## Attributes Reference

The following attributes are returned by the datasource:

* `clusters` - The list of PostgreSQL v2 clusters. Each cluster has the following attributes:
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
