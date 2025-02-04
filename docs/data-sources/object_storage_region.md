---
subcategory: "Object storage management"
layout: "ionoscloud"
page_title: "IonosCloud : object_storage_region"
sidebar_current: "docs-datasource-object_storage_region"
description: |-
  Get information on a IonosCloud Object Storage Region
---

# ionoscloud_object_storage_region

The **Object storage region data source** can be used to search for and return an existing S3 Regions.

## Example Usage

### By ID 
```hcl
data "ionoscloud_object_storage_region" "example" {
  id       = "region_id"
}
```

## Argument Reference

 * `id` - (Required) Id of an existing object storage Region that you want to search for.

## Attributes Reference

The following attributes are returned by the datasource:

- `id` - The id of the region
- `version` - The version of the region properties
- `endpoint` - The endpoint URL for the region
- `website` - The website URL for the region
- `storage_classes` - The available classes in the region
- `location` - The data center location of the region as per [Get Location](/docs/cloud/v6/#tag/Locations/operation/locationsGet). *Can't be used as `LocationConstraint` on bucket creation.*
- `capability` - The capabilities of the region
  * `iam` - Indicates if IAM policy based access is supported
  * `s3select` - Indicates if S3 Select is supported

> **⚠ WARNING:** `IONOS_API_URL_OBJECT_STORAGE_MANAGEMENT` can be used to set a custom API URL for the Object Storage Management SDK. Setting `endpoint` or `IONOS_API_URL` does not have any effect
