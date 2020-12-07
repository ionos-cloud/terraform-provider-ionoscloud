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
