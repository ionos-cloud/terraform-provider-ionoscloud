---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_container_registry_locations"
sidebar_current: "docs-ionoscloud_container_registry_locations"
description: |-
Get list of Container Registry Locations
---

# ionoscloud_container_registry_locations

The **Container Registry Token data source** can be used to get a list of Container Registry Locations

⚠️ **Note:** Container Registry is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

## Example Usage

```hcl
data "ionoscloud_container_registry_locations" "example" {
}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `locations` - list of container registry locations