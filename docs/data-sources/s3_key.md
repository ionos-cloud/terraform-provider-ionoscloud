---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: s3_key"
sidebar_current: "docs-resource-s3-key"
description: |-
  Get Information on a IonosCloud s3 key
---

# ionoscloud_s3_key

The **S3 key data source** can be used to search for and return an existing s3 key of an user.

## Example Usage

```hcl
data "ionoscloud_s3_key" "example" {
  id         = <s3_key_id>
  user_id    = <user-uuid>
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the S3 Key.
- `id` - (Required) ID of the s3 key you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the s3 key
* `active` - The state of the s3 key
* `user_id` - The ID of the user that owns the key
* `secret_key` - (Computed)The S3 Secret key.