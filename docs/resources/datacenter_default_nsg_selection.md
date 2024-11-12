---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: Datacenter Default NSG Selection"
sidebar_current: "docs-resource-datacenter-default-nsg-selection"
description: |-
  Manages the selection of the default Network Security Group for IonosCloud datacenters.
---

# ionoscloud_datacenter_default_nsg_selection

Manages the selection of the default Network Security Group for IonosCloud datacenters.

## Example Usage

The boot device of a `ionoscloud_server`, `ionoscloud_vcpu_server` or `ionoscloud_cube_server` can be selected with this resource.
The default Network Security Group of a `ionoscloud_datacenter` can be selected with this resource.
Deleting this resource will unset the default NSG of the datacenter.

### Select an external volume
```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter Default NSG Example"
  location        = "de/fra"
}

resource "ionoscloud_nsg" "example" {
  name              = "NSG"
  description       = "NSG"
  datacenter_id     = ionoscloud_datacenter.example.id
}

resource "ionoscloud_datacenter_default_nsg_selection" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  nsg_id            = ionoscloud_nsg.example.id
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `nsg_id` - (Required)[string] The ID of a Network Security Group.
