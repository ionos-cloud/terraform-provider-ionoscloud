---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_container_registry_locations"
sidebar_current: "docs-ionoscloud_container_registry_locations"
description: |-
  Get list of Container Registry Locations
---

# ionoscloud_container_registry_locations

The **Container Registry Locations data source** can be used to get a list of Container Registry Locations

## Example Usage

```hcl
data "ionoscloud_container_registry_locations" "example" {
}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `locations` - list of container registry locations