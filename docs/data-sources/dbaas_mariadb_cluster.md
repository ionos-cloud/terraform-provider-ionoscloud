---
subcategory: "Database as a Service - MariaDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mariadb_cluster"
sidebar_current: "docs-datasource-mariadb_cluster"
description: |-
  Get information on a DBaaS MariaDB Cluster
---

# ionoscloud_mariadb_cluster

The **DBaaS MariaDB Cluster data source** can be used to search for and return an existing DBaaS MariaDB Cluster.

## Example Usage

### By ID 
```hcl
data "ionoscloud_mariadb_cluster" "example" {
  id = <cluster_id>
}
```

### By Name

```hcl
data "ionoscloud_mariadb_cluster" "example" {
  display_name = "MariaDB_cluster"
}
```

## Argument Reference

* `display_name` - (Optional)[string] Display Name of an existing cluster that you want to search for.
* `id` - (Optional)[string] ID of the cluster you want to search for.

Either `display_name` or `id` must be provided. If none or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `mariadb_version` - [string] The MariaDB version of your cluster.
* `instances` - [int] The total number of instances in the cluster (one primary and n-1 secondary).
* `cores` - [int] The number of CPU cores per instance.
* `ram` - [int] The amount of memory per instance in gigabytes (GB).
* `storage_size` - [int] The amount of storage per instance in gigabytes (GB).
* `connections` - The network connection for your cluster. Only one connection is allowed.
  * `datacenter_id` - [string] The datacenter to connect your cluster to.
  * `lan_id` - [string] The LAN to connect your cluster to.
  * `cidr` - [string] The IP and subnet for your cluster.
* `display_name` - [string] The friendly name of your cluster.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - [string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - [string] The name of the week day.
* `dns_name` - [string] The DNS name pointing to your cluster.
  