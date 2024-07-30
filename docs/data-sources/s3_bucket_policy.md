---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_policy"
sidebar_current: "docs-ionoscloud-datasource-s3_bucket_policy"
description: |-
    Get information about IonosCloud S3 Buckets policies.
---

# ionoscloud_s3_bucket_policy

The **Bucket Policy data source** can be used to search for and return existing bucket policies.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl

data "ionosclud_s3_bucket_policy" "example" {
  bucket = "example"
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.

## Attributes Reference

The following attributes are returned by the datasource:

- `policy` - The policy of the bucket.