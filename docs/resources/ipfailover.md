---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: ipfailover"
sidebar_current: "docs-resource-ipfailover"
description: |-
  Creates and manages ipfailover objects.
---

# ionoscloud_ipfailover

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
  cpu_family            = "INTEL_XEON"
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

> **⚠ WARNING:** Do not modify the IP for an IP failover group (that was provisioned via Pulumi)
> using the DCD, the API or other means because it may lead to unexpected behavior. If you provisioned
> an IP failover group using Pulumi, please use only Pulumi in order to manage the created
> IP failover group.

> **⚠ WARNING:** For creating multiple IP failover groups at the same time, you can use one of the
> following options:
1. Create multiple IP failover groups resources and use `depends_on` meta-argument to specify the order
of creation, for example:
```example
resource "ionoscloud_ipfailover" "firstexample" {
  datacenter_id         = "datacenter ID"
  lan_id                = "LAN ID"
  ip                    = "IP address"
  nicuuid               = "NIC UUID"
}

 resource "ionoscloud_ipfailover" "secondexample" {
   depends_on = [ ionoscloud_ipfailover.firstexample ]
   datacenter_id         = "datacenter ID"
   lan_id                = "LAN ID"
   ip                    = "IP address"
   nicuuid               = "NIC UUID"
 }
```
2. Define the resources as presented above, but without using the `depends_on` meta-argument and run the apply command using
`-parallelism=1` as presented below:
```shell
terraform apply -parallelism=1
```

## Import

Resource IpFailover can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_ipfailover.myipfailover datacenter uuid/lan uuid
```


## A note on multiple NICs on an IP Failover
If you want to add a secondary NIC to an IP Failover, follow these steps:
1) Creating NIC A with failover IP on LAN 1
2) Create NIC B unde the same LAN but with a different IP
3) Create the IP Failover on LAN 1 with NIC A and failover IP of NIC A (A becomes now "master", no slaves)
4) Update NIC B IP to be the failover IP ( B becomes now a slave, A remains master)

After this you can create a new NIC C, NIC D and so on, in LAN 1, directly with the failover IP.

Please check [examples](../../examples/ip_failover) for a full example with the above steps.