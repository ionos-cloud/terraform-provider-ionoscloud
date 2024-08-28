---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_access_block"
sidebar_current: "docs-resource-s3_bucket_access_block"
description: |-
  Creates and manages IonosCloud S3 Public Access Block for buckets.
---

# ionoscloud_s3_public_access_block

Manages **public acccess for Buckets** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket_public_access_block" "example"{
  bucket = ionoscloud_s3_bucket.example.name
  ignore_public_acls = true
  restrict_public_buckets = true
  block_public_policy = false
  block_public_acls = false
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `ignore_public_acls` - (Optional)[bool] Instructs the system to ignore any ACLs that grant public access. Even if ACLs are set to allow public access, they will be disregarded.
- `restrict_public_buckets` - (Optional)[bool] Restricts access to buckets that have public policies. Buckets with policies that grant public access will have their access restricted.
- `block_public_policy` - (Optional)[bool] Blocks public access to the bucket via bucket policies. Bucket policies that grant public access will not be allowed.
- `block_public_acls` - (Optional)[bool] Indicates that access to the bucket via Access Control Lists (ACLs) that grant public access is blocked. In other words, ACLs that allow public access are not permitted.
## Import

Resource Bucket access block can be imported using the `bucket name`

```shell
terraform import ionoscloud_s3_bucket_public_access_block.example example
```