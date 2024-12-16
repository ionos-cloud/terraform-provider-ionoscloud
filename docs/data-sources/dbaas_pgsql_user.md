---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_user"
sidebar_current: "docs-ionoscloud-datasource-pg_user"
description: |-
  Get information on DBaaS PgSql User.
---

# ionoscloud_pg_user

The **PgSql User data source** can be used to search for and return an existing PgSql user.

## Example Usage

```hcl
data "ionoscloud_pg_user" "example" {
   cluster_id = "cluster_id"
   username   = "username"
}
```

## Argument Reference

* `cluster_id` - (Required)[string] The ID of the cluster.
* `username` - (Required)[string] Name of an existing user that you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - [string] The id of the user.
* `is_system_user` - [bool] Describes whether this user is a system user or not. A system user cannot be updated or deleted.

