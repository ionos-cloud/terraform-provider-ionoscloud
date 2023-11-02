---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: private_crossconnect"
sidebar_current: "docs-resource-private-crossconnect"
description: |-
  Creates and manages Cross Connections between virtual datacenters.
---

# ionoscloud_private_crossconnect

Manages a **Cross Connect** on IonosCloud.
Cross Connect allows you to connect virtual data centers (VDC) with each other using a private LAN. 
The VDCs to be connected need to belong to the same IONOS Cloud contract and location. 
You can only use private LANs for a Cross Connect connection. A LAN can only be a part of one Cross Connect.

The IP addresses of the NICs used for the Cross Connect connection may not be used in more than one NIC and they need to belong to the same IP range.

## Example Usage

To connect two datacenters we need 2 lans defined, one in each datacenter. After, we reference the cross-connect through which we want the connection to be established.

```hcl
resource ionoscloud_private_crossconnect CrossConnectTestResource {
  name        = "CrossConnectTestResource"
  description = "CrossConnectTestResource"
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
  pcc = ionoscloud_private_crossconnect.CrossConnectTestResource.id
}

resource ionoscloud_lan dc2lan {
  datacenter_id = ionoscloud_datacenter.dc2.id
  public = false
  name = "dc2lan"
  pcc = ionoscloud_private_crossconnect.CrossConnectTestResource.id
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required)[string] The name of the cross-connection.
- `description` - (Optional)[string] A short description for the cross-connection.
- `connectable datacenters` - (Computed) A list containing all the connectable datacenters
  - `id` - The UUID of the connectable datacenter
  - `name` - The name of the connectable datacenter
  - `location` - The physical location of the connectable datacenter
- `peers` - (Computed) Lists LAN's joined to this cross connect
  - `lan_id` - The id of the cross-connected LAN
  - `lan_name` - The name of the cross-connected LAN
  - `datacenter_id` - The id of the cross-connected datacenter
  - `datacenter_name` - The name of the cross-connected datacenter
  - `location` - The location of the cross-connected datacenter
  
## Import

A Cross Connect resource can be imported using its `resource id`, e.g.

```shell
terraform import ionoscloud_private_crossconnect.demo {ionoscloud_private_crossconnect_uuid}
```

This can be helpful when you want to import cross-connects which you have already created manually or using other means, outside of terraform.
