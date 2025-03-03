---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: group"
sidebar_current: "docs-resource-group"
description: |-
  Creates and manages group objects.
---

# ionoscloud_group

Manages **Groups** and **Group Privileges** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_user" "example1" {
  first_name              = "user1"
  last_name               = "user1"
  email                   = "unique_email.com"
  password                = random_password.user1_password.result
  administrator           = false
  force_sec_auth          = false
}

resource "ionoscloud_user" "example2" {
  first_name              = "user2"
  last_name               = "user2"
  email                   = "unique_email.com"
  password                = random_password.user2_password.result
  administrator           = false
  force_sec_auth          = false
}

resource "ionoscloud_group" "example" {
  name                                   = "Group Example"
  create_datacenter                      = true
  create_snapshot                        = true
  reserve_ip                             = true
  access_activity_log                    = true
  create_pcc                             = true
  s3_privilege                           = true
  create_backup_unit                     = true
  create_internet_access                 = true
  create_k8s_cluster                     = true
  create_flow_log                        = true
  access_and_manage_monitoring           = true
  access_and_manage_certificates         = true
  access_and_manage_logging              = true
  access_and_manage_cdn                  = true
  access_and_manage_vpn                  = true
  access_and_manage_api_gateway          = true
  access_and_manage_kaas                 = true
  access_and_manage_network_file_storage = true
  access_and_manage_ai_model_hub         = true
  access_and_manage_iam_resources        = true
  create_network_security_groups         = true
  manage_dbaas                           = true
  user_ids                               = [ ionoscloud_user.example1.id, ionoscloud_user.example2.id ] 
}

resource "random_password" "user1_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "random_password" "user2_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `name` - (Required) [string] A name for the group.
* `create_datacenter` - (Optional) [Boolean] The group will be allowed to create virtual data centers.
* `create_snapshot` - (Optional) [Boolean] The group will be allowed to create snapshots.
* `reserve_ip` - (Optional) [Boolean] The group will be allowed to reserve IP addresses.
* `access_activity_log` - (Optional) [Boolean] The group will be allowed to access the activity log.
* `create_pcc` - (Optional) [Boolean] The group will be allowed to create Cross Connects privilege.
* `s3_privilege` - (Optional) [Boolean] The group will have S3 privilege.
* `create_backup_unit` - (Optional) [Boolean] The group will be allowed to create backup unit privilege.
* `create_internet_access` - (Optional) [Boolean] The group will be allowed to create internet access privilege.
* `create_k8s_cluster` - (Optional) [Boolean]  The group will be allowed to create kubernetes cluster privilege.
* `create_flow_log` - (Optional) [Boolean]  The group will be allowed to create flow log.
* `access_and_manage_monitoring` - (Optional) [Boolean]  The group will be allowed to access and manage monitoring.
* `access_and_manage_certificates` - (Optional) [Boolean]  The group will be allowed to access and manage certificates.
* `access_and_manage_dns` - (Optional) [Boolean]  The group will be allowed to access and manage dns records.
* `manage_registry` - (Optional) [Boolean]  The group will be allowed to access container registry related functionality.
* `manage_dataplatform` - (Optional) [Boolean]  The group will be allowed to access and manage the Data Platform.
* `access_and_manage_logging` - (Optional) [Boolean]  The group will be allowed to access and manage logging.
* `access_and_manage_cdn` - (Optional) [Boolean]  The group will be allowed to access and manage cdn.
* `access_and_manage_vpn` - (Optional) [Boolean]  The group will be allowed to access and manage vpn.
* `access_and_manage_api_gateway` - (Optional) [Boolean]  The group will be allowed to access and manage api gateway.
* `access_and_manage_kaas` - (Optional) [Boolean]  The group will be allowed to access and manage kaas.
* `access_and_manage_network_file_storage` - (Optional) [Boolean]  The group will be allowed to access and manage network file storage.
* `access_and_manage_ai_model_hub` - (Optional) [Boolean]  The group will be allowed to access and manage ai model hub.
* `access_and_manage_iam_resources` - (Optional) [Boolean]  The group will be allowed to access and manage iam resources.
* `create_network_security_groups` - (Optional) [Boolean]  The group will be allowed to create network security groups.
* `manage_dbaas` - (Optional) [Boolean]  Privilege for a group to manage DBaaS related functionality.
* `user_ids` - (Optional) [list] A list of users to add to the group.
* `user_id` - (Optional) [string] The ID of the specific user to add to the group. Please use user_ids argument since this is **DEPRECATED**
* `users` - (Computed) List of users - See the [User](user.md) section

**NOTE:** user_id/user_ids field cannot be used at the same time with group_ids field in user resource. Trying to add the same user to the same group in both ways in the same plan will result in a cyclic dependency error.

## Import

Resource Group can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_group.mygroup group uuid
```

> :warning: **If you are upgrading to v6.2.0**: You have to modify you plan for user_ids to match the new structure, by renaming the field old field, **user_id**, to user_ids and put the old value into an array. This is not backwards compatible.
