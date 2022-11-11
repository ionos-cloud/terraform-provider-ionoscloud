---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_container_registry"
sidebar_current: "docs-resource-ionoscloud_container_registry"
description: |-
  Creates and manages IonosCloud Container Registry.
---

# ionoscloud_container_registry

⚠️ **Note:** Container Registry is currently in the Early Access (EA) phase. We recommend keeping usage and testing to non-production critical applications.
Please contact your sales representative or support for more information.

Manages an **Container Registry** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_container_registry" "example" {
  garbage_collection_schedule {
    days    = ["Monday", "Tuesday"]
    time    = "05:19:00+00:00"
  }
  location  = "de/fra"
  name      = "container-registry-example"
}
```

## Argument Reference

The following arguments are supported:

* `name`     - The name of the container registry. Immutable, update forces re-creation of the resource.
* `garbage_collection_schedule` - (Optional)[Map]
    * `time` - (Required)[string]
    * `days` - (Required)[list] Elements of list must have one of the values: `Saturday`, `Sunday`, `Monday`, `Tuesday`,  `Wednesday`,  `Thursday`,  `Friday` 
* `location` - (Required)[string] Immutable, update forces re-creation of the resource.


## Import

Resource Container Registry can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_container_registry.mycr {container_registry uuid}
```