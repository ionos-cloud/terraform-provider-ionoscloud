---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IONOS CLOUD: object_storage_accesskey"
description: |-
  Lists IONOS CLOUD Object Storage Access Keys.
---

# List Resource: ionoscloud_object_storage_accesskey

-> **Note:** List Resources are supported in HashiCorp Terraform version 1.14 and later.

Lists [Object Storage Access Keys](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage/concepts/key-management) on IONOS CLOUD.

## Example Usage

``` hcl
list "ionoscloud_object_storage_accesskey" "example" {
  provider = ionoscloud
}
```

## Argument Reference

This list resource has no configuration arguments.

## Attributes Reference

Each result exposes the following attributes, matching the `ionoscloud_object_storage_accesskey` resource schema:

- `id` - The ID (UUID) of the Access Key.
- `accesskey` - Access key metadata string (92 characters).
- `canonical_user_id` - The canonical user ID valid for user-owned buckets.
- `contract_user_id` - The contract user ID valid for contract-owned buckets.
- `description` - Description of the Access Key.

> **Note:** `secretkey` is not available via the list resource as it is only returned at creation time.

> **⚠ WARNING:** `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` can be used to set a custom API URL for the Object Storage Management SDK. Setting `endpoint` or `IONOS_API_URL` does not have any effect.
