---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_object"
sidebar_current: "docs-resource-s3_object"
description: |-
  Creates and manages IONOS Object Storage Objects.
---

# ionoscloud_s3_object

Manages **IONOS Object Storage Objects** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_s3_bucket" "example" {
  name = "example"
  object_lock_enabled = true
}

resource "ionoscloud_s3_object" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  key = "object"
  content = "body"
  content_type        = "text/plain"
  cache_control       = "no-cache"
  content_disposition = "attachment"
  content_encoding    = "identity"
  content_language    = "en-GB"
  expires			 = "2024-10-07T12:34:56Z"
  website_redirect = "https://www.ionos.com"
  server_side_encryption = "AES256"
  
  object_lock_mode = "GOVERNANCE"
  object_lock_retain_until_date = "2024-10-07T12:34:56Z"
  object_lock_legal_hold = "ON"

  tags = {
    tk = "tv"
  }

  metadata = {
    "mk" = "mv"
  }
  
  force_destroy = true
}

// Upload from file
resource "ionoscloud_s3_object" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  key = "file-object"
  source = "path/to/file"
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored. Must be between 3 and 63 characters.
- `key`  - (Required)[string] The key of the object. Must be at least 1 character long.
- `source` - (Optional)[string] The path to the file to upload.
- `content` - (Optional)[string] Inline content of the object.
- `cache_control` - (Optional)[string] Specifies caching behavior along the request/reply chain.
- `content_disposition` - (Optional)[string] Specifies presentational information for the object.
- `content_encoding` - (Optional)[string] Specifies what content encodings have been applied to the object.
- `content_language` - (Optional)[string] The natural language or languages of the intended audience for the object.
- `content_type` - (Optional)[string] A standard MIME type describing the format of the contents.
- `expires` - (Optional)[string] The date and time at which the object is no longer cacheable.
- `server_side_encryption` - (Optional)[string] The server-side encryption algorithm used when storing this object in IONOS Object Storage. Valid value is AES256.
- `storage_class` - (Optional)[string] The storage class of the object. Valid value is STANDARD. Default is STANDARD.
- `website_redirect` - (Optional)[string] Redirects requests for this object to another object in the same bucket or to an external URL.
- `server_side_encryption_customer_algorithm` - (Optional)[string] Specifies the algorithm to use for encrypting the object. Valid value is AES256.
- `server_side_encryption_customer_key` - (Optional)[string] Specifies the 256-bit, base64-encoded encryption key to use to encrypt and decrypt your data.
- `server_side_encryption_customer_key_md5` - (Optional)[string] Specifies the 128-bit MD5 digest of the encryption key.
- `server_side_encryption_context` - (Optional)[string] Specifies the IONOS Object Storage Encryption Context for object encryption.
- `request_payer` - (Optional)[string] Confirms that the requester knows that they will be charged for the request.
- `object_lock_mode` - (Optional)[string] The object lock mode that you want to apply to the object. Valid values are `GOVERNANCE` and `COMPLIANCE`.
- `object_lock_retain_until_date` - (Optional)[string] The date and time when the object lock retention expires.Must be in RFC3999 format
- `object_lock_legal_hold` - (Optional)[string] Indicates whether a legal hold is in effect for the object. Valid values are `ON` and `OFF`.
- `etag` - (Computed)[string] An entity tag (ETag) is an opaque identifier assigned by a web server to a specific version of a resource found at a URL.
- `metadata` - (Optional)[map] A map of metadata to store with the object in IONOS Object Storage. Metadata keys must be lowercase alphanumeric characters.
- `tags` - (Optional)[map] The tag-set for the object.
- `version_id` - (Computed)[string] The version of the object.
- `mfa` - (Optional) [string]The concatenation of the authentication device's serial number, a space, and the value displayed on your authentication device.
- `force_destroy` - (Optional)[bool] If true, the object will be destroyed if versioning is enabled then all versions will be destroyed. Default is `false`.

## Import

Resource Object can be imported using the `bucket name` and `object key`

```shell
terraform import ionoscloud_s3_object.example example/object
```