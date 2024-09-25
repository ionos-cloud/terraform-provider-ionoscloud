---
subcategory: "S3 management"
layout: "ionoscloud"
page_title: "IonosCloud: s3_accesskey"
sidebar_current: "docs-resource-s3_accesskey"
description: |-
  Creates and manages IonosCloud S3 Accesskeys.
---

# ionoscloud_s3_accesskey

Manages an **S3 Accesskey** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_accesskey" "example" {
    description = "my description"
}
```

## Argument Reference

The following arguments are supported:

- `description` - (Optional)[string] Description of the Access key.
- `id` - (Computed)  The ID (UUID) of the AccessKey.
- `accesskey` - (Computed)  Access key metadata is a string of 92 characters.
- `secretkey` - (Computed)  The secret key of the Access key.
- `canonical_user_id` - (Computed)  The canonical user ID which is valid for user-owned buckets.
- `contract_user_id` - (Computed)  The contract user ID which is valid for contract-owned buckets
- `timeouts` - (Optional) Timeouts for this resource.
  - `create` - (Optional)[string] Time to wait for the bucket to be created. Default is `10m`.
  - `delete` - (Optional)[string] Time to wait for the bucket to be deleted. Default is `10m`.
- `force_destroy` - (Optional)[bool] If true, the bucket and the contents of the bucket will be destroyed. Default is `false`.

## Import

An S3 accesskey resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_s3_accesskey.demo {s3AccesskeyId}
```

This can be helpful when you want to import S3 Accesskeys which you have already created manually or using other means, outside of terraform.
