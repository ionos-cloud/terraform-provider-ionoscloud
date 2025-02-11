---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IonosCloud : object_storage_accesskey"
sidebar_current: "docs-datasource-object_storage_accesskey"
description: |-
  Get information on a IonosCloud Object storage Accesskey
---

# ionoscloud_object_storage_accesskey

The **Object Storage Accesskey data source** can be used to search for and return an existing Object Storage Accesskeys.

## Example Usage

### By ID 
```hcl
data "ionoscloud_object_storage_accesskey" "example" {
  id       = "accesskey_id"
}
```

## Argument Reference

 * `id` - (Optional) Id of an existing object storage accesskey that you want to search for.
 * `accesskey` - (Optional) Access key metadata is a string of 92 characters.
 * `description` - (Optional) Description of the Access key.

## Attributes Reference

The following attributes are returned by the datasource:

- `id` - The ID (UUID) of the AccessKey.
- `description` - Description of the Access key.
- `accesskey` - Access key metadata is a string of 92 characters.
- `canonical_user_id` - The canonical user ID which is valid for user-owned buckets.
- `contract_user_id` - The contract user ID which is valid for contract-owned buckets

> **⚠ WARNING:** `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` can be used to set a custom API URL for the Object Storage Management SDK. Setting `endpoint` or `IONOS_API_URL` does not have any effect
