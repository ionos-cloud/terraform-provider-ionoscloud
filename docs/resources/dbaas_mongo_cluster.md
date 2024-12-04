---
subcategory: "Database as a Service - MongoDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mongo_cluster"
sidebar_current: "docs-resource_mongo_cluster"
description: |-
  Creates and manages DbaaS MongoDB Cluster objects.
---

# ionoscloud\_mongo_cluster

Manages a **DbaaS Mongo Cluster**.

## Example Usage for playground or business editions. They require template_id defined.

```hcl
resource "ionoscloud_datacenter" "datacenter_example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan"  "lan_example" {
  datacenter_id           = ionoscloud_datacenter.datacenter_example.id
  public                  = false
  name                    = "example"
}

resource "ionoscloud_mongo_cluster" "example_mongo_cluster" {
  maintenance_window {
    day_of_the_week  = "Sunday"
  time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 1
  display_name = "example_mongo_cluster"
  location = ionoscloud_datacenter.datacenter_example.location
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    cidr_list            =  ["192.168.1.108/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
}

resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```
## Example Usage for enterprise edition

**Enterprise Support: With MongoDB Enterprise, you gain access to professional support from the MongoDB team ensuring that you receive timely assistance and expert guidance when needed. IONOS offers enterprise-grade Service Level Agreements (SLAs), guaranteeing rapid response times and 24/7 support to address any critical issues that may arise.**

```hcl
resource "ionoscloud_datacenter" "datacenter_example" {
  name                    = "example"
  location                = "de/txl"
  description             = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan"  "lan_example" {
  datacenter_id           = ionoscloud_datacenter.datacenter_example.id
  public                  = false
  name                    = "example"
}

resource ionoscloud_mongo_cluster "example_mongo_cluster" {
  maintenance_window {
    day_of_the_week  = "Sunday"
    time             = "09:00:00"
  }
  mongodb_version = "5.0"
  instances          = 3
  display_name = "example_mongo_cluster"
  location = ionoscloud_datacenter.datacenter_example.location
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    cidr_list       =  ["192.168.1.108/24", "192.168.1.109/24", "192.168.1.110/24"]
  }
  type = "sharded-cluster"
  shards = 2
  edition = "enterprise"
  ram = 2048
  cores = 1
  storage_size = 5120
  storage_type = "HDD"
}

resource "random_password" "cluster_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `edition` - (Optional)(Computed)[string] Cluster edition. Playground, business or enterprise.
* `mongodb_version` - (Required)[string] The MongoDB version of your cluster. Updates to the value of the field force the cluster to be re-created.
* `template_id` - (Optional)[string] The unique ID of the template, which specifies the number of cores, storage size, and memory. Updates to the value of the field force the cluster to be re-created. Required for playground and business editions. Must not be provided for enterprise edition.
* `instances` - (Required)[int] The total number of instances in the cluster (one master and n-1 standbys). Example: 1, 3, 5, 7. Updates to the value of the field force the cluster to be re-created.
* `display_name` - (Required)[string] The name of your cluster. Updates to the value of the field force the cluster to be re-created.
* `location` - (Required)[string] The physical location where the cluster will be created. Property cannot be modified after datacenter creation (disallowed in update requests). Available locations: de/txl, gb/lhr, es/vit. Update forces cluster re-creation.
* `connections` - (Required)[List] Details about the network connection for your cluster. Updates to the value of the field force the cluster to be re-created.
    * `datacenter_id` - (Required)[string] The datacenter to connect your cluster to.
    * `lan_id` - (Required)[string] The LAN to connect your cluster to.
    * `cidr_list` - (Required)[List] The list of IPs and subnet for your cluster. Note the following unavailable IP ranges:10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. example: [192.168.1.100/24, 192.168.1.101/24]. See [Private IPs](https://www.ionos.com/help/server-cloud-infrastructure/private-network/private-ip-address-ranges/) and [Cluster Setup - Preparing the network](https://docs.ionos.com/cloud/databases/mongodb/api-howtos/create-a-cluster#preparing-the-network).
* `maintenance_window` - (Optional)(Computed) A weekly 4 hour-long window, during which maintenance might occur.  Updates to the value of the field force the cluster to be re-created.
    * `time` - (Required)[string]
    * `day_of_the_week` - (Required)[string]
* `connection_string` - (Computed)[string] The physical location where the cluster will be created. This will be where all of your instances live. Updates to the value of the field force the cluster to be re-created. Available locations: de/txl, gb/lhr, es/vit
* `ram` - (Optional)(Computed)[int]The amount of memory per instance in megabytes. Required for enterprise edition.
* `storage_size` - (Optional)(Computed)[int] The amount of storage per instance in MB. Required for enterprise edition.
* `storage_type` - (Optional)(Computed)[String] The storage type used in your cluster. Required for enterprise edition.
* `cores` - (Optional)(Computed)[int] The number of CPU cores per replica. Required for enterprise edition.
* `shards` - (Optional)[int]The total number of shards in the cluster.
* `type` - (Optional)(Computed)[string]The cluster type, either `replicaset` or `sharded-cluster`.
* `bi_connector` - (Optional)(Computed)The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands to aid in data analysis.
  * `enabled`: (Optional)[bool] - The status of the BI Connector. If not set, the BI Connector is disabled. 
  * `host`: (Computed)[string] - The host where this new BI Connector is installed.
  * `port`: (Computed)[string] - Port number used when connecting to this new BI Connector.
* `backup` - (Optional)[list]
  * `location`: (Optional)[string] - The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster. Possible values are de, eu-south-2, or eu-central-2.


## Import

Resource DbaaS MongoDb Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import ionoscloud_mongo_cluster.mycluser {cluster uuid}
```
