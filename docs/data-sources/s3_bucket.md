---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket"
sidebar_current: "docs-ionoscloud-datasource-s3_bucket"
description: |-
  Get information about IonosCloud IONOS Object Storage Buckets.
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

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.

## Attributes Reference

The following attributes are returned by the datasource:

- `region` - The region where the bucket is located.