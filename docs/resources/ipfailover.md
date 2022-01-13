---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipfailover"
sidebar_current: "docs-resource-ipfailover"
description: |-
  Creates and manages ipfailover objects.
---

# ionoscloud\_ipfailover

Manages IP Failover groups on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_ipfailover" "failovertest" {
  datacenter_id = "datacenterId"
  lan_id="lanId"
  ip ="reserved IP"
  nicuuid= "nicId"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `ip` - (Required)[string] The reserved IP address to be used in the IP failover group.
* `lan_id` - (Required)[string] The ID of a LAN.
* `nicuuid` - (Required)[string] The ID of a NIC.

## Import

Resource IpFailover can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_ipfailover.myipfailover {datacenter uuid}/{lan uuid}
```
