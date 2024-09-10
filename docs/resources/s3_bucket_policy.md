---
subcategory: "S3"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket_policy"
sidebar_current: "docs-resource-s3_bucket_policy"
description: |-
  Creates and manages IonosCloud S3 Buckets policies.
---

# ionoscloud_s3_bucket_policy

Manages **S3 Buckets policies** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_s3_bucket" "example" {
  name = "example"
}

resource "ionoscloud_s3_bucket_policy" "example" {
  bucket = ionoscloud_s3_bucket.example.name
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid = "Delegate certain actions to another user"
        Action = [
          "s3:ListBucket",
          "s3:PutObject",
          "s3:GetObject"
        ]
        Effect = "Allow"
        Resource = [
          "arn:aws:s3:::example",
          "arn:aws:s3:::example/*"
        ]
        Condition = {
          IpAddress = [
            "123.123.123.123/32"
          ]
        }
        Principal = [
          "arn:aws:iam:::user/31000000:9acd8251-2857-410e-b1fd-ca86462bdcec"
        ]
      }
    ]
  })
}

```

## Argument Reference

The following arguments are supported:

- `bucket` - (Required)[string] The name of the bucket where the object will be stored.
- `policy` - (Required)[string] The policy document. This is a JSON formatted string.

## Import

Resource Policy can be imported using the `bucket name`

```shell
terraform import ionoscloud_s3_bucket_policy.example example
```