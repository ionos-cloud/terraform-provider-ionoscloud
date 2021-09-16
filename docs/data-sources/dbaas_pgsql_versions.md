---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_dbaas_pgsql_versions"
sidebar_current: "docs-ionoscloud-dbaas_pgsql_versions"
description: |-
Get information on DbaaS PgSql Versions
---

# ionoscloud\_dbaas_pgsql_versions

The DbaaS PgSql Versions data source can be used to search for and retrieve list of available pgsql versions for a specific cluster or for all clusters.

## Example Usage

### Retrieve list of postgres versions for a specific cluster
```hcl
data "ionoscloud_dbaas_pgsql_versions" "test_ds_dbaas_versions" {
	cluster_id = ionoscloud_dbaas_pgsql_cluster.test_dbaas_cluster.id
}
```

### Retrieve list of postgres versions for all clusters
```hcl
data "ionoscloud_dbaas_pgsql_versions" "test_ds_dbaas_versions" {
}
```

## Argument Reference

* `cluster_id` - (Optional) The unique ID of the cluster.

If `cluster_id` is not provided the data source will return the list of postgres version for all cluster.

## Attributes Reference

The following attributes are returned by the datasource:

* `cluster_id` - Id of the cluster
* `postgres_version` - list of PostgreSQL versions.