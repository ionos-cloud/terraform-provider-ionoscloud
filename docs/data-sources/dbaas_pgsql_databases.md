---
subcategory: "Database as a Service - Postgres"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_pg_databases"
sidebar_current: "docs-ionoscloud-datasource-pg_databases"
description: |-
  Get information on DBaaS PgSql Databases.
---

# ionoscloud_pg_databases

The **PgSql Databases data source** can be used to search for and return multiple existing PgSql databases.

## Example Usage

### All databases from a specific cluster
```hcl
data "ionoscloud_pg_databases" "example" {
   cluster_id = "cluster_id"
}
```

### Filter by owner
```hcl
data "ionoscloud_pg_databases" "example" {
   cluster_id = "cluster_id"
   owner = "owner"
}
```

## Argument Reference

* `cluster_id` - (Required)[string] The ID of the cluster.
* `owner` - (Optional)[string] Filter using a specific owner.

## Attributes Reference

The following attributes are returned by the datasource:

* `databases` - [list] A list that contains either all databases, either some of them (filter by owner). A database from list has the following format:
  * `name` - [string] The name of the database.
  * `owner` - [string] The owner of the database.
  * `id` - [string] The ID of the database.
