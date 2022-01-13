---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipfailover"
sidebar_current: "docs-resource-ipfailover"
description: |-
Get Information on ipfailover objects.
---

# ionoscloud\_ipfailover

The share data source can be used to search for and return an existing ip failover object.
You need to provide the datacenter_id and failover_id to get the ip failover object for the provided datacenter.


## Example Usage

```hcl
data "ionoscloud_ipfailover" "failovertest" {
  datacenter_id = "datacenterId"
  id = "failover_resource_id"
}
```

## Argument Reference

The following arguments are supported:

* `datacenter_id` - (Required)The ID of the datacenter containing the ip failover datasource
* `id` - (Required)The id of the ip failover object


## Attributes Reference

The following attributes are returned by the datasource:

* `datacenter_id` - The ID of a Data Center.
* `ip` - The reserved IP address to be used in the IP failover group.
* `lan_id` - The ID of a LAN.
* `nicuuid` - The ID of a NIC.
