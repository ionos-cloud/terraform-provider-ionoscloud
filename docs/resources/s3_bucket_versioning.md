---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_versioning"
sidebar_current: "docs-resource-s3_bucket_versioning"
description: |-
  Manages Buckets versioning on IonosCloud.
---

# ionoscloud_s3_versioning

Manages **Buckets versioning** on IonosCloud.

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_versioning" "example"{
  bucket = ionoscloud_s3_bucket.example.name
  versioning_configuration {
    status = "Enabled"
  }
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `versioning_configuration` - (Required)[block] A block of versioning_configuration as defined below.
  - `status` - (Required)[string] The versioning state of the bucket. Can be `Enabled` or `Suspended`.
  - `mfa_delete` - (Optional)[string] Specifies whether MFA delete is enabled or not. Can be `Enabled` or `Disabled`.

## Import

IONOS Object Storage Bucket Versioning can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_versioning.example example
```