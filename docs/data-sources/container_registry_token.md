---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_container_registry_token"
sidebar_current: "docs-ionoscloud_container_registry_token"
description: |-
Get information on a Container Registry Token
---

# ionoscloud_container_registry_token

⚠️ **Note:** Container Registry is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

The **Container Registry Token data source** can be used to search for and return an existing Container Registry Token.
You can provide a string for the name parameter which will be compared with provisioned Container Registry Token.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By Id
```hcl
data "ionoscloud_container_registry_token" "example" {
  registry_id = ionoscloud_container_registry.example.id
  id			= <token_id>
}
```

### By Name
```hcl
data "ionoscloud_container_registry_token" "example" {
  registry_id   = ionoscloud_container_registry.example.id
  name			= "container-registry-token-example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_container_registry_token" "example" {
  registry_id   = ionoscloud_container_registry.example.id
  name			= "-example"
  partial_match = true
}
```

## Argument Reference

* `registry_id` - (Required) Registry's UUID.
* `id` - (Optional) ID of the container registry token you want to search for.
* `name` - (Optional) Name of an existing container registry token that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true.
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

`registry_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the container registry token.
* `name` - The name of the container registry token.
* `credentials` 
    * `username`
* `expiry-date`
* `scopes`
  * `actions`
  * `name`
  * `type`
* `status` 
