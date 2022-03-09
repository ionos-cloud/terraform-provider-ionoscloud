---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: user"
sidebar_current: "docs-resource-user"
description: |-
  Creates and manages user objects.
---

# ionoscloud\_user

Manages users and list users and groups associated with that user.

## Example Usage

```hcl
resource "ionoscloud_user" "user" {
  first_name = "user"
  last_name = "user"
  email = "user@email.com"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}
```

## Argument reference

* `administrator` - (Required)[Boolean] The group has permission to edit privileges on this resource.
* `email` - (Required)[string] An e-mail address for the user.
* `first_name` - (Required)[string] A first name for the user.
* `force_sec_auth` - (Required)[Boolean] Indicates if secure (two-factor) authentication should be enabled for the user (true) or not (false).
* `last_name` - (Required)[string] A last name for the user.
* `password` - (Required)[string] A password for the user.
* `sec_auth_active` - (Optional)[Boolean] Indicates if secure authentication is active for the user or not. *it can not be used in create requests - can be used in update*
* `s3_canonical_user_id` - (Computed) Canonical (S3) id of the user for a given identity
* `active` - (Optional)[Boolean] Indicates if the user is active


## Import

Resource User can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_user.myuser {user uuid}
```