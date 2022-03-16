---
subcategory: "User Management"
layout: "ionoscloud"
page_title: "IonosCloud : resource"
sidebar_current: "docs-datasource-resource"
description: |-
  Get information on a IonosCloud Resource
---

# ionoscloud\_resource

The **Resource data source** can be used to search for and return any existing IonosCloud resource and optionally their group associations.
You can provide a string for the resource type (datacenter,image,snapshot,ipblock) and/or resource id parameters which will be queries against available resources.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID 
```hcl
data "ionoscloud_resource" "example" {
  resource_id   = <resource_id>
}
```

### By Type
```hcl
data "ionoscloud_resource" "example" {
  resource_type = "datacenter"
}
```

### By ID & Type
```hcl
data "ionoscloud_resource" "example" {
  resource_type = "datacenter"
  resource_id   = <resource_id>
}
```

## Argument Reference

 * `resource_type` - (Optional) The specific type of resources to retrieve information about.
 * `resource_id` - (Optional) The ID of the specific resource to retrieve information about.

## Attributes Reference

 * `id` - UUID of the Resource
