---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IONOS CLOUD: s3_bucket"
description: |-
  Lists IONOS Object Storage Buckets.
---

# List Resource: ionoscloud_s3_bucket

-> **Note:** List Resources are supported in HashiCorp Terraform version 1.14 and later.

Lists [IONOS Object Storage Buckets](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage) on IONOS CLOUD.

## Example Usage

``` hcl
list "ionoscloud_s3_bucket" "example" {
  provider = ionoscloud
}
```

## Argument Reference

This list resource has no configuration arguments.

## Attributes Reference

Each result exposes the following attributes, matching the `ionoscloud_s3_bucket` resource schema:

- `id` - Name of the bucket (same as `name`).
- `name` - The name of the bucket.
- `region` - The region where the bucket is located.
- `object_lock_enabled` - Whether object lock is enabled for the bucket.
- `tags` - A mapping of tags assigned to the bucket.

> **Note:** `force_destroy` is not available via the list resource as it is a local Terraform-only attribute and not returned by the API.
