---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: lan"
sidebar_current: "docs-resource-lan"
description: |-
  Creates and manages LAN objects.
---

# ionoscloud_lan

Manages a **LAN** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_private_crossconnect" "example" {
  name                  = "Cross Connect Example"
  description           = "Cross Connect Description"
}

resource "ionoscloud_lan" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = false
  name                  = "Lan Example"
  pcc                   = ionoscloud_private_crossconnect.example.id
}
```

## Example Usage With IPv6 Enabled

```hcl
resource "ionoscloud_datacenter" "example" {
  name                  = "Datacenter Example"
  location              = "de/txl"
  description           = "Datacenter Description"
  sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = true
  name                  = "Lan IPv6 Example"
  ipv6_cidr_block       = "AUTO"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the LAN.
* `public` - (Optional)[Boolean] Indicates if the LAN faces the public Internet (true) or not (false).
* `pcc` - (Optional)[String] The unique id of a `ionoscloud_private_crossconnect` resource, in order. It needs to be ensured that IP addresses of the NICs of all LANs connected to a given Cross Connect is not duplicated and belongs to the same subnet range
* `ipv4_cidr_block` - (Computed)[String] For public LANs this property is null, for private LANs it contains the private IPv4 CIDR range. This property is a read only property.
* `ipv6_cidr_block` - (Computed, Optional) Contains the LAN's /64 IPv6 CIDR block if this LAN is IPv6 enabled. 'AUTO' will result in enabling this LAN for IPv6 and automatically assign a /64 IPv6 CIDR block to this LAN. If you specify your own IPv6 CIDR block then you must provide a unique /64 block, which is inside the IPv6 CIDR block of the virtual datacenter and unique inside all LANs from this virtual datacenter.
* `ip_failover` - (Computed) IP failover configurations for lan
  * `ip`
  * `nic_uuid`
  
## Import

Resource Lan can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_lan.mylandatacenter uuid/lan id
```

## Important Notes

- Please note that only LANs datacenters found in the same physical location can be connected through a Cross-connect
- A LAN cannot be a part of two Cross-connects