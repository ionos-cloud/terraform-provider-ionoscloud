---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: user"
sidebar_current: "docs-resource-user"
description: |-
  Creates and manages user objects.
---

# ionoscloud_user

Manages **Users** and list users and groups associated with that user.

## Example Usage

```hcl
resource "ionoscloud_user" "example" {
  first_name              = "example"
  last_name               = "example"
  email                   = "unique@email.com"
  password                = random_password.user_password.result
  administrator           = false
  force_sec_auth          = false
  active                  = true
  group_ids 		          = [ ionoscloud_group.group1.id, ionoscloud_group.group2.id, ionoscloud_group.group3.id]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group2" {
  name = "group2"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group3" {
  name = "group3"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "random_password" "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
```

## Argument reference

* `administrator` - (Required)[Boolean] Indicates if the user has administrative rights. Administrators do not need to be managed in groups, as they automatically have access to all resources associated with the contract.
* `email` - (Required)[string] An e-mail address for the user.
* `first_name` - (Required)[string] A first name for the user.
* `force_sec_auth` - (Required)[Boolean] Indicates if secure (two-factor) authentication should be forced for the user (true) or not (false).
* `last_name` - (Required)[string] A last name for the user.
 * `password` - (Optional)[string] A password for the user. If you are using terraform 1.11 or higher, you can use `password_wo` instead of `password` to avoid storing the password in the state file.
* `password_wo` - (Optional, Write-Only, exactly one of `password` or `password_wo` is required) user password. This value is always marked as sensitive in the plan output, regardless of `type`. Additionally, `write-only` values are never stored to state. `password_wo_version` can be used to trigger an update and is required with this argument. In Terraform CLI version 0.15 and later, this may require additional configuration handling for certain scenarios. For more information, see the [Terraform v0.15 Upgrade Guide](https://www.terraform.io/upgrade-guides/0-15.html#sensitive-output-values).
* `password_wo_version` - (Optional) Used together with `password_wo` to trigger an update. Increment this value when an update to the `password_wo` is required.
* `sec_auth_active` - (Optional)[Boolean] Indicates if secure authentication is active for the user or not. *it can not be used in create requests - can be used in update*
* `s3_canonical_user_id` - (Computed) Canonical (IONOS Object Storage) id of the user for a given identity
* `active` - (Optional)[Boolean] Indicates if the user is active
* `group_ids` - (Optional)[Set] The groups that this user will be a member of
**NOTE:** Group_ids field cannot be used at the same time with user_ids field in group resource. Trying to add the same user to the same group in both ways in the same plan will result in a cyclic dependency error.
**NOTE:** `password_wo` requires Teraform 1.11 or higher.
## Import

Resource User can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_user.myuser user uuid
```