---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_group"
sidebar_current: "docs-ionoscloud-datasource-group"
description: |-
  Get information on a Ionos Cloud Groups
---

# ionoscloud_group

The **Group data source** can be used to search for and return existing groups.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned. 
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_group" "example" {
  id			= "group_id"
}
```

### By Name
```hcl
data "ionoscloud_group" "example" {
  name			= "Group Example"
}
```

## Argument Reference

* `name` - (Optional) Name of an existing group that you want to search for.
* `id` - (Optional) ID of the group you want to search for.

Either `name` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the group.
* `name` - A name for the group.
* `create_datacenter` - The group will be allowed to create virtual data centers.
* `create_snapshot` - The group will be allowed to create snapshots.
* `reserve_ip` - The group will be allowed to reserve IP addresses.
* `access_activity_log` - The group will be allowed to access the activity log.
* `create_pcc` - The group will be allowed to create Cross Connects privilege.
* `s3_privilege` - The group will have S3 privilege.
* `create_backup_unit` - The group will be allowed to create backup unit privilege.
* `create_internet_access` - The group will be allowed to create internet access privilege.
* `create_k8s_cluster` - The group will be allowed to create kubernetes cluster privilege.
* `create_flow_log` -  The group will be allowed to create flow log.
* `access_and_manage_monitoring`  The group will be allowed to access and manage monitoring.
* `access_and_manage_certificates` - The group will be allowed to access and manage certificates.
* `access_and_manage_dns` - The group will be allowed to access and manage dns records.
* `manage_registry` - The group will be allowed to access container registry related functionality.
* `manage_dataplatform` - The group will be allowed to access and manage the Data Platform.
* `access_and_manage_logging` - The group will be allowed to access and manage logging.
* `access_and_manage_cdn` - The group will be allowed to access and manage cdn.
* `access_and_manage_vpn` - The group will be allowed to access and manage vpn.
* `access_and_manage_api_gateway` - The group will be allowed to access and manage api gateway.
* `access_and_manage_kaas` - The group will be allowed to access and manage kaas.
* `access_and_manage_network_file_storage` - The group will be allowed to access and manage network file storage.
* `access_and_manage_ai_model_hub` - The group will be allowed to access and manage ai model hub.
* `access_and_manage_iam_resources` - The group will be allowed to access and manage iam resources.
* `create_network_security_groups` - The group will be allowed to create network security groups.
* `manage_dbaas` - Privilege for a group to manage DBaaS related functionality.
* `users` - List of users in group.