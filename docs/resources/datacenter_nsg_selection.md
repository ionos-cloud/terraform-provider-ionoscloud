---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: Datacenter NSG Selection"
sidebar_current: "docs-resource-datacenter-nsg-selection"
description: |-
  Links a Network Security Group to a IonosCloud datacenter.
---

# ionoscloud_datacenter_nsg_selection

Links a Network Security Group to a IonosCloud datacenter. The datacenter can only have one linked NSG. To set a new NSG for the datacenter, the current one will be unlinked.

## Example Usage

A Network Security Group can be linked to a `ionoscloud_datacenter` with this resource.
Deleting the resource will unlink the NSG from the datacenter.

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

resource "ionoscloud_datacenter_nsg_selection" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  nsg_id            = ionoscloud_nsg.example.id
}
```

## Argument reference

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `nsg_id` - (Required)[string] The ID of a Network Security Group.
