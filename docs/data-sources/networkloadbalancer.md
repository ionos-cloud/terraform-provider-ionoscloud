---
subcategory: "Network Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_networkloadbalancer"
sidebar_current: "docs-ionoscloud-datasource-networkloadbalancer"
description: |-
Get information on a Network Load Balancer
---

# ionoscloud_networkloadbalancer

The network load balancer data source can be used to search for and return existing network load balancers.

## Example Usage

```hcl
data "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "example_"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing network load balancer that you want to search for.
* `id` - (Optional) ID of the network load balancer you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of that Network Load Balancer
* `name`- Name of that Network Load Balancer
* `listener_lan` - Id of the listening LAN. (inbound)
* `target_lan` - Id of the balanced private target LAN. (outbound)
* `ips` - Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
* `lb_private_ips` - Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.

