---
layout: "ionoscloud"
page_title: "IonosCloud: dbaas_pgsql_cluster"
sidebar_current: "docs-resource-dbaas_pgsql_cluster"
description: |-
Creates and manages DbaaS PgSql Cluster objects.
---

# ionoscloud\_dbaas_pgsql_cluster

Manages a DbaaS PgSql Cluster.

## Example Usage

```hcl
resource "ionoscloud_dbaas_pgsql_cluster" "example" {
  postgres_version   = 12
  replicas           = 1
  cpu_core_count     = 4
  ram_size           = "2Gi"
  storage_size       = "1Gi"
  storage_type       = "HDD"
  vdc_connections   {
	vdc_id          =  ionoscloud_datacenter.example.id 
    lan_id          =  ionoscloud_lan.example.id 
    ip_address      =  "192.168.1.100/24"
  }
  location = ionoscloud_datacenter.example.location
  display_name = "PostgreSQL_cluster"
  backup_enabled = true
  maintenance_window {
    weekday = "Sunday"
    time            = "09:00:00"
  }
  credentials {
  	username = "username"
	password = "password"
  }
}
```

## Argument reference

* `postgres_version` - (Required)[string] The PostgreSQL version of your cluster.
* `replicas` - (Required)[int] The number of replicas in your cluster.
* `cpu_core_count` - (Required)[int] The number of CPU cores per replica.
* `ram_size` - (Required)[string] The amount of memory per replica.
* `storage_size` - (Required)[string] The amount of storage per replica.
* `storage_type` - (Required)[string] The storage type used in your cluster.
* `vdc_connections` - (Required)[string] The VDC to connect to your cluster.
  * `vdc_id` - (Required)[true] 
  * `lan_id` - (Required)[true]
  * `ip_address` - (Optional)[true] The IP and subnet for the database. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24
* `location` - (Required)[string] The physical location where the cluster will be created. This will be where all of your instances live. Property cannot be modified after datacenter creation (disallowed in update requests)
* `display_name` - (Required)[string] The friendly name of your cluster.
* `backup_enabled` - (Optional)[string] Deprecated: backup is always enabled. Enables automatic backups of your cluster.
* `maintenance_window` - (Optional)[string] A weekly 4 hour-long window, during which maintenance might occur
  * `time` - (Required)[string]
  * `weekday` - (Required)[string]
* `credentials` - (Required)[string] Credentials for the database user to be created.
    * `username` - (Required)[string] The username for the initial postgres user. some system usernames are restricted (e.g. "postgres", "admin", "standby")
    * `password` - (Required)[string]
    
## Import

Resource DbaaS PgSql Cluster can be imported using the `cluster_id`, e.g.

```shell
terraform import ionoscloud_dbaas_pgsql_cluster.mycluser {cluster uuid}
```
