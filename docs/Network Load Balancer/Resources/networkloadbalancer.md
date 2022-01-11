---
layout: "ionoscloud"
page_title: "IonosCloud: networkloadbalancer"
sidebar_current: "docs-resource-networkloadbalancer"
description: |-
Creates and manages Network Load Balancer objects.
---

# ionoscloud_networkloadbalancer

Manages a Network Load Balancer  on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name          = "example"
  listener_lan  = ionoscloud_lan.example1.id
  target_lan    = ionoscloud_lan.example2.id
  ips           = ["10.12.118.224"]
  lb_private_ips = ["10.13.72.225/24"]
}
```

## Argument reference

- `name` - (Required)[string] A name of that Network Load Balancer.
- `listener_lan` - (Required)[int] Id of the listening LAN. (inbound)
- `ips` - (Optional)[list] Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
- `target_lan` - (Required)[int] Id of the balanced private target LAN. (outbound)
- `lb_private_ips` - (Optional)[list] Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.
- `datacenter_id` - (Required)[string] A Datacenter's UUID.

## Import

A Network Load Balancer resource can be imported using its `resource id` and the `datacenter id` e.g.

```shell
terraform import ionoscloud_networkloadbalancer.my_networkloadbalancer {datacenter uuid}/{networkloadbalancer uuid}
```