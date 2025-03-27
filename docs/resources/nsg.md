---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: network security group"
sidebar_current: "docs-resource-nsg"
description: |-
  Creates and manages IonosCloud Network Security Group.
---

# ionoscloud_nsg

Manages a **Network Security Group** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
}

resource "ionoscloud_nsg" "example" {
  name              = "Example NSG"
  description       = "Example NSG Description"
  datacenter_id     = ionoscloud_datacenter.example.id
}
```

## Argument Reference

The following arguments are supported:
* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the Network Security Group.
* `description` - (Optional)[string] Description for the Network Security Group.
* `rule_ids` - (Computed) List of Firewall Rules that are part of the Network Security Group

## Import

Resource Server can be imported using the `resource id` and the `datacenter id`, e.g.

```shell
terraform import ionoscloud_nsg.mynsg datacenter uuid/nsg uuid
```

Or by using an `import` block. Here is an example that allows you to import the default created nsg into pulumi.
```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
}

import {
  to = ionoscloud_nsg.imported
  id = "datacenter uuid/default nsg uuid" 
}
  
resource "ionoscloud_nsg" "imported_default" {  # Imported here
  datacenter_id     = ionoscloud_datacenter.example.id
}
```