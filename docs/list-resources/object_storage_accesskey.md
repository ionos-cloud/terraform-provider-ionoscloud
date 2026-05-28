---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IONOS CLOUD: object_storage_accesskey"
description: |-
  Lists IONOS CLOUD Object Storage Access Keys.
---

# List Resource: ionoscloud_object_storage_accesskey

⚠️ **Note:** List Resources require HashiCorp Terraform version 1.14 or later and are queried using `terraform query`.

Lists [Object Storage Access Keys](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage/concepts/key-management) on IONOS CLOUD.

## Example Usage

⚠️ **Note:** `list` blocks must be placed in a dedicated `tfquery.hcl` file, separate from your main Terraform configuration.

### List access keys

```hcl
list "ionoscloud_object_storage_accesskey" "example" {
  provider = ionoscloud
  include_resource = true
}
```

### Filter access keys by description

```hcl
list "ionoscloud_object_storage_accesskey" "filtered" {
  provider = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "description"
      field_value = "my-key"
    }]
  }
}
```

### Generate resource configuration from existing access keys

Use `terraform query` with `-generate-config-out` to produce ready-to-use `ionoscloud_object_storage_accesskey` resource blocks for all existing access keys:

```shell
terraform query -generate-config-out=imported.tf
```

Terraform will write an `ionoscloud_object_storage_accesskey` resource block for each discovered access key into `imported.tf`, which can then be used directly in your configuration.

## Argument Reference

The `config` block supports the following arguments:

- `filters` - (Optional) List of filters to apply. All filters must match (AND logic). Each filter supports:
  - `field_name` - (Required) The field to filter on. Supported values: `id`, `description`, `accesskey`.
  - `field_value` - (Required) The exact value to match against.

## Attributes Reference

Each result exposes the following attributes, matching the `ionoscloud_object_storage_accesskey` resource schema:

- `id` - The ID (UUID) of the Access Key.
- `accesskey` - Access key metadata string (92 characters).
- `canonical_user_id` - The canonical user ID valid for user-owned buckets.
- `contract_user_id` - The contract user ID valid for contract-owned buckets.
- `description` - Description of the Access Key.

> **Note:** `secretkey` is not available via the list resource as it is only returned at creation time.

> **⚠ WARNING:** `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` can be used to set a custom API URL for the Object Storage Management SDK. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
