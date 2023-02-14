---
subcategory: "Database as a Service - MongoDB"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_mongo_cluster"
sidebar_current: "docs-ionoscloud_mongo_cluster"
description: |-
  Get information on DbaaS MongoDB Cluster objects.
---

# ionoscloud\_mongo_cluster

The **DbaaS Mongo Cluster data source** can be used to search for and return an existing DbaaS MongoDB Cluster.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_mongo_cluster" "example" {
  id	= <cluster_id>
}
```
### By display_name

```hcl
data "ionoscloud_mongo_cluster" "example" {
  display_name	= <display_name>
}
```

* `display_name` - (Optional) Display Name of an existing cluster that you want to search for.
* `id` - (Optional) ID of the cluster you want to search for.

Either `display_name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `mongodb_version` - The MongoDB version of your cluster. Updates to the value of the field force the cluster to be re-created.
* `template_id` - The unique ID of the template, which specifies the number of cores, storage size, and memory. Updates to the value of the field force the cluster to be re-created.
* `instances` - The total number of instances in the cluster (one master and n-1 standbys). Example: 3, 5, 7. Updates to the value of the field force the cluster to be re-created.
* `display_name` - The name of your cluster. Updates to the value of the field force the cluster to be re-created.
* `location` - The connection string for your cluster. Updates to the value of the field force the cluster to be re-created.
* `connections` - Details about the network connection for your cluster. Updates to the value of the field force the cluster to be re-created.
    * `datacenter_id` - The datacenter to connect your cluster to.
    * `lan_id` -  The LAN to connect your cluster to.
    * `cidr` - The IP and subnet for the database. Must be same number as instances. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. Please input in the correct format like IP/Subnet, exp: 192.168.10.0/24. See [Private IPs](https://www.ionos.com/help/server-cloud-infrastructure/private-network/private-ip-address-ranges/) and [Cluster Setup - Preparing the network](https://docs.ionos.com/reference/product-information/api-automation-guides/database-as-a-service/create-a-database#preparing-the-network).
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur.  Updates to the value of the field force the cluster to be re-created.
    * `time` 
    * `day_of_the_week`
* `credentials` - Credentials for the database user to be created. This attribute is immutable(disallowed in update requests). Updates to the value of the field force the cluster to be re-created.
    * `username` -  The username for the initial mongoDB user.
    * `password`
* `connection_string` - The physical location where the cluster will be created. This will be where all of your instances live. Updates to the value of the field force the cluster to be re-created. Available locations: de/txl, gb/lhr, es/vit"
