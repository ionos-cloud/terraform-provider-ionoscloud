---
subcategory: "Application Load Balancer"
layout: "ionoscloud"
page_title: "IonosCloud: application_loadbalancer"
sidebar_current: "docs-resource-application_loadbalancer"
description: |-
  Creates and manages IonosCloud Application Load Balancer.
---

# ionoscloud_application_loadbalancer

Manages an **Application Load Balancer** on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_datacenter" "example" {
    name                  = "Datacenter Example"
    location              = "us/las"
    description           = "datacenter description"
    sec_auth_protection   = false
}

resource "ionoscloud_lan" "example_1" {
    datacenter_id         = ionoscloud_datacenter.example.id
    public                = true
    name                  = "Lan Example"
}

resource "ionoscloud_lan" "example_2" {
    datacenter_id         = ionoscloud_datacenter.example.id
    public                = true
    name                  = "Lan Example"
}

resource "ionoscloud_application_loadbalancer" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    name                  = "ALB Example"
    listener_lan          = ionoscloud_lan.example_1.id
    ips                   = [ "10.12.118.224"]
    target_lan            = ionoscloud_lan.example_2.id
    lb_private_ips        = [ "10.13.72.225/24"]
}

```

## Argument Reference

The following arguments are supported:

- `datacenter_id` - (Required)[string] ID of the datacenter.
- `name` - (Required)[string] The name of the Application Load Balancer.
- `listener_lan` - (Required)[int] ID of the listening (inbound) LAN.
- `ips` - (Optional)[set] Collection of the Application Load Balancer IP addresses. (Inbound and outbound) IPs of the listenerLan are customer-reserved public IPs for the public Load Balancers, and private IPs for the private Load Balancers.
- `target_lan` - (Required)[int] ID of the balanced private target LAN (outbound).
- `lb_private_ips` - (Optional)[set] Collection of private IP addresses with the subnet mask of the Application Load Balancer. IPs must contain valid a subnet mask. If no IP is provided, the system will generate an IP with /24 subnet.


## Import

Resource Application Load Balancer can be imported using the `resource id` and `datacenter id`, e.g.

```shell
terraform import ionoscloud_application_loadbalancer.myalb {datacenter uuid}/{applicationLoadBalancer uuid}
```