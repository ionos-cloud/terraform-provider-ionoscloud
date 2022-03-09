---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_group"
sidebar_current: "docs-ionoscloud-datasource-group"
description: |-
  Get information on a Ionos Cloud Groups
---

# ionoscloud\_group

The group data source can be used to search for and return existing groups.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned. 
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_group" "group_example" {
  name			= "my_group"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing group that you want to search for.
* `id` - (Optional) ID of the group you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the LAN.
* `name` - A name for the group.
* `create_datacenter` - The group will be allowed to create virtual data centers.
* `create_snapshot` - The group will be allowed to create snapshots.
* `reserve_ip` - The group will be allowed to reserve IP addresses.
* `access_activity_log` - The group will be allowed to access the activity log.
* `create_pcc` - The group will be allowed to create pcc privilege.
* `s3_privilege` - The group will have S3 privilege.
* `create_backup_unit` - The group will be allowed to create backup unit privilege.
* `create_internet_access` - The group will be allowed to create internet access privilege.
* `create_k8s_cluster` - The group will be allowed to create kubernetes cluster privilege.
* `users` - List of users in group.