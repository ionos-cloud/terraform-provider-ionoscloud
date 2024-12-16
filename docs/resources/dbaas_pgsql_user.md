---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_user"
sidebar_current: "docs-resource-pg_user"
description: |-
  Creates and manages DbaaS Postgres User objects.
---

# ionoscloud_pg_user

Manages a **DbaaS PgSql User**.

## Example Usage

Create a `PgSQL` cluster as presented in the documentation for the cluster, then define a user resource
and link it with the previously created cluster:

```hcl
resource "ionoscloud_pg_user" "example_pg_user" {
  cluster_id = ionoscloud_pg_cluster.example.id
  username = "exampleuser"
  password = random_password.user_password.result
}

resource "random_password" "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `cluster_id` - (Required)[string] The unique ID of the cluster. Updates to the value of the field force the cluster to be re-created.
* `username` - (Required)[string] Used for authentication. Updates to the value of the field force the cluster to be re-created.
* `password` - (Required)[string] User password.
* `is_system_user` - (Computed)[bool] Describes whether this user is a system user or not. A system user cannot be updated or deleted.

## Import

In order to import a PgSql user, you can define an empty user resource in the plan:

```hcl
resource "ionoscloud_pg_user" "example" {
  
}
```

The resource can be imported using the `clusterId` and the `username`, for example:

```shell
terraform import ionoscloud_pg_user.example {clusterId}/{username}
```
