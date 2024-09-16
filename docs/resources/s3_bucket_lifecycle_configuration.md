---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_lifecycle_configuration"
sidebar_current: "docs-resource-s3_bucket_lifecycle_configuration"
description: |-
  Manages Buckets lifecycle configuration on IonosCloud.
---

# ionoscloud_s3_bucket_lifecycle_configuration

Manages Lifecycle Configuration for Buckets on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_lifecycle_configuration" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  rule {
    id     = "Logs delete"
    status = "Enabled"

    prefix = "/logs"

    expiration {
      days = 90
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `lifecycle_rule` - (Required)[block] A block of lifecycle_rule as defined below.
  - `id` - (Optional)[int] Container for the Contract Number of the owner
  - `prefix` - (Required)[string] Prefix identifying one or more objects to which the rule applies.
  - `status` - (Required)[string] The lifecycle rule status. Valid values are `Enabled` or `Disabled`.
  - `expiration` - (Optional)[block]  A lifecycle rule for when an object expires.
    - `days` - (Optional)[int] Specifies the number of days after object creation when the object expires. Required if 'date' is not specified.
    - `date` - (Optional)[string] Specifies the date after which you want the specific rule action to take effect.
    - `expired_object_delete_marker` - (Optional)[bool] Indicates whether IONOS S3 Object Storage will remove a delete marker with no noncurrent versions. If set to true, the delete marker will be expired; if set to false the policy takes no operation. This cannot be specified with Days or Date in a Lifecycle Expiration Policy.
  - `noncurrent_version_expiration` - (Optional)[block] A lifecycle rule for when non-current object versions expire.
    - `noncurrent_days` - (Optional)[int] Specifies the number of days an object is noncurrent before Amazon S3 can perform the associated action.
  - `abort_incomplete_multipart_upload` - (Optional)[block] Specifies the days since the initiation of an incomplete multipart upload that IONOS S3 Object Storage will wait before permanently removing all parts of the upload.
    - `days_after_initiation` - (Optional)[int] Specifies the number of days after which IONOS S3 Object Storage aborts an incomplete multipart upload.

## Import

S3 Bucket lifecycle configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_lifecycle_configuration.example example
```