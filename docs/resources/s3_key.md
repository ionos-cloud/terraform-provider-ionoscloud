---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud: s3_key"
sidebar_current: "docs-resource-s3-key"
description: |-
  Creates and manages IONOS Object Storage keys.
---

# ionoscloud_s3_key

Manages an **IONOS Object Storage Key** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_user" "example" {
    first_name              = "example"
    last_name               = "example"
    email                   = "unique@email.com"
    password                = "abc123-321CBA"
    administrator           = false
    force_sec_auth          = false
    active                  = true
}

resource "ionoscloud_s3_key" "example" {
    user_id                 = ionoscloud_user.example.id
    active                  = true
}
```

## Argument Reference

The following arguments are supported:

- `user_id` - (Required)[string] The UUID of the user owning the IONOS Object Storage Key.
- `active` - (Optional)[boolean] Whether the IONOS Object Storage is active / enabled or not - Please keep in mind this is only required on create. Default value in true
- `secret_key` - (Computed)  The IONOS Object Storage Secret key.

## Import

An IONOS Object Storage Unit resource can be imported using its user id as well as its `resource id`, e.g.

```shell
terraform import ionoscloud_s3_key.demo userid/s3Keyid
```

This can be helpful when you want to import IONOS Object Storage Keys which you have already created manually or using other means, outside of terraform.

⚠️ **Note:** Due to eventual consistency in the IONOS Cloud API, creating an S3 key immediately after assigning S3 privileges to a user may result in the error: "The user needs to be part of a group that has ACCESS_S3_OBJECT_STORAGE privilege".
To mitigate this, set the environment variable `IONOS_S3_KEY_CREATION_RETRY=true` to enable automatic retries during S3 key creation. The provider will retry the creation until it succeeds or the configured timeout is reached.
This is a temporary workaround until the API consistency issue is resolved.
