---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: user_object_storage_bucket"
sidebar_current: "docs-resource-user_object_storage_bucket"
description: |-
  Creates and manages IONOS User-Owned Object Storage Buckets.
---

# ionoscloud_user_object_storage_bucket

Manages user-owned IONOS Object Storage Buckets on IonosCloud.

> ⚠️ **Deprecation notice:** User-owned buckets are a legacy bucket type. IONOS recommends using **contract-owned buckets** ([`ionoscloud_s3_bucket`](s3_bucket.md)) for all new workloads. Contract-owned buckets offer broader regional availability, full Terraform sub-resource support (versioning, lifecycle, CORS, SSE, etc.), and are the actively developed offering. Migrate existing user-owned buckets to contract-owned buckets where possible.

## Example Usage

```hcl
resource "ionoscloud_user_object_storage_bucket" "example" {
  name   = "example-user-bucket"
  region = "de"

  force_destroy = true

  timeouts {
    create = "10m"
    delete = "10m"
  }
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The bucket name. Must be between 3 and 63 characters.
- `region` - (Optional)[string] The region where the bucket is created. Defaults to `de` (Frankfurt). Available regions:
  - `de` — Frankfurt, Germany (`https://s3.eu-central-1.ionoscloud.com`)
  - `eu-central-2` — Berlin, Germany (`https://s3.eu-central-2.ionoscloud.com`)
  - `eu-south-2` — Logroño, Spain (`https://s3.eu-south-2.ionoscloud.com`)
- `object_lock_enabled` - (Optional)[bool] Whether Object Lock is enabled for the bucket. Defaults to `false`. **Cannot be changed after creation** — changing this value forces a new resource.
- `force_destroy` - (Optional)[bool] Defaults to `false`. When set to `true`, all objects in the bucket are deleted before the bucket itself is destroyed, allowing Terraform to remove a non-empty bucket. **Use with caution** — this irreversibly deletes all bucket contents.
- `timeouts` - (Optional) Timeouts for this resource.
  - `create` - (Optional)[string] Time to wait for the bucket to be created. Default is `10m`.
  - `delete` - (Optional)[string] Time to wait for the bucket to be deleted. Default is `10m`.

## Attributes Reference

- `id` - (Computed) Same value as `name`.

> ⚠️ **Note:** The bucket name must be unique. Follow IONOS [bucket naming conventions](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage/concepts/buckets#naming-conventions).

## Limitations

- **Region cannot be read back from the API** after import due to a known SDK issue with the `GetBucketLocation` response XML model. The region defaults to `de` when importing by bucket name alone. Use the `region:name` import format to specify the correct region explicitly.
- **`object_lock_enabled` cannot be read back from the API** after import. It defaults to `false` in state. If your bucket was created with Object Lock enabled, set `object_lock_enabled = true` in your configuration and re-import.
- Tags and sub-resources (versioning, lifecycle, CORS, SSE, website, policy) are not currently supported for user-owned buckets. Use `ionoscloud_s3_bucket` and its associated resources for full feature coverage.

## Import

A bucket can be imported using the bucket name:

```shell
terraform import ionoscloud_user_object_storage_bucket.example bucket_name
```

If the bucket is not in the default region (`de`), specify the region explicitly:

```shell
terraform import ionoscloud_user_object_storage_bucket.example region:bucket_name
```

For example:

```shell
terraform import ionoscloud_user_object_storage_bucket.example eu-central-2:my-bucket
```
