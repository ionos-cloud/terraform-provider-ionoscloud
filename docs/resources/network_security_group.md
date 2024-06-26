---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: network security group"
sidebar_current: "docs-resource-nsg"
description: |-
  Creates and manages IonosCloud Network Security Group.
---

# ionoscloud\_nsg

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

## Creating or setting default NSG for a datacenter
#### Either one `create_default_security_group` or `default_security_group_id` should be set, setting both will result in an error.
#### If `create_default_security_group` is set at Datacenter creation, a default NSG is created together with the datacenter, it can also be set at update to create it later. 
#### The ID is then set by terraform on the `default_security_group_id` field, this field is not `Computed` so the plan will have to be updated with the value. 
#### To set a custom NSG as default for the datacenter, set an ID value for `default_security_group_id` 
###### (Note: must specify ID as string, referencing a NSG is not possible due to resource reference cycle between datacenter and nsg)
#### Unsetting `default_security_group_id` will unset the default security group from the datacenter.
```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
  create_default_security_group = true
#   default_security_group_id = "or your security group ID"
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
terraform import ionoscloud_nsg.mynsg {datacenter uuid}/{nsg uuid}
```

Or by using an `import` block. Here is an example that allows you to import the default created nsg into terraform.
```hcl
resource "ionoscloud_datacenter" "example" {
  name            = "Datacenter NSG Example"
  location        = "de/txl"
  create_default_security_group = true    # NSG created by this flag
}

import {
  to = ionoscloud_nsg.imported
  id = "{datacenter uuid}/{default nsg uuid}" 
}
  
resource "ionoscloud_nsg" "imported_default" {  # Imported here
  datacenter_id     = ionoscloud_datacenter.example.id
}
```