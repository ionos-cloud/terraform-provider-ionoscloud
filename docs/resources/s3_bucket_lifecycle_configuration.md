---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_lifecycle_configuration"
sidebar_current: "docs-resource-s3_bucket_lifecycle_configuration"
description: |-
  Manages Buckets lifecycle configuration on IonosCloud.
---

# ionoscloud_s3_bucket_lifecycle_configuration

Manages Lifecycle Configuration for Buckets on IonosCloud.

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_lifecycle_configuration" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  rule {
    id     = "1"
    status = "Enabled"
    filter {
      prefix = "/logs"
    }
    expiration {
      days = 90
    }
  }

  rule {
    id     = "2"
    status = "Enabled"
    filter {
      prefix = "/logs"
    }
    noncurrent_version_expiration {
      noncurrent_days = 90
    }
  }

  rule {
    id     = "3"
    status = "Enabled"
    filter {
      prefix = "/logs"
    }
    abort_incomplete_multipart_upload {
      days_after_initiation = 90
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `lifecycle_rule` - (Required)[block] A block of lifecycle_rule as defined below.
  - `id` - (Optional)[int] Container for the Contract Number of the owner
  - `prefix` - (Required)[string] DEPRECATED! This field does not do anything! Will be removed in a future version, use `filter` instead. Prefix identifying one or more objects to which the rule applies.
  - `filter - (Optional)[block] A filter identifying one or more objects to which the rule applies.`
    - `prefix` - (Optional)[string] Prefix identifying one or more objects to which the rule applies. Cannot be used at the same time as `prefix` in the lifecycle rule.
  - `status` - (Required)[string] The lifecycle rule status. Valid values are `Enabled` or `Disabled`.
  - `expiration` - (Optional)[block]  A lifecycle rule for when an object expires.
    - `days` - (Optional)[int] Specifies the number of days after object creation when the object expires. Required if 'date' is not specified.
    - `date` - (Optional)[string] Specifies the date after which you want the specific rule action to take effect.
    - `expired_object_delete_marker` - (Optional)[bool] Indicates whether IONOS Object Storage will remove a delete marker with no noncurrent versions. If set to true, the delete marker will be expired; if set to false the policy takes no operation. This cannot be specified with Days or Date in a Lifecycle Expiration Policy.
  - `noncurrent_version_expiration` - (Optional)[block] A lifecycle rule for when non-current object versions expire.
    - `noncurrent_days` - (Optional)[int] Specifies the number of days an object is noncurrent before the associated action can be performed.
  - `abort_incomplete_multipart_upload` - (Optional)[block] Specifies the days since the initiation of an incomplete multipart upload that IONOS Object Storage will wait before permanently removing all parts of the upload.
    - `days_after_initiation` - (Optional)[int] Specifies the number of days after which IONOS Object Storage aborts an incomplete multipart upload.

## Import

IONOS Object Storage Bucket lifecycle configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_lifecycle_configuration.example example
```