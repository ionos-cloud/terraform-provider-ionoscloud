---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_container_registry"
sidebar_current: "docs-resource-ionoscloud_container_registry"
description: |-
Creates and manages IonosCloud Container Registry.
---

# ionoscloud_container_registry

Manages an **Container Registry** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_container_registry" "example" {
    garbage_collection_schedule {
        days			 = ["Monday", "Tuesday"]
        time             = "10:00:00"
    }
  location               = "de/fra"
    name		         = "container-registry-example"
}

```

## Argument Reference

The following arguments are supported:

* `name` - The name of the container registry.
* `garbage_collection_schedule` - (Optional)[Map]
    * `time` - (Required)[string]
    * `days` - (Required)[list] Elements of list must have one of the values: `Saturday`, `Sunday`, `Monday`, `Tuesday`,  `Wednesday`,  `Thursday`,  `Friday` 
* `location` - (Required)[string]
    * `time` - (Required)[string]
    * `days` - (Required)[list] Elements of list must have one of the values: `Saturday`, `Sunday`, `Monday`, `Tuesday`,  `Wednesday`,  `Thursday`,  `Friday`


## Import

Resource Container Registry can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_container_registry.mycr {container_registry uuid}
```