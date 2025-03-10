---
subcategory: "Compute Engine"
layout: "ionoscloud"
page_title: "IonosCloud: loadbalancer"
sidebar_current: "docs-resource-loadbalancer"
description: |-
  Creates and manages Load Balancers
---

# ionoscloud_loadbalancer

Manages a [Load Balancer](https://docs.ionos.com/cloud/network-services/application-load-balancer/overview) on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
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
  }
}

resource "ionoscloud_loadbalancer" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  nic_ids               = [ ionoscloud_server.example.primary_nic ]
  name                  = "Load Balancer Example"
  dhcp                  = true
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
```

## Argument reference

* `name` - (Required)[string] The name of the load balancer.
* `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
* `nic_ids` - (Required)[list] A list of NIC IDs that are part of the load balancer.
* `dhcp` - (Optional)[Boolean] Indicates if the load balancer will reserve an IP using DHCP.
* `ip` - (Optional)[string] IPv4 address of the load balancer.

## Import

Resource Load Balancer can be imported using the `resource id`, e.g.

```shell
terraform import ionoscloud_loadbalancer.myloadbalancer datacenter uuid/loadbalancer uuid
```

## A note on nics

When declaring NIC resources to be used with the load balancer, please make sure
you use the "lifecycle meta-argument" to make sure changes to the lan attribute
of the nic are ignored. 

Please see the [Nic](nic.md) resource's documentation for an example on how to do that. 