---
subcategory: "Container Registry Token"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_container_registry_token"
sidebar_current: "docs-resource-ionoscloud_container_registry_token"
description: |-
Creates and manages IonosCloud Container Registry Token.
---

# ionoscloud_container_registry_token

Manages an **Container Registry Token** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_container_registry" "example" {
    garbage_collection_schedule {
        days			 = ["Monday", "Tuesday"]
        time             = "10:00:00"
    }
  location               = "de/txl"
    maintenance_window {
        days			 = ["Sunday"]
        time             = "09:00:00"
    }
    name		         = "container-registry-example"
}

resource "ionoscloud_container_registry_token" "example" {
    credentials {
        username	   = "username"
        password       = "password"
    }
    expiry_date        = "2023-01-13T16:27:42Z"
    name			   = "container-registry-token-example"
    scopes  {
        actions		   = ["push"]
        name           = "Scope1"
        type           = "repository"
    }
    status	           = "enabled"
    registry_id        = ionoscloud_container_registry.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required)[string] The name of the container registry token.
* `expiry-date` - (Optional)[string] The value must be supplied as ISO 8601 timestamp
* `scopes` - (Optional)[map]
  * `actions` - (Required)[string]
  * `name` - (Required)[string]
  * `type` - (Required)[string]
* `status` - (Optional)[string] Must have on of the values: `enabled`, `disabled`

## Import

Resource Container Registry Token can be imported using the `container registry id` and `resource id`, e.g.

```shell
terraform import ionoscloud_container_registry_token.mycrtoken {container_registry uuid}/{container_registry_token uuid}
```