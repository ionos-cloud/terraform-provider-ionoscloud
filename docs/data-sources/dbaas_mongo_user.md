---
subcategory: "Database as a Service - MongoDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mongo_user"
sidebar_current: "docs-ionoscloud_mongo_user"
description: |-
  Creates and manages DbaaS MongoDB users.
---

# ionoscloud\_mongo_user

The **DbaaS Mongo User data source** can be used to search for and return an existing DbaaS MongoDB User.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

⚠️ **Note:** DBaaS - MongoDB is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications. 
Please contact your sales representative or support for more information.

## Example Usage

### By display_name
```hcl
data "ionoscloud_mongo_user" "example" {
  cluster_id	= <cluster_id>
  display_name	= <display_name>
}
```

## Argument reference

* `cluster_id` - (Required)[string] The unique ID of the cluster. Updates to the value of the field force the cluster to be re-created.
* `username` - (Required)[string] Used for authentication. Updates to the value of the field force the cluster to be re-created.
* `database` - (Required)[string] The user database to use for authentication. Updates to the value of the field force the cluster to be re-created.
* `password` - (Required)[string] User password. Updates to the value of the field force the cluster to be re-created.
* `roles` - (Required)[string] a list of mongodb user roles. Updates to the value of the field force the cluster to be re-created.
    * `role` - (Required)[true] Mongodb user role. Examples: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase and clusterMonitor.
    * `database` - (Required)[true] Database on which to apply the role.

**NOTE:** MongoDb users do not support update at the moment. Changing any attribute will result in the user being re-created.

## Import

Resource DbaaS MongoDb User can be imported using the `cluster_id`, the `database` and the `username` e.g.

```shell
terraform import ionoscloud_mongo_cluster.mycluser {cluster uuid} {database} {username}
```
