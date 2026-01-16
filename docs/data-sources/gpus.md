---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_gpus"
sidebar_current: "docs-ionoscloud-datasource-gpus"
description: |-
  Get information on Ionos Cloud GPUs
---

# ionoscloud_gpus

The **GPUs data source** can be used to retrieve a list of all GPUs attached to a specific server within a datacenter.

## Example Usage

```hcl
data "ionoscloud_gpus" "example" {
  datacenter_id = "datacenter_id"
  server_id     = "server_id"
}
```

## Argument Reference

* `datacenter_id` - (Required) The ID of the datacenter.
* `server_id` - (Required) The ID of the server.

## Attributes Reference

The following attributes are returned by the datasource:

* `gpus` - A list of GPUs. Each GPU has the following attributes:
  * `id` - The id of the GPU.
  * `name` - The name of the GPU.
  * `vendor` - The vendor of the GPU.
  * `type` - The type of the GPU.
  * `model` - The model of the GPU.
