---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IONOS CLOUD: object_storage_accesskey"
sidebar_current: "docs-resource-object_storage_accesskey"
description: |-
  Creates and manages IONOS CLOUD Object Storage Accesskeys.
---

# ionoscloud_object_storage_accesskey

Manages an [Object Storage Accesskey](https://docs.ionos.com/cloud/storage-and-backup/ionos-object-storage/concepts/key-management) on IONOS CLOUD.

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

In Terraform v1.12.0 and later, the [`import` block](https://developer.hashicorp.com/terraform/language/import) can be used with the `identity` attribute. For example:

```hcl
import {
  to = ionoscloud_object_storage_accesskey.example
  identity = {
    id = "objectStorageAccesskeyid"
  }
}

resource "ionoscloud_object_storage_accesskey" "example" {
  ### Configuration omitted for brevity ###
}
```

### Identity Schema

#### Required

* `id` (String) The ID (UUID) of the AccessKey.

---

An object storage accesskey resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_object_storage_accesskey.demo objectStorageAccesskeyid
```

This can be helpful when you want to import Object Storage Accesskeys which you have already created manually or using other means, outside of terraform.

## Query (List Resource)

Object Storage Access Keys can be listed using `terraform query` (requires Terraform 1.14+). List blocks must be placed in a dedicated `tfquery.hcl` file.

```hcl
list "ionoscloud_object_storage_accesskey" "all" {
  provider         = ionoscloud
  include_resource = true
}
```

Filter by description:

```hcl
list "ionoscloud_object_storage_accesskey" "filtered" {
  provider         = ionoscloud
  include_resource = true
  config {
    filters = [{
      field_name  = "description"
      field_value = "my-key"
    }]
  }
}
```

See the [ionoscloud_object_storage_accesskey list resource](../list-resources/object_storage_accesskey.md) documentation for the full filter reference.
