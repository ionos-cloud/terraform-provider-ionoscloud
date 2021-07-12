---
layout: "ionoscloud"
page_title: "IonosCloud: loadbalancer"
sidebar_current: "docs-resource-loadbalancer"
description: |-
  Creates and manages Load Balancers
---

# ionoscloud\_loadbalancer

Manages a Load Balancer on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_loadbalancer" "example" {
  datacenter_id = "${ionoscloud_datacenter.example.id}"
  nic_ids        = ["${ionoscloud_nic.example.id}"]
  name          = "load balancer name"
  dhcp          = true
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
terraform import ionoscloud_loadbalancer.myloadbalancer {datacenter uuid}/{loadbalancer uuid}
```

## A note on nics

When declaring NIC resources to be used with the load balancer, please make sure
you use the "lifecycle meta-argument" to make sure changes to the lan attribute
of the nic are ignored. 

Please see the nic resource's documentation for an example on how to do that. 