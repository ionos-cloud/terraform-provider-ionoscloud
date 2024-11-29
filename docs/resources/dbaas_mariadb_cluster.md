---
subcategory: "Database as a Service - MariaDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mariadb_cluster"
sidebar_current: "docs-resource-mariadb_cluster"
description: |-
  Creates and manages DBaaS MariaDB Cluster objects.
---

# ionoscloud_mariadb_cluster

Manages a **DBaaS MariaDB Cluster**. 

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for testing DBaaS cluster"
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
    size                  = 10
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

resource "ionoscloud_mariadb_cluster" "example" {
  mariadb_version         = "10.6"
  location                = "de/txl"
  instances               = 1
  cores                   = 4
  ram                     = 4
  storage_size            = 10
  connections   {
    datacenter_id         =  ionoscloud_datacenter.example.id 
    lan_id                =  ionoscloud_lan.example.id 
    cidr                  =  local.database_ip_cidr
  }
  display_name            = "MariaDB_cluster"
  maintenance_window {
    day_of_the_week       = "Sunday"
    time                  = "09:00:00"
  }
  credentials {
    username              = "username"
    password              = random_password.cluster_password.result
  }
}
resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `mariadb_version` - (Required)[string] The MariaDB version of your cluster. Cannot be downgraded.
* `instances` - (Required)[int] The total number of instances in the cluster (one primary and n-1 secondary).
* `location`- (Optional)[string] The location in which the cluster will be created. Different service endpoints are used based on location, possible options are: "de/fra", "de/txl", "es/vit", "fr/par", "gb/lhr", "us/ewr", "us/las", "us/mci". If not set, the endpoint will be the one corresponding to "de/txl".
* `cores` - (Required)[int] The number of CPU cores per instance.
* `ram` - (Required)[int] The amount of memory per instance in gigabytes (GB).
* `storage_size` - (Required)[int] The amount of storage per instance in gigabytes (GB).
* `connections` - (Required) The network connection for your cluster. Only one connection is allowed.
  * `datacenter_id` - (Required)[true] The datacenter to connect your cluster to.
  * `lan_id` - (Required)[true] The numeric LAN ID to connect your cluster to.
  * `cidr` - (Required)[true] The IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Please enter in the correct format like IP/Subnet, exp: 192.168.10.0/24. See [Private IPs](https://www.ionos.com/help/server-cloud-infrastructure/private-network/private-ip-address-ranges/) and [Configuring the network](https://docs.ionos.com/cloud/compute-engine/networks/how-tos/configure-networks).
* `display_name` - (Required)[string] The friendly name of your cluster.
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur
  * `time` - (Required)[string] Start of the maintenance window in UTC time.
  * `day_of_the_week` - (Required)[string] The name of the week day.
* `credentials` - (Required) Credentials for the database user to be created.
    * `username` - (Required)[string] The username for the initial MariaDB user. Some system usernames are restricted (e.g 'mariadb', 'admin', 'standby').
    * `password` - (Required)[string] The password for a MariaDB user.
* `dns_name` - (Computed)[string] The DNS name pointing to your cluster.

> **⚠ WARNING:** `IONOS_API_URL_MARIADB` can be used to set a custom API URL for the MariaDB Cluster. `location` field needs to be empty, otherwise it will override the custom API URL. Setting `token` or `IONOS_API_URL` does not have any effect.

## Import

Resource DBaaS MariaDB Cluster can be imported using the `cluster_id` and the `location`, separated by `:`, e.g.

```shell
terraform import ionoscloud_mariadb_cluster.mycluster {location}:{cluster UUID}
```
