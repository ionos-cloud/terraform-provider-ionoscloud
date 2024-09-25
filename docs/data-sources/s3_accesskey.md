---
subcategory: "S3 management"
layout: "ionoscloud"
page_title: "IonosCloud : s3_accesskey"
sidebar_current: "docs-datasource-s3_accesskey"
description: |-
  Get information on a IonosCloud S3 Accesskeys
---

# ionoscloud\_s3\_accesskey

The **S3 Accesskey data source** can be used to search for and return an existing S3 Accesskeys.

## Example Usage

### By ID 
```hcl
data "ionoscloud_s3_accesskey" "example" {
  id       = <accesskey_id>
}
```

## Argument Reference

 * `id` - (Required) Id of an existing S3 accesskey that you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

- `id` - The ID (UUID) of the AccessKey.
- `description` - Description of the Access key.
- `accesskey` - Access key metadata is a string of 92 characters.
- `canonical_user_id` - The canonical user ID which is valid for user-owned buckets.
- `contract_user_id` - The contract user ID which is valid for contract-owned buckets
