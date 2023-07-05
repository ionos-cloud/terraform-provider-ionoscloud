---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: private_crossconnect"
sidebar_current: "docs-resource-private-crossconnect"
description: |-
  Creates and manages Private Cross Connections between virtual datacenters.
---

# ionoscloud_private_crossconnect

Manages a **Private Cross Connect** on IonosCloud.

## Example Usage

To connect two datacenters we need 2 lans defined, one in each datacenter. After, we reference the pcc through which we want the connection to be established.

```hcl
resource ionoscloud_private_crossconnect PCCTestResource {
  name        = "PCCTestResource"
  description = "PCCTestResource"
}

resource ionoscloud_datacenter dc1 {
  location = "de/txl"
  name = "dc1"
}

resource ionoscloud_datacenter dc2 {
  location = "de/txl"
  name = "dc2"
}

resource ionoscloud_lan dc1lan {
  datacenter_id = ionoscloud_datacenter.dc1.id
  public = false
  name = "dc1lan"
  pcc = ionoscloud_private_crossconnect.PCCTestResource.id
}

resource ionoscloud_lan dc2lan {
  datacenter_id = ionoscloud_datacenter.dc2.id
  public = false
  name = "dc2lan"
  pcc = ionoscloud_private_crossconnect.PCCTestResource.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the private cross-connection.
- `description` - (Optional)[string] A short description for the private cross-connection.
- `connectable datacenters` - (Computed) A list containing all the connectable datacenters
  - `id` - The UUID of the connectable datacenter
  - `name` - The name of the connectable datacenter
  - `location` - The physical location of the connectable datacenter
- `peers` - (Computed) Lists LAN's joined to this private cross connect
  - `lan_id` - The id of the cross-connected LAN
  - `lan_name` - The name of the cross-connected LAN
  - `datacenter_id` - The id of the cross-connected datacenter
  - `datacenter_name` - The name of the cross-connected datacenter
  - `location` - The location of the cross-connected datacenter
  
## Import

A Private Cross Connect resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_private_crossconnect.demo {ionoscloud_private_crossconnect_uuid}
```

This can be helpful when you want to import private cross-connects which you have already created manually or using other means, outside of terraform.
