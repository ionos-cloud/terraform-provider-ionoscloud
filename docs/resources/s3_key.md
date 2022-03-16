---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: s3_key"
sidebar_current: "docs-resource-s3-key"
description: |-
  Creates and manages IonosCloud S3 keys.
---

# ionoscloud_s3_key

Manages an **S3 Key** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_user" "example" {
    first_name              = "example"
    last_name               = "example"
    email                   = <unique_email>
    password                = "abc123-321CBA"
    administrator           = false
    force_sec_auth          = false
}

resource "ionoscloud_s3_key" "example" {
    user_id                 = ionoscloud_user.example.id
    active                  = true
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the S3 Key.
- `active` - (Optional)[boolean] Whether the S3 is active / enabled or not - Please keep in mind this is only required on create. Default value in true
- `secret_key` - (Computed)  The S3 Secret key.

## Import

An S3 Unit resource can be imported using its user id as well as its `resource id`, e.g.

```shell
terraform import ionoscloud_s3_key.demo {userId}/{s3KeyId}
```

This can be helpful when you want to import S3 Keys which you have already created manually or using other means, outside of terraform.
