---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_website_configuration"
sidebar_current: "docs-resource-s3_bucket_website_configuration"
description: |-
  Manages Buckets website configuration on IonosCloud.
---

# ionoscloud_s3_bucket_website_configuration

Manages Website Configuration for Buckets on IonosCloud.

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_website_configuration" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }

  routing_rule {
    condition {
      key_prefix_equals = "docs/"
    }
    redirect {
      replace_key_prefix_with = "documents/"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `index_document` - (Optional) Container for the Suffix element.
  - `suffix` - (Required) A suffix that is appended to a request that is for a directory on the website endpoint (for example, if the suffix is index.html and you make a request to samplebucket/images/ the data that is returned will be for the object with the key name images/index.html) The suffix must not be empty and must not include a slash character. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests.
- `error_document` - (Optional) The object key name to use when a 4XX class error occurs. Replacement must be made for object keys containing special characters (such as carriage returns) when using XML requests
    - `key` - (Required) The object key
- `redirect_all_requests_to` - (Optional) Container for redirect information. You can redirect requests to another host, to another page, or with another protocol. In the event of an error, you can can specify a different error code to return.
  - `host_name` - (Optional) Name of the host where requests will be redirected.
  - `protocol` - (Optional) Protocol to use (http, https).
- `routing_rule` - (Optional) A container for describing a condition that must be met for the specified redirect to apply.
  - `condition` - (Required) A container for describing a condition that must be met for the specified redirect to apply.
    - `http_error_code_returned_equals` - (Optional) The HTTP error code when the redirect is applied. In the event of an error, if the error code equals this value, then the specified redirect is applied.
    - `key_prefix_equals` - (Optional) The object key name prefix when the redirect is applied. For example, to redirect requests for ExamplePage.html, the key prefix will be ExamplePage.html. To redirect request for all pages with the prefix example, the key prefix will be /example.
  - `redirect` - (Required) Container for the redirect information.
    - `host_name` - (Optional) The host name to use in the redirect request.
    - `http_redirect_code` - (Optional) The HTTP redirect code to use on the response. Not required if one of the siblings is present.
    - `protocol` - (Optional) Protocol to use (http, https).
    - `replace_key_prefix_with` - (Optional) The object key to be used in the redirect request. For example, redirect request to error.html, the replace key prefix will be /error.html. Not required if one of the siblings is present.
    - `replace_key_with` - (Optional) The specific object key to use in the redirect request. For example, redirect request for error.html, the replace key will be /error.html. Not required if one of the siblings is present.
    - `http_redirect_code` - (Optional) The HTTP redirect code to use on the response. Not required if one of the siblings is present.

## Import

IONOS Object Storage Bucket website configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_website_configuration.example example
```