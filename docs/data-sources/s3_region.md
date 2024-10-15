---
subcategory: "S3 management"
layout: "ionoscloud"
page_title: "IonosCloud : s3_region"
sidebar_current: "docs-datasource-s3_region"
description: |-
  Get information on a IonosCloud S3 Region
---

# ionoscloud\_s3\_region

The **S3 region data source** can be used to search for and return an existing S3 Regions.

## Example Usage

### By ID 
```hcl
data "ionoscloud_s3_region" "example" {
  id       = <region_id>
}
```

## Argument Reference

 * `id` - (Required) Id of an existing S3 Region that you want to search for.

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
