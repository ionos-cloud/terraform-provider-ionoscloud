---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_database"
sidebar_current: "docs-ionoscloud-datasource-pg_database"
description: |-
  Get information on DBaaS PgSql Database.
---

# ionoscloud\_pg_database

The **PgSql Database data source** can be used to search for and return an existing PgSql database.

## Example Usage

```hcl
data "ionoscloud_pg_database" "example" {
   cluster_id = "cluster_id"
   name   = "databasename"
}
```

## Argument Reference

* `cluster_id` - (Required)[string] The ID of the cluster.
* `name` - (Required)[string] Name of an existing database that you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - [string] The id of the database.
* `owner` - [string] The owner of the database.
