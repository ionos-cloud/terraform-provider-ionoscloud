---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_object"
sidebar_current: "docs-ionoscloud-datasource-s3_object"
description: |-
  Get information about  IONOS Object Storage Objects.
---

# ionoscloud_s3_object

The **Object data source** can be used to search for and return existing objects.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

```hcl

data "ionoscloud_s3_object" "example" {
  bucket = "example"
  key = "object"
}

 ```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object is stored.
- `key` - (Required)[string] The name of the object.
- `version_id` - (Optional)[string] The version of the object.
- `range` - (Optional)[string] Downloads the specified range bytes of an object. For more information about the HTTP Range header


## Attributes Reference

The following attributes are returned by the datasource:

- `body` - The content of the object.
- `cache_control` - The caching behavior along the request/reply chain.
- `content_length` - The size of the object in bytes.
- `content_disposition` - Presentational information for the object.
- `content_encoding` - The content encodings applied to the object.
- `content_language` - The natural language of the intended audience for the object.
- `content_type` - The MIME type describing the format of the contents.
- `expires` - The date and time at which the object is no longer cacheable.
- `server_side_encryption` - The server-side encryption algorithm used when storing this object.
- `storage_class` - The storage class of the object.
- `website_redirect` - If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL.
- `server_side_encryption_customer_algorithm` - The algorithm to use for encrypting the object (e.g., AES256).
- `server_side_encryption_customer_key` - The 256-bit, base64-encoded encryption key to encrypt and decrypt your data. This attribute is sensitive.
- `server_side_encryption_customer_key_md5` - The 128-bit MD5 digest of the encryption key.
- `server_side_encryption_context` - The encryption context to use for object encryption. This attribute is sensitive.
- `request_payer` - Confirms that the requester knows that they will be charged for the request.
- `object_lock_mode` - The object lock mode, which can be either GOVERNANCE or COMPLIANCE.
- `object_lock_retain_until_date` - The date until which the object will remain locked.
- `object_lock_legal_hold` - The legal hold status of the object, which can be either ON or OFF.
- `etag` - An entity tag (ETag) assigned by a web server to a specific version of a resource.
- `tags` - The tag-set for the object, represented as a map of string key-value pairs.
- `metadata` - A map of metadata stored with the object.
- `version_id` - The version of the object. This attribute is optional.