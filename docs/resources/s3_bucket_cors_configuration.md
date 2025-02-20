---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_cors_configuration"
sidebar_current: "docs-resource-s3_bucket_cors_configuration"
description: |-
  Manages Buckets cors_configuration on IonosCloud.
---

# ionoscloud_s3_bucket_cors_configuration

Manages Object Lock Configuration for Buckets on IonosCloud.

⚠️ **Note:** The Terraform provider **only supports contract-owned buckets. User-owned buckets are not supported,** and there are no plans to introduce support for them. As a result, **user-owned buckets cannot be created, updated, deleted, read, or imported** using this provider.

## Example Usage

```hcl
resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_cors_configuration" "test" {
  bucket = ionoscloud_s3_bucket.example.name
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
    id = 1234
  }
}
```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `cors_rule` - (Required)[block] A block of cors_rule as defined below.
  - `allowed_headers` - (Optional)[list] Specifies which headers are allowed in a preflight OPTIONS request through the Access-Control-Request-Headers header
  - `allowed_methods` - (Required)[list] An HTTP method that you allow the origin to execute. Valid values are GET, PUT, HEAD, POST, DELETE.
  - `allowed_origins` - (Required)[list] Specifies which origins are allowed to make requests to the resource.
  - `expose_headers` - (Optional)[list] Specifies which headers are exposed to the browser.
  - `max_age_seconds` - (Optional)[int] Specifies how long the results of a pre-flight request can be cached in seconds.
  - `id` - (Optional)[int] Container for the Contract Number of the owner

Days and years are mutually exclusive. You can only specify one of them.
## Import

IONOS Object Storage Bucket cors configuration can be imported using the `bucket` name.

```shell
terraform import ionoscloud_s3_bucket_cors_configuration.example example
```