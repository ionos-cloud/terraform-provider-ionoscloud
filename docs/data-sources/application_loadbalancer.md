---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_application_loadbalancer"
sidebar_current: "docs-ionoscloud_application_loadbalancer"
description: |-
  Get information on an Application Load Balancer
---

# ionoscloud_application_loadbalancer

The **Application Load Balancer data source** can be used to search for and return an existing Application Load Balancer.
You can provide a string for the name parameter which will be compared with provisioned Application Load Balancers.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please make sure that your resources have unique names.

## Example Usage

### By Id
```hcl
data "ionoscloud_application_loadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  id			= <alb_id>
}
```
### By Name
```hcl
data "ionoscloud_application_loadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "ALB Example"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing application load balancer that you want to search for. Search by name is case-insensitive, but the whole resource name is required (we do not support partial matching).
* `id` - (Optional) ID of the application load balancer you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the application load balancer.
- `name` - The name of the Application Load Balancer.
- `listener_lan` - ID of the listening (inbound) LAN.
- `ips` - Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
- `target_lan` - ID of the balanced private target LAN (outbound).
- `lb_private_ips` - Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.
