---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_objects"
sidebar_current: "docs-ionoscloud-datasource-s3_objects"
description: |-
  Get information about IONOS Object Storage Objects.
---

# ionoscloud_s3_objects

The **Objects data source** can be used to search for and return existing objects.

## Example Usage

```hcl

data "ionoscloud_s3_objects" "example" {
  bucket = "example"
  prefix    = "prefix1/"
  delimiter = "/"
  max_keys  = 100
  fetch_owner = true
}

 ```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the objects are stored.
- `encoding_type` - (Optional)[string] Specifies the encoding method used to encode the object keys in the response. Default is url.
- `prefix` - (Optional)[string] Limits the response to keys that begin with the specified prefix.
- `delimiter` - (Optional)[string] A character used to group keys.
- `max_keys` - (Optional)[int] Sets the maximum number of keys returned in the response body.Default is 1000.
- `fetch_owner` - (Optional)[bool] If set to true, the response includes the owner field in the metadata.
- `start_after` - (Optional)[string] Specifies the key to start after when listing objects in a bucket.

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.


## Attributes Reference

The following attributes are returned by the datasource:

- `keys` - A list of objects in the bucket.
- `common_prefixes` - A list of keys that act as a delimiter for grouping keys.
- `owner` - The owner of the bucket.