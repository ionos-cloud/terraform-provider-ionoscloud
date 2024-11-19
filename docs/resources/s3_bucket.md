---
subcategory: "Object Storage"
layout: "ionoscloud"
page_title: "IonosCloud: s3_bucket"
sidebar_current: "docs-resource-s3_bucket"
description: |-
  Creates and manages IONOS Object Storage Buckets.
---

# ionoscloud_s3_bucket

Manages **IONOS Object Storage Buckets** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_s3_bucket" "example" {
  name = "example"
  region = "eu-central-3"
  object_lock_enabled = true
  force_destroy = true
  
  tags = {
    key1 = "value1"
    key2 = "value2"
  }

  timeouts {
    create = "10m"
    delete = "10m"
  }
}

```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The bucket name. [ 3 .. 63 ] characters
- `region` - (Optional)[string] Specifies the Region where the bucket will be created. Please refer to the list of available regions
- `object_lock_enabled` - (Optional)[bool] The object lock configuration status of the bucket. Must be `true` or `false`.
- `tags` - (Optional) A mapping of tags to assign to the bucket.
- `timeouts` - (Optional) Timeouts for this resource.
  - `create` - (Optional)[string] Time to wait for the bucket to be created. Default is `10m`.
  - `delete` - (Optional)[string] Time to wait for the bucket to be deleted. Default is `10m`.
- `force_destroy` - (Optional)[bool] Default is `false`.By setting force_destroy to true, you instruct Terraform to delete the bucket and all its contents during the terraform destroy process. This is particularly useful when dealing with buckets that contain objects, as it allows for automatic cleanup without requiring the manual deletion of objects beforehand. If force_destroy is not set or is set to false, Terraform will refuse to delete a bucket that still contains objects. You must manually empty the bucket before Terraform can remove it.There is a significant risk of accidental data loss when using this attribute, as it irreversibly deletes all contents of the bucket. It's crucial to ensure that the bucket does not contain critical data before using force_destroy.

## Attributes Reference

- `id` - (Computed) Name of the bucket

⚠️ **Note:** The name must be unique across all IONOS accounts in all IONOS Object Storage regions. The name should adhere to the following [restrictions](https://docs.ionos.com/cloud/storage-and-backup/s3-object-storage/concepts/buckets#naming-conventions).

## Import

Resource Bucket can be imported using the `bucket name`

```shell
terraform import ionoscloud_s3_bucket.example example
```