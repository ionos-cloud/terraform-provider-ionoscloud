---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_object_copy"
sidebar_current: "docs-resource-s3_object_copy"
description: |-
  Creates a copy of an object that is already stored in IONOS S3 Object Storage.
---

# ionoscloud_s3_object_copy

Creates a copy of an object that is already stored in IONOS S3 Object Storage.

## Example Usage

```hcl

resource "ionoscloud_s3_bucket" "source" {
  name = "source"
}

resource "ionoscloud_s3_bucket" "target" {
  name = "target"
}

resource "ionoscloud_s3_object" "source" {
  bucket  = ionoscloud_s3_bucket.source.name
  key     = "source_object"
  content = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  content_type = "application/octet-stream"
}

resource "ionoscloud_s3_object_copy" "example" {
  bucket = ionoscloud_s3_bucket.target.name
  key    = "example"
  source = "${ionoscloud_s3_bucket.source.name}/${ionoscloud_s3_object.source.key}"
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored. Must be between 3 and 63 characters.
- `key`  - (Required)[string] The key of the object. Must be at least 1 character long.
- `source` - (Optional)[string] The source of the object to be copied
- `copy_source_if_match` - (Optional)[string] Copies the object if its entity tag (ETag) matches the specified tag.
- `copy_source_if_none_match` - (Optional)[string] Copies the object if its entity tag (ETag) is different than the specified ETag.
- `copy_source_if_modified_since` - (Optional)[string] Copies the object if it has been modified since the specified time.
- `copy_source_if_unmodified_since` - (Optional)[string] Copies the object if it hasn't been modified since the specified time.
- `cache_control` - (Optional)[string] Specifies caching behavior along the request/reply chain.
- `content_disposition` - (Optional)[string] Specifies presentational information for the object.
- `content_encoding` - (Optional)[string] Specifies what content encodings have been applied to the object.
- `content_language` - (Optional)[string] The natural language or languages of the intended audience for the object.
- `content_type` - (Optional)[string] A standard MIME type describing the format of the contents.
- `expires` - (Optional)[string] The date and time at which the object is no longer cacheable.
- `server_side_encryption` - (Optional)[string] The server-side encryption algorithm used when storing this object in IONOS S3 Object Storage. Valid value is AES256.
- `storage_class` - (Optional)[string] The storage class of the object. Valid value is STANDARD. Default is STANDARD.
- `website_redirect` - (Optional)[string] Redirects requests for this object to another object in the same bucket or to an external URL.
- `server_side_encryption_customer_algorithm` - (Optional)[string] Specifies the algorithm to use for encrypting the object. Valid value is AES256.
- `server_side_encryption_customer_key` - (Optional)[string] Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data.
- `server_side_encryption_customer_key_md5` - (Optional)[string] Specifies the 128-bit MD5 digest of the encryption key.
- `server_side_encryption_context` - (Optional)[string] Specifies the IONOS S3 Object Storage Encryption Context for object encryption.
- `source_customer_algorithm` - (Optional)[string] Specifies the algorithm used for source object encryption. Valid value is AES256.
- `source_customer_key` - (Optional)[string] Specifies the 256-bit, base64-encoded encryption key for source object encryption.
- `source_customer_key_md5` - (Optional)[string] Specifies the 128-bit MD5 digest of the encryption key for source object encryption.
- `object_lock_mode` - (Optional)[string] The object lock mode that you want to apply to the object. Valid values are `GOVERNANCE` and `COMPLIANCE`.
- `object_lock_retain_until_date` - (Optional)[string] The date and time when the object lock retention expires.Must be in RFC3999 format
- `object_lock_legal_hold` - (Optional)[string] Indicates whether a legal hold is in effect for the object. Valid values are `ON` and `OFF`.
- `etag` - (Computed)[string] An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.
- `last_modified` - (Computed)[string] The date and time at which the object was last modified.
- `metadata_directive` - (Optional)[string] Specifies whether the metadata is copied from the source object or replaced with metadata provided in the request. Valid values are `COPY` and `REPLACE`.
- `metadata` - (Optional)[map] A map of metadata to store with the object in IONOS S3 Object Storage. Metadata keys must be lowercase alphanumeric characters.
- `tagging_directive` - (Optional)[string] Specifies whether the object tag-set is copied from the source object or replaced with tag-set provided in the request. Valid values are `COPY` and `REPLACE`.
- `tags` - (Optional)[map] The tag-set for the object.
- `version_id` - (Computed)[string] The version of the object.
- `force_destroy` - (Optional)[bool] If true, the object will be destroyed if versioning is enabled then all versions will be destroyed. Default is `false`.