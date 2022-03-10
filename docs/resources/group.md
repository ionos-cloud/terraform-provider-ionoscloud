---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: group"
sidebar_current: "docs-resource-group"
description: |-
  Creates and manages group objects.
---

# ionoscloud\_group

Manages **Groups** and **Group Privileges** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_user" "example1" {
  first_name              = "user1"
  last_name               = "user1"
  email                   = <unique_email>
  password                = "abc123-321CBA"
  administrator           = false
  force_sec_auth          = false
}

resource "ionoscloud_user" "example2" {
  first_name              = "user2"
  last_name               = "user2"
  email                   = <unique_email>
  password                = "abc123-321CBA"
  administrator           = false
  force_sec_auth          = false
}

resource "ionoscloud_group" "example" {
  name                    = "Group Example"
  create_datacenter       = true
  create_snapshot         = true
  reserve_ip              = true
  access_activity_log     = true
  create_pcc              = true
  s3_privilege            = true
  create_backup_unit      = true
  create_internet_access  = true
  create_k8s_cluster      = true
  user_ids                = [ ionoscloud_user.example1.id, ionoscloud_user.example2.id ] 
}
```

##Argument reference

* `name` - (Required) [string] A name for the group.
* `create_datacenter` - (Optional) [Boolean] The group will be allowed to create virtual data centers.
* `create_snapshot` - (Optional) [Boolean] The group will be allowed to create snapshots.
* `reserve_ip` - (Optional) [Boolean] The group will be allowed to reserve IP addresses.
* `access_activity_log` - (Optional) [Boolean] The group will be allowed to access the activity log.
* `create_pcc` - (Optional) [Boolean] The group will be allowed to create pcc privilege.
* `s3_privilege` - (Optional) [Boolean] The group will have S3 privilege.
* `create_backup_unit` - (Optional) [Boolean] The group will be allowed to create backup unit privilege.
* `create_internet_access` - (Optional) [Boolean] The group will be allowed to create internet access privilege.
* `create_k8s_cluster` - (Optional) [Boolean]  The group will be allowed to create kubernetes cluster privilege.
* `create_flow_log` - (Optional) [Boolean]  The group will be allowed to create flow log.
* `access_and_manage_monitoring` - (Optional) [Boolean]  The group will be allowed to access and manage monitoring.
* `access_and_manage_certificates` - (Optional) [Boolean]  The group will be allowed to access and manage certificates.
* `user_ids` - (Optional) [list] A list of users to add to the group.
* `users` - (Computed) List of users - See the user section

## Import

Resource Group can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_group.mygroup {group uuid}
```

> :warning: **If you are upgrading to v6.2.0**: You have to modify you plan for user_ids to match the new structure, by renaming the field old field, **user_id**, to user_ids and put the old value into an array. This is not backwards compatible.
