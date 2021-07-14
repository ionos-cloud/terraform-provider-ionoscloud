---
layout: "ionoscloud"
page_title: "IonosCloud: s3_key"
sidebar_current: "docs-resource-s3-key"
description: |-
  Creates and manages IonosCloud S3 keys.
---

# ionoscloud_s3_key

Manages an S3 Key on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_key" "demo" {
  user_id    = <user-uuid>
  active     = true
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the S3 Key.
- `active` - (Required)[boolean] Whether the S3 is active / enabled or not - Please keep in mind this is only required on create.
- `secret_key` - (Computed) Whether this key should be active or not

## Import

An S3 Unit resource can be imported using its user id as well as its `resource id`, e.g.

```shell
terraform import ionoscloud_s3_key.demo {userId}/{s3KeyId}
```

This can be helpful when you want to import S3 Keys which you have already created manually or using other means, outside of terraform.
