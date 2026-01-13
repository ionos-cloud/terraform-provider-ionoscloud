---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_gpu"
sidebar_current: "docs-ionoscloud-datasource-gpu"
description: |-
  Get information on a Ionos Cloud GPU
---

# ionoscloud_gpu

The **GPU data source** can be used to search for and return an existing GPU by either its ID or name.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID
```hcl
data "ionoscloud_gpu" "example" {
  datacenter_id = "datacenter_id"
  server_id     = "server_id"
  id            = "gpu_id"
}
```

### By Name
```hcl
data "ionoscloud_gpu" "example" {
  datacenter_id = "datacenter_id"
  server_id     = "server_id"
  name          = "GPU Name"
}
```

## Argument Reference

* `datacenter_id` - (Required) The ID of the datacenter.
* `server_id` - (Required) The ID of the server.
* `name` - (Optional) Name of the GPU.
* `id` - (Optional) ID of the GPU.

`datacenter_id` and `server_id` are required. Either `name` or `id` must be provided. If both `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The id of the GPU.
* `name` - The name of the GPU.
* `vendor` - The vendor of the GPU.
* `type` - The type of the GPU.
* `model` - The model of the GPU.
