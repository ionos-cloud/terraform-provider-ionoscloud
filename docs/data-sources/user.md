---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_user"
sidebar_current: "docs-ionoscloud-datasource-user"
description: |-
  Get information on a Ionos Cloud Users
---

# ionoscloud\_user

The **User data source** can be used to search for and return existing users.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_user" "example" {
  id			= <user_id>
}
```

### By Email
```hcl
data "ionoscloud_user" "example" {
  email			= "example@email.com"
}
```

### By Email from Env Variables - Current User
data "ionoscloud_user" "example" {
}

## Argument Reference

* `email` - (Optional) Email of an existing user that you want to search for.
* `id` - (Optional) ID of the user you want to search for.

Either `email` or `id` can be provided. If no argument is set, the provider will search for the **email that was provided for the configuration**. If none is found, the provider will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the user.
* `administrator` - Indicates if the user has administrative rights. Administrators do not need to be managed in groups, as they automatically have access to all resources associated with the contract.
* `email` - The e-mail address for the user.
* `first_name` - The first name for the user.
* `force_sec_auth` - Indicates if secure (two-factor) authentication should be forced for the user (true) or not (false).
* `last_name` - The last name for the user.
* `password` - The password for the user.
* `sec_auth_active` - Indicates if secure authentication is active for the user or not
* `s3_canonical_user_id` - Canonical (S3) id of the user for a given identity
* `active` - Indicates if the user is active
* `groups` - Shows the id and name of the groups that the user is a member of