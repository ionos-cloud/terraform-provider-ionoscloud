---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_server_side_encryption_configuration"
sidebar_current: "docs-resource-s3_bucket_server_side_encryption_configuration"
description: |-
  Manages Buckets server side encryption configuration on IonosCloud.
---

# ionoscloud_s3_bucket_server_side_encryption_configuration

Manages Server Side Encryption Configuration for Buckets on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_server_side_encryption_configuration" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `rule` - (Required)[block] A block of rule as defined below.
  - `apply_server_side_encryption_by_default` - (Required)[block] Defines the default encryption settings.
    - `sse_algorithm` - (Required)[string] Server-side encryption algorithm to use. Valid values are 'AES256'
## Import

S3 Bucket server side encryption configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_server_side_encryption_configuration.example example
```