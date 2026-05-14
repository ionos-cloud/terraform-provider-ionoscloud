---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: user_object_storage_bucket"
sidebar_current: "docs-ionoscloud-datasource-user_object_storage_bucket"
description: |-
  Get information about IONOS User-Owned Object Storage Buckets.
---

# ionoscloud_user_object_storage_bucket

The **User Object Storage Bucket data source** can be used to look up an existing user-owned bucket by name and region.

> ⚠️ **Deprecation notice:** User-owned buckets are a legacy bucket type. Use [`ionoscloud_s3_bucket`](../resources/s3_bucket.md) (contract-owned) for new workloads.

## Example Usage

```hcl
data "ionoscloud_user_object_storage_bucket" "example" {
  name   = "my-bucket"
  region = "de"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the bucket.
- `region` - (Optional)[string] The region where the bucket resides. Defaults to `de` (Frankfurt). Valid values: `de`, `eu-central-2`, `eu-south-2`. Must match the bucket's actual region or the lookup will return not found.

## Attributes Reference

The following attributes are returned by the data source:

- `name` - The name of the bucket.
- `object_lock_enabled` - Whether Object Lock is enabled for the bucket.
- `region` - The region of the bucket.
- `tags` - Tags assigned to the bucket.
