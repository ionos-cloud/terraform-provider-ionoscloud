---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_user"
sidebar_current: "docs-ionoscloud-datasource-user"
description: |-
Get information on a Ionos Cloud Users
---

# ionoscloud\_user

The users data source can be used to search for and return existing users.

## Example Usage

```hcl
data "ionoscloud_user" "user_example" {
  email			= "my_user@email.com"
}
```

## Argument Reference

* `email` - (Optional) Email of an existing user that you want to search for.
* `id` - (Optional) ID of the user you want to search for.

Either `email` or `id` must be provided. If none, or both are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the LAN.
* `administrator` - The group has permission to edit privileges on this resource.
* `email` - The e-mail address for the user.
* `first_name` - The first name for the user.
* `force_sec_auth` - Indicates if secure (two-factor) authentication should be enabled for the user (true) or not (false).
* `last_name` - The last name for the user.
* `password` - The password for the user.
* `sec_auth_active` - Indicates if secure authentication is active for the user or not
* `s3_canonical_user_id` - Canonical (S3) id of the user for a given identity
* `active` - Indicates if the user is active