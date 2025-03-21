---
subcategory: "Container Registry"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_container_registry"
sidebar_current: "docs-resource-ionoscloud_container_registry"
description: |-
  Creates and manages IonosCloud Container Registry.
---

# ionoscloud_container_registry

Manages an [Container Registry](https://docs.ionos.com/cloud/containers/private-container-registry/overview) on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_container_registry" "example" {
  garbage_collection_schedule {
    days                = ["Monday", "Tuesday"]
    time                = "05:19:00+00:00"
  }
  location              = "de/fra"
  name                  = "container-registry-example"
  api_subnet_allow_list = ["1.2.3.4/32"]
}
```

## Argument Reference

The following arguments are supported:

* `name`     - The name of the container registry. Immutable, update forces re-creation of the resource.
* `garbage_collection_schedule` - (Optional)[Map]
    * `time` - (Required)[string]
    * `days` - (Required)[list] Elements of list must have one of the values: `Saturday`, `Sunday`, `Monday`, `Tuesday`,  `Wednesday`,  `Thursday`,  `Friday` 
* `location` - (Required)[string] Immutable, update forces re-creation of the resource.
* `api_subnet_allow_list` - (Optional)[list] The subnet CIDRs that are allowed to connect to the registry.  Specify "a.b.c.d/32" for an individual IP address. __Note__: If this list is empty or not set, there are no restrictions.
* `features` - (Optional)[Map]
    * `vulnerability_scanning` - (Optional)[bool] Enables or disables the Vulnerability Scanning feature for the Container Registry. To disable this feature, set the attribute to false when creating the CR resource.
  
> **⚠ WARNING** `Container Registry Vulnerability Scanning` is a paid feature which is enabled by default, and cannot be turned off after activation. To disable this feature for a Container Registry, ensure `vulnerability_scanning` is set to false on resource creation.

## Import

Resource Container Registry can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_container_registry.mycr container_registry uuid
```