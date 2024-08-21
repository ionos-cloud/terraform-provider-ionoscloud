---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_access_block"
sidebar_current: "docs-ionoscloud-datasource-s3_bucket_access_block"
description: |-
    Get information about IonosCloud S3 Buckets access block.
---

# ionoscloud_s3_public_access_block

The **Bucket access block data source** can be used to search for and return existing bucket access block.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_s3_bucket_public_access_block" "example"{
  bucket = "example"
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.

## Attributes Reference

The following attributes are returned by the datasource:

- `block_public_acls` - Indicates that access to the bucket via Access Control Lists (ACLs) that grant public access is blocked. In other words, ACLs that allow public access are not permitted.
- `block_public_policy` - Blocks public access to the bucket via bucket policies. Bucket policies that grant public access will not be allowed.
- `ignore_public_acls` - Instructs the system to ignore any ACLs that grant public access. Even if ACLs are set to allow public access, they will be disregarded.
- `restrict_public_buckets` - Restricts access to buckets that have public policies. Buckets with policies that grant public access will have their access restricted.
