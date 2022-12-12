---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipfailover"
sidebar_current: "docs-resource-ipfailover"
description: |-
  Creates and manages ipfailover objects.
---

# ionoscloud\_ipfailover

Manages **IP Failover** groups on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_ipblock" "example" {
  location              = "us/las"
  size                  = 1
  name                  = "IP Block Example"
}

resource "ionoscloud_lan" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  public                = true
  name                  = "Lan Example"
}

resource "ionoscloud_server" "example" {
  name                  = "Server Example"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
  availability_zone     = "ZONE_1"
  cpu_family            = "AMD_OPTERON"
  image_name            = "Ubuntu-20.04"
  image_password        = random_password.server_image_password.result
  volume {
    name                = "system"
    size                = 14
    disk_type           = "SSD"
  }
  nic {
    lan                 = "1"
    dhcp                = true
    firewall_active     = true
    ips                 = [ ionoscloud_ipblock.example.ips[0] ]
  }
}

resource "ionoscloud_ipfailover" "example" {
  depends_on            = [ ionoscloud_lan.example ]
  datacenter_id         = ionoscloud_datacenter.example.id
  lan_id                = ionoscloud_lan.example.id
  ip                    = ionoscloud_ipblock.example.ips[0]
  nicuuid               = ionoscloud_server.example.primary_nic
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = false
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


## A note on multiple NICs on an IP Failover
If you want to add a secondary NIC to an IP Failover, follow these steps:
1) Creating NIC A with failover IP on LAN 1
2) Create NIC B unde the same LAN but with a different IP
3) Create the IP Failover on LAN 1 with NIC A and failover IP of NIC A (A becomes now "master", no slaves)
4) Update NIC B IP to be the failover IP ( B becomes now a slave, A remains master)

After this you can create a new NIC C, NIC D and so on, in LAN 1, directly with the failover IP.

Please check [examples](../../examples/ip_failover) for a full example with the above steps.