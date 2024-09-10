---
subcategory: "Database as a Service - MongoDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mongo_user"
sidebar_current: "docs-resource_mongo_user"
description: |-
  Creates and manages DbaaS MongoDB users.
---

# ionoscloud\_mongo_user

Manages a **DbaaS Mongo User**. .

## Example Usage

```hcl
# Basic example

resource "ionoscloud_datacenter" "datacenter_example" {
  name        = "example"
  location    = "de/txl"
  description = "Datacenter for testing dbaas cluster"
}

resource "ionoscloud_lan" "lan_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  public        = false
  name          = "example"
}

resource "ionoscloud_mongo_cluster" "example_mongo_cluster" {
  maintenance_window {
    day_of_the_week = "Sunday"
    time            = "09:00:00"
  }
  mongodb_version = "5.0"
  instances       = 1
  display_name    = "example_mongo_cluster"
  location        = ionoscloud_datacenter.datacenter_example.location
  connections {
    datacenter_id = ionoscloud_datacenter.datacenter_example.id
    lan_id        = ionoscloud_lan.lan_example.id
    cidr_list = ["192.168.1.108/24"]
  }
  template_id = "6b78ea06-ee0e-4689-998c-fc9c46e781f6"
}

resource "ionoscloud_mongo_user" "example_mongo_user" {
  cluster_id = ionoscloud_mongo_cluster.example_mongo_cluster.id
  username   = "myUser"
  password   = "strongPassword"
  roles {
    role     = "read"
    database = "db1"
  }
  roles {
    role     = "readWrite"
    database = "db2"
  }
}
```

```hcl
# Complete example

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

resource "random_password" "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "ionoscloud_mongo_user" "example_mongo_user" {
  cluster_id = ionoscloud_mongo_cluster.example_mongo_cluster.id
  username = "myUser"
  password = random_password.user_password.result
  roles {
    role = "read"
    database = "db1"
  }
  roles {
    role = "readWrite"
    database = "db2"
  }
}
```

## Argument reference

* `cluster_id` - (Required)[string] The unique ID of the cluster. Updates to the value of the field force the cluster to be re-created.
* `username` - (Required)[string] Used for authentication. Updates to the value of the field force the cluster to be re-created.
* `database` - (Required)[string] The user database to use for authentication. Updates to the value of the field force the cluster to be re-created.
* `password` - (Required)[string] User password. Updates to the value of the field force the cluster to be re-created.
* `roles` - (Required)[string] a list of mongodb user roles. Updates to the value of the field force the cluster to be re-created.
    * `role` - (Required)[true] Mongodb user role. Examples: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor.
    * `database` - (Required)[true] Database on which to apply the role.

**NOTE:** MongoDb users do not support update at the moment. Changing any attribute will result in the user being re-created.

## Import

Resource DBaaS MongoDB User can be imported using the `clusterID` and the `username`.
First, define an empty resource in the plan:
```hcl
resource "ionoscloud_mongo_user" "importeduser" {
  
}
```
Then you can import the user using the following command:
```shell
terraform import ionoscloud_mongo_user.mycluser {clusterId}/{username}
```
