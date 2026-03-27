---
subcategory: "Database as a Service - PostgreSQL v2"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_cluster_v2"
sidebar_current: "docs-resource-pg_cluster_v2"
description: |-
  Creates and manages DBaaS PostgreSQL v2 Cluster objects.
---

# ionoscloud_pg_cluster_v2

Manages a DBaaS PostgreSQL v2 Cluster.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name     = "example"
  location = "de/txl"
}

resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public        = false
  name          = "example"
}

resource "ionoscloud_pg_cluster_v2" "example" {
  name              = "PostgreSQL_cluster"
  description       = "Production PostgreSQL cluster"
  version           = "17"
  location          = "de/txl"
  backup_location   = "de"
  replication_mode  = "ASYNCHRONOUS"
  connection_pooler = "ENABLED"
  logs_enabled      = true
  metrics_enabled   = true

  instances = {
    count        = 1
    cores        = 2
    ram          = 4
    storage_size = 10
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.example.id
    lan_id                   = ionoscloud_lan.example.id
    primary_instance_address = "192.168.1.100/24"
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "username"
    password = random_password.cluster_password.result
    database = "mydb"
  }

  restore_from_backup = {
    source_backup_id         = "backup_uuid"
    recovery_target_datetime = "2024-12-06T13:54:08Z"
  }

  timeouts = {
    create = "60m"
    update = "60m"
    delete = "60m"
  }
}

resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  override_special = "@$!%*?&"
}
```

## Argument reference

* `name` - (Required)[string] The name of the PostgreSQL cluster.
* `description` - (Optional)[string] Human-readable description for the cluster.
* `version` - (Optional)(Computed)[string] The PostgreSQL version of the cluster. If omitted, the API assigns a default version.
* `location` - (Required)[string] The location of the PostgreSQL cluster. This is used for routing to the regional API endpoint. Changing this value will destroy the existing cluster and create a new one in the specified location. Available locations: `de/fra`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `backup_location` - (Required)[string] The S3 location where the backups will be created. Supported locations are provided by the `ionoscloud_pg_backup_location` data source.
* `replication_mode` - (Required)[string] Replication mode across the instances. Possible values: `ASYNCHRONOUS`, `STRICTLY_SYNCHRONOUS`.
* `connection_pooler` - (Optional)(Computed)[string] Defines how database connections are managed and reused. Possible values: `DISABLED`, `TRANSACTION`, `SESSION`.
* `logs_enabled` - (Optional)(Computed)[bool] Enables or disables the collection and reporting of logs for observability of this cluster.
* `metrics_enabled` - (Optional)(Computed)[bool] Enables or disables the collection and reporting of metrics for observability of this cluster.
* `instances` - (Required)[object] The instance configuration for the PostgreSQL cluster.
  * `count` - (Required)[int] The total number of instances in the cluster (one primary and n-1 secondary).
  * `cores` - (Required)[int] The number of CPU cores per instance.
  * `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB).
  * `storage_size` - (Required)[int] The amount of storage per instance in gigabytes (GB).
* `connections` - (Required)[object] Connection information of the PostgreSQL cluster.
  * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
  * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
  * `primary_instance_address` - (Required)[string] The IP and netmask that will be assigned to the cluster primary instance.
* `maintenance_window` - (Required)[object] A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `credentials` - (Required)[object] Credentials for the master database user to be created.
  * `username` - (Required)[string] The username of the master database user.
  * `password` - (Required, WriteOnly)[string] The password for the master database user. This value is never stored in state (requires Terraform 1.11+).
  * `database` - (Required)[string] The name of the initial database to be created.
* `restore_from_backup` - (Optional)[object] Configures the cluster to be initialized with data from an existing backup.
  * `source_backup_id` - (Optional)[string] The UUID of the backup to restore data from.
  * `recovery_target_datetime` - (Optional)[string] If supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.
* `dns_name` - (Computed)[string] The DNS name used to access the cluster.

## Timeouts

This resource supports the following `Timeouts` configuration options:

* `create` - (Default `60m`) Time to wait for the cluster to be provisioned.
* `read` - (Default `60m`) Time to wait for the cluster to be read.
* `update` - (Default `60m`) Time to wait for the cluster to be updated.
* `delete` - (Default `60m`) Time to wait for the cluster to be deleted.

## Import

Resource DBaaS PostgreSQL v2 Cluster can be imported using the format `location:cluster_id`, e.g.

```shell
terraform import ionoscloud_pg_cluster_v2.mycluster de/txl:cluster_uuid
```
