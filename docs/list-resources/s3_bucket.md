---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IONOS CLOUD: s3_bucket"
description: |-
  Lists IONOS Object Storage Buckets.
---

# List Resource: ionoscloud_s3_bucket

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [IONOS Object Storage Buckets](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `tfquery.hcl` file, separate from your main Terraform configuration.

### List buckets

```hcl
list "ionoscloud_s3_bucket" "example" {
  provider = ionoscloud
  include_resource = true
}
```

### Filter buckets by region

```hcl
list "ionoscloud_s3_bucket" "eu_central" {
  provider = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "region"
      field_value = "eu-central-3"
    }]
  }
}
```

### Generate resource configuration from existing buckets

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_s3_bucket` resource blocks for all existing buckets:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_s3_bucket` resource block for each discovered bucket into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `name`, `region`.
  - `field_value` - (Required) The exact value to match against.

## Attributes Reference

Each result exposes the following attributes, matching the `ionoscloud_s3_bucket` resource schema:

- `id` - Name of the bucket (same as `name`).
- `name` - The name of the bucket.
- `region` - The region where the bucket is located.
- `object_lock_enabled` - Whether object lock is enabled for the bucket.
- `tags` - A mapping of tags assigned to the bucket.

> **Note:** `force_destroy` is not available via the list resource as it is a local Terraform-only attribute and not returned by the API.
