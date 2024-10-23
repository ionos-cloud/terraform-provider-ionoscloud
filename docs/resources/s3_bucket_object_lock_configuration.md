---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_object_lock_configuration"
sidebar_current: "docs-resource-s3_bucket_object_lock_configuration"
description: |-
  Manages Buckets object_lock_configuration on IonosCloud.
---

# ionoscloud_s3_object_lock_configuration

Manages Object Lock Configuration for Buckets on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
  object_lock_enabled = true
}

resource "ionoscloud_s3_bucket_object_lock_configuration" "test" {
  bucket = ionoscloud_s3_bucket.example.name
  object_lock_enabled = "Enabled"
  rule {
    default_retention {
      mode = "GOVERNANCE"
      days = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `object_lock_enabled` - (Required)[Optional] The object lock configuration status of the bucket. Must be `Enabled`.
- `rule` - (Optional)[block] A block of rule as defined below.
  - `default_retention` - (Required)[block] A block of default_retention as defined below.
    - `mode` - (Optional)[string] The default retention mode of the bucket. Can be `GOVERNANCE` or `COMPLIANCE`.
    - `days` - (Optional)[int] The default retention period of the bucket in days.
    - `years` - (Optional)[int] The default retention period of the bucket in years.

Days and years are mutually exclusive. You can only specify one of them.
## Import

IONOS Object Storage Bucket object lock configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_object_lock_configuration.example example
```