---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_application_loadbalancer"
sidebar_current: "docs-ionoscloud_application_loadbalancer"
description: |-
Get information on an Application Load Balancer
---

# ionoscloud_application_loadbalancer

The Application Load Balancer data source can be used to search for and return existing Application Load Balancers.

## Example Usage

### By Id
```hcl
data "ionoscloud_application_loadbalancer" "alb_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  id			= <alb_uuid>
}
```
### By Name
```hcl
data "ionoscloud_application_loadbalancer" "alb_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  name			= "alb_example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing application load balancer that you want to search for.
* `id` - (Optional) ID of the application load balancer you want to search for.

`datacenter_id`  and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the application load balancer.
- `name` - Name of the application load balancer.
- `listener_lan` - Id of the listening LAN. (inbound).
- `ips` - Collection of IP addresses of the Application Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
- `target_lan` - Id of the balanced private target LAN. (outbound).
- `lb_private_ips` - Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.
