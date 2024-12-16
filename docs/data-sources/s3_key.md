---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: s3_key"
sidebar_current: "docs-resource-s3-key"
description: |-
  Get Information on a IonosCloud Object Storage key
---

# ionoscloud_s3_key

The **IONOS Object Storage key data source** can be used to search for and return an existing IONOS Object Storage key.
You can provide a string id which will be compared with provisioned IONOS Object Storage keys.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl
data "ionoscloud_s3_key" "example" {
  id         = "key_id"
  user_id    = "user-uuid"
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the IONOS Object Storage Key.
- `id` - (Required) ID of the IONOS Object Storage key you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the IONOS Object Storage key
* `active` - The state of the IONOS Object Storage key
* `user_id` - The ID of the user that owns the key
* `secret_key` - (Computed)The IONOS Object Storage Secret key.