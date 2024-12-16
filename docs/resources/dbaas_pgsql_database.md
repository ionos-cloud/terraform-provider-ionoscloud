---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_database"
sidebar_current: "docs-resource-pg_database"
description: |-
  Creates and manages DbaaS Postgres Database objects.
---

# ionoscloud_pg_database

Manages a **DbaaS PgSql Database**.

## Example Usage

Create a `PgSQL` cluster as presented in the documentation for the cluster, then define a database resource
and link it with the previously created cluster:

```hcl
resource "ionoscloud_pg_database" "example_pg_database" {
  cluster_id = ionoscloud_pg_cluster.example.id
  name = "exampledatabase"
  owner = "exampleuser"
}
```

## Argument reference

* `cluster_id` - (Required)[string] The unique ID of the cluster.
* `name` - (Required)[string] The name of the database.
* `owner` - (Required)[string] The owner of the database.

## Import

In order to import a PgSql database, you can define an empty database resource in the plan:

```hcl
resource "ionoscloud_pg_database" "example" {
  
}
```

The resource can be imported using the `clusterId` and the `name`, for example:

```shell
terraform import ionoscloud_pg_database.example {clusterId}/{name}
```
