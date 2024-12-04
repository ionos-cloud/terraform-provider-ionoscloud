---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_cluster"
sidebar_current: "docs-resource-pg_cluster"
description: |-
  Creates and manages DbaaS Postgres Cluster objects.
---

# ionoscloud\_pg_cluster

Manages a **DbaaS PgSql Cluster**.

## Example Usage

```hcl
# Basic example

resource "ionoscloud_datacenter" "example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan"  "example" {
  datacenter_id           = ionoscloud_datacenter.example.id
  public                  = false
  name                    = "example"
}

resource "ionoscloud_pg_cluster" "example" {
  postgres_version        = "12"
  instances               = 1
  cores                   = 4
  ram                     = 2048
  storage_size            = 2048
  storage_type            = "HDD"
  connection_pooler {
    enabled = true
    pool_mode = "session"
  }
  connections   {
    datacenter_id         =  ionoscloud_datacenter.example.id
    lan_id                =  ionoscloud_lan.example.id
    cidr                  =  "192.168.100.1/24"
  }
  location                = ionoscloud_datacenter.example.location
  display_name            = "PostgreSQL_cluster"
  maintenance_window {
    day_of_the_week       = "Sunday"
    time                  = "09:00:00"
  }
  credentials {
    username              = "username"
    password              = "strongPassword"
  }
  synchronization_mode    = "ASYNCHRONOUS"
}
```

```hcl
# Complete example

resource "ionoscloud_datacenter" "example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan"  "example" {
  datacenter_id           = ionoscloud_datacenter.example.id 
  public                  = false
  name                    = "example"
}

resource "ionoscloud_server" "example" {
  name                    = "example"
  datacenter_id           = ionoscloud_datacenter.example.id
  cores                   = 2
  ram                     = 2048
  availability_zone       = "ZONE_1"
  cpu_family              = "INTEL_SKYLAKE"
  image_name              = "rockylinux-8-GenericCloud-20230518"
  image_password          = "password"
  volume {
    name                  = "example"
    size                  = 6
    disk_type             = "SSD Standard"
  }
  nic {
    lan                   = ionoscloud_lan.example.id
    name                  = "example"
    dhcp                  = true
  }
}

locals {
 prefix                   = format("%s/%s", ionoscloud_server.example.nic[0].ips[0], "24")
 database_ip              = cidrhost(local.prefix, 1)
 database_ip_cidr         = format("%s/%s", local.database_ip, "24")
}

resource "ionoscloud_pg_cluster" "example" {
  postgres_version        = "12"
  instances               = 1
  cores                   = 4
  ram                     = 2048
  storage_size            = 2048
  storage_type            = "HDD"
  connection_pooler {
    enabled = true
    pool_mode = "session"
  }
  connections   {
    datacenter_id         =  ionoscloud_datacenter.example.id 
    lan_id                =  ionoscloud_lan.example.id 
    cidr                  =  local.database_ip_cidr
  }
  location                = ionoscloud_datacenter.example.location
  display_name            = "PostgreSQL_cluster"
  maintenance_window {
    day_of_the_week       = "Sunday"
    time                  = "09:00:00"
  }
  credentials {
    username              = "username"
    password              = random_password.cluster_password.result
  }
  synchronization_mode    = "ASYNCHRONOUS"
  from_backup {
    backup_id             = "backup_uuid"
    recovery_target_time  = "2021-12-06T13:54:08Z"
  }
}
resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `postgres_version` - (Required)[string] The PostgreSQL version of your cluster.
* `instances` - (Required)[int] The total number of instances in the cluster (one master and n-1 standbys)
* `cores` - (Required)[int] The number of CPU cores per replica.
* `ram` - (Required)[int] The amount of memory per instance in megabytes. Has to be a multiple of 1024.
* `storage_size` - (Required)[int] The amount of storage per instance in MB. Has to be a multiple of 2048.
* `storage_type` - (Required)[string] SSD, SSD Standard, SSD Premium, or HDD. Value "SSD" is deprecated, use the equivalent "SSD Premium" instead. This attribute is immutable(disallowed in update requests).
* `connection_pooler` - (Optional)[object]
  * `enabled` - (Required)[bool]
  * `pool_mode` - (Required)[string] Represents different modes of connection pooling for the connection pooler.
* `connections` - (Required)[string] Details about the network connection for your cluster.
  * `datacenter_id` - (Required)[true] The datacenter to connect your cluster to.
  * `lan_id` - (Required)[true] The LAN to connect your cluster to.
  * `cidr` - (Required)[true] The IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Please enter in the correct format like IP/Subnet, exp: 192.168.10.0/24. See [Private IPs](https://www.ionos.com/help/server-cloud-infrastructure/private-network/private-ip-address-ranges/) and [Configuring the network](https://docs.ionos.com/cloud/compute-engine/networks/how-tos/configure-networks).
* `location` - (Required)[string] The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation. Possible values are: `de/fra`, `de/txl`, `gb/lhr`, `es/vit`, `us/ewr`, `us/las`. This attribute is immutable(disallowed in update requests).
* `backup_location` - (Optional)(Computed)[string] The IONOS Object Storage location where the backups will be stored. Possible values are: `de`, `eu-south-2`, `eu-central-2`. This attribute is immutable (disallowed in update requests).
* `display_name` - (Required)[string] The friendly name of your cluster.
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur
  * `time` - (Required)[string]
  * `day_of_the_week` - (Required)[string]
* `credentials` - (Required)[string] Credentials for the database user to be created. This attribute is immutable(disallowed in update requests).
    * `username` - (Required)[string] The username for the initial postgres user. Some system usernames are restricted (e.g. "postgres", "admin", "standby")
    * `password` - (Required)[string]
* `synchronization_mode` - (Required) [string] Represents different modes of replication. Can have one of the following values: ASYNCHRONOUS, SYNCHRONOUS, STRICTLY_SYNCHRONOUS. This attribute is immutable(disallowed in update requests).
* `from_backup` - (Optional)[string] The unique ID of the backup you want to restore. This attribute is immutable(disallowed in update requests).
  * `backup_id` - (Required)[string] The unique ID of the backup you want to restore.
  * `recovery_target_time` - (Optional)[string] If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely.
* `dns_name` - (Computed)[string] The DNS name pointing to your cluster.

## Import

Resource DbaaS Postgres Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import ionoscloud_pg_cluster.mycluser {cluster uuid}
```
