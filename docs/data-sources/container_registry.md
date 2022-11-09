---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_container_registry"
sidebar_current: "docs-ionoscloud_container_registry"
description: |-
  Get information on a Container Registry
---

# ionoscloud_container_registry

⚠️ **Note:** Container Registry is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

The **Container Registry data source** can be used to search for and return an existing Container Registry.
You can provide a string for the name parameter which will be compared with provisioned Container Registry.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By Id
```hcl
data "ionoscloud_container_registry" "example" {
  id  = <registry_id>
}
```

### By Name
```hcl
data "ionoscloud_container_registry" "example" {
  name  = "container-registry-example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_container_registry" "example" {
  name          = "-example"
  partial_match = true
}
```

## Argument Reference

* `id` - (Optional) ID of the container registry you want to search for.
* `name` - (Optional) Name of an existing container registry that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

Either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the container registry.
* `name` - The name of the container registry.
* `location` 
* `garbage_collection_schedule`
    * `time`
    * `days`
* `hostname`
* `maintenance_window`
  * `time`
  * `days`
* `storage_usage`
  * `bytes`
  * `updated_at`
