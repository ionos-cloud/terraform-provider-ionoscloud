---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IONOS CLOUD: lan"
sidebar_current: "docs-resource-lan"
description: |-
  Creates and manages LAN objects.
---

# ionoscloud_lan

Manages a **LAN** on IONOS CLOUD.

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

## Example Usage With a Custom IPv4 CIDR Block

By default, a private LAN gets a /23 block from the `10.0.0.0/8` range. Set `ipv4_cidr_block` to choose the block yourself, e.g. a /24 for general use or a /28 for a small subnet:

```hcl
resource "ionoscloud_datacenter" "example" {
  name                  = "Datacenter Example"
  location              = "de/txl"
  description           = "Datacenter Description"
  sec_auth_protection   = false
}

resource "ionoscloud_lan" "general" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = false
  name                  = "Lan /24 Example"
  ipv4_cidr_block       = "10.5.0.0/24"
}

resource "ionoscloud_lan" "small" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = false
  name                  = "Lan /28 Example"
  ipv4_cidr_block       = "192.168.42.0/28"
}
```

## Argument reference

* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `name` - (Optional)[string] The name of the LAN.
* `public` - (Optional)[Boolean] Indicates if the LAN faces the public Internet (true) or not (false).
* `pcc` - (Optional)[String] The unique id of a `ionoscloud_private_crossconnect` resource, in order. It needs to be ensured that IP addresses of the NICs of all LANs connected to a given Cross Connect is not duplicated and belongs to the same subnet range
* `ipv4_cidr_block` - (Computed, Optional)[String] The IPv4 CIDR block of a private LAN, e.g. `10.5.0.0/24`. If omitted, a /23 block from the `10.0.0.0/8` range is assigned automatically. A custom block must be within a private RFC 1918 range (`10.0.0.0/8`, `172.16.0.0/12` or `192.168.0.0/16`), have a prefix length between /23 and /28, and be aligned to its prefix boundary (`10.5.0.0/24`, not `10.5.0.1/24`). It can only be set on private LANs. Changing it updates the LAN in place; removing it from the configuration keeps the current block. If the API rejects a block, the error contains the API error code: `536` (not aligned to its prefix boundary), `602` (overlaps a block already in use), `654` (not a private range), `655` (larger than /23) or `656` (set on a public LAN).
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