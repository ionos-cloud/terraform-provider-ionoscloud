---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IonosCloud: object_storage_accesskey"
sidebar_current: "docs-resource-object_storage_accesskey"
description: |-
  Creates and manages IonosCloud Object Storage Accesskeys.
---

# ionoscloud_object_storage_accesskey

Manages an [Object Storage Accesskey](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage/concepts/key-management) on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_object_storage_accesskey" "example" {
    description = "my description"
}
```

## Argument Reference

The following arguments are supported:

- `description` - (Optional)[string] Description of the Access key.
- `id` - (Computed)  The ID (UUID) of the AccessKey.
- `accesskey` - (Computed)  Access key metadata is a string of 92 characters.
- `secretkey` - (Computed)  The secret key of the Access key.
- `canonical_user_id` - (Computed)  The canonical user ID which is valid for user-owned buckets.
- `contract_user_id` - (Computed)  The contract user ID which is valid for contract-owned buckets
- `timeouts` - (Optional) Timeouts for this resource.
  - `create` - (Optional)[string] Time to wait for the bucket to be created. Default is `10m`.
  - `delete` - (Optional)[string] Time to wait for the bucket to be deleted. Default is `10m`.

> **⚠ WARNING:** `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` can be used to set a custom API URL for the Object Storage Management SDK. Setting `endpoint` or `IONOS_API_URL` does not have any effect.

## Import

An object storage accesskey resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_object_storage_accesskey.demo objectStorageAccesskeyid
```

This can be helpful when you want to import Object Storage Accesskeys which you have already created manually or using other means, outside of terraform.
