---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipfailover"
sidebar_current: "docs-resource-ipfailover"
description: |-
  Get Information on ipfailover objects.
---

# ionoscloud_ipfailover

The **IP Failover data source** can be used to search for and return an existing IP Failover object.
You need to provide the datacenter_id and the id of the lan to get the ip failover object for the provided datacenter.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.

## Example Usage

```hcl
data "ionoscloud_ipfailover" "example" {
  datacenter_id   = "datacenter_id"
  lan_id              = "lan_id"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter_id` - (Required) The ID of the datacenter containing the ip failover datasource
* `lan_id` - (Required) The id of the lan of which the IP failover belongs 


## Attributes Reference

The following attributes are returned by the datasource:

* `datacenter_id` - The ID of a Data Center.
* `ip` - The reserved IP address to be used in the IP failover group.
* `lan_id` - The ID of a LAN.
* `nicuuid` - The ID of a NIC.
