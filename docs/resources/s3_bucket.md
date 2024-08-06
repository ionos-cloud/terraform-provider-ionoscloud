---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket"
sidebar_current: "docs-resource-s3_bucket"
description: |-
  Creates and manages IonosCloud S3 Buckets.
---

# ionoscloud_s3_bucket

Manages **S3 Buckets** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_s3_bucket" "example" {
  name = "example"
  region = "eu-central-3"
}

```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The bucket name. [ 3 .. 63 ] characters
- `region` - (Optional)[string] Specifies the Region where the bucket will be created. Please refer to the list of available regions.

⚠️ **Note:** The name must be unique across all IONOS accounts in all S3 regions. The name should adhere to the following [restrictions](https://docs.ionos.com/cloud/storage-and-backup/s3-object-storage/concepts/buckets#naming-conventions).

## Import

Resource Bucket can be imported using the `bucket name`

```shell
terraform import ionoscloud_s3_bucket.example example
```