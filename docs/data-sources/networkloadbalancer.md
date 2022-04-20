---
subcategory: "Network Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_networkloadbalancer"
sidebar_current: "docs-ionoscloud-datasource-networkloadbalancer"
description: |-
  Get information on a Network Load Balancer
---

# ionoscloud_networkloadbalancer

The **Network Load Balancer data source** can be used to search for and return existing network load balancers.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage

### By ID
```hcl
data "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  id			= <networkloadbalancer_id>
}
```

### By Name
```hcl
data "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "Network Load Balancer Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "Example"
  partial_match = true
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `id` - (Optional) ID of the network load balancer you want to search for.
* `name` - (Optional) Name of an existing network load balancer that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

`datacenter_id` and either `id` or `name` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of that Network Load Balancer
* `name`- Name of that Network Load Balancer
* `listener_lan` - Id of the listening LAN. (inbound)
* `target_lan` - Id of the balanced private target LAN. (outbound)
* `ips` - Collection of IP addresses of the Network Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
* `lb_private_ips` - Collection of private IP addresses with subnet mask of the Network Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.

