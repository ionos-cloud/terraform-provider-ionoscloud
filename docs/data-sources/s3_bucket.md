---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket"
sidebar_current: "docs-ionoscloud-datasource-s3_bucket"
description: |-
  Get information about IonosCloud S3 Buckets.
---

# ionoscloud_s3_bucket

The **Bucket data source** can be used to search for and return existing buckets.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl

data "ionoscloud_s3_bucket" "example" {
  name = "example"
}

```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The bucket name. [ 3 .. 63 ] characters

## Attributes Reference

The following attributes are returned by the datasource:

- `region` - The region where the bucket is located.