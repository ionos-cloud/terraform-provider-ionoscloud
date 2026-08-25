---
subcategory: "Database as a Service - MariaDB v2"
layout: "ionoscloud"
page_title: "IONOS CLOUD: ionoscloud_mariadb_cluster_v2"
sidebar_current: "docs-resource-mariadb_cluster_v2"
description: |-
  Creates and manages an IONOS CLOUD MariaDB V2 Cluster.
---

# ionoscloud_mariadb_cluster_v2

Manages a DBaaS MariaDB V2 Cluster.

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

resource "ionoscloud_mariadb_cluster_v2" "example" {
  location    = "de/txl"
  name        = "my-mariadb-cluster"
  description = "MariaDB cluster"
  version     = "11.4"

  logs_enabled    = true
  metrics_enabled = true

  instances = {
    count        = 1
    cores        = 1
    ram          = 4
    storage_size = 10
  }

  connections = {
    datacenter_id            = ionoscloud_datacenter.example.id
    lan_id                   = ionoscloud_lan.example.id
    primary_instance_address = "192.168.2.101/24"
  }

  backup = {
    location       = "eu-central-3"
    retention_days = 7
  }

  maintenance_window = {
    time            = "09:00:00"
    day_of_the_week = "Sunday"
  }

  credentials = {
    username = "dbadmin"
    password = "P@ssw0rd123!"
    database = "my_database"
  }
}
```

## Argument Reference

* `location` - (Required)[string] The location of the cluster. Changing this forces a new resource. Available: `de/fra`, `de/fra/1`, `de/fra/2`, `de/txl`, `es/vit`, `fr/par`, `gb/bhx`, `gb/lhr`, `us/ewr`, `us/las`, `us/mci`.
* `name` - (Required)[string] The name of the MariaDB cluster.
* `version` - (Required)[string] The MariaDB version for the cluster. See the [`ionoscloud_mariadb_versions_v2`](../data-sources/mariadb_versions_v2.md) data source for the list of supported versions and their available upgrade paths.
* `description` - (Optional)[string] Human-readable description for the cluster.
* `logs_enabled` - (Optional)(Computed)[bool] Allows or disallows the collection and reporting of logs for this cluster's observability. If not set, the API default is used.
* `metrics_enabled` - (Optional)(Computed)[bool] Allows or disallows the collection and reporting of metrics for this cluster's observability. If not set, the API default is used.
* `instances` - (Required)[object] Compute and storage configuration for each instance in the cluster.
  * `count` - (Required)[int] The total number of instances in the cluster (one primary and n-1 secondary).
  * `cores` - (Required)[int] The number of CPU cores per instance.
  * `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB).
  * `storage_size` - (Required)[int] The amount of storage per instance in gigabytes (GB).
* `connections` - (Required)[object] Connection information of the MariaDB cluster. Only one connection is allowed.
  * `datacenter_id` - (Required)[string] The datacenter to connect your instance to.
  * `lan_id` - (Required)[string] The numeric LAN ID to connect your instance to.
  * `primary_instance_address` - (Required)[string] The IP address and netmask of the cluster's primary instance, in CIDR notation.
* `backup` - (Required)[object] Configures backup location and retention.
  * `location` - (Required)[string] The Object Storage location where the backup will be created. Changing this forces re-creation of the cluster.
  * `retention_days` - (Required)[int] Configures how many days cluster backups are retained.
* `maintenance_window` - (Required)[object] A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - (Required)[string] Start of the maintenance window in UTC time (`HH:MM:SS`).
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `credentials` - (Required)[object] Credentials for the initial database user to be created.
  * `username` - (Required)[string] The username of the initial MariaDB user.
  * `password` - (Required)[string] **Sensitive.** The password for the initial MariaDB user. Not returned by the API — will be null in state after `terraform import`.
  * `database` - (Required)[string] The name of the initial database to be created.
* `restore_from_backup` - (Optional)[object] Restores the cluster from a backup.
  * `source_backup_id` - (Optional)[string] UUID of the backup to restore from. Required when `restore_from_backup` is set during cluster creation; not valid for in-place restore during an update.
  * `recovery_target_datetime` - (Optional)[string] ISO 8601 timestamp causing the system to replay backups up to the specified time. Optional for create-time restore; required for in-place restore during an update.

> **Note:** `restore_from_backup` is not returned by the API. The values are stored in state as configured but will be null after `terraform import`.

* `timeouts` - (Optional) Standard Terraform timeouts: `create`, `update`, `delete`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID (UUID) of the cluster.
* `dns_name` - The DNS name used to access the cluster.

## Import

MariaDB V2 clusters can be imported using `<location>:<cluster_id>`:

```shell
terraform import ionoscloud_mariadb_cluster_v2.example de/txl:example-id
```

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can be used with the `identity` attribute:

```hcl
import {
  to = ionoscloud_mariadb_cluster_v2.example
  identity = {
    id       = "example-id"
    location = "de/txl"
  }
}

resource "ionoscloud_mariadb_cluster_v2" "example" {
  ### Configuration omitted for brevity ###
}
```
