---
layout: "ionoscloud"
page_title: "IonosCloud: application_loadbalancer"
sidebar_current: "docs-resource-application_loadbalancer"
description: |-
Creates and manages IonosCloud Application Load Balancer.
---

# ionoscloud_application_loadbalancer

Manages an Application Load Balancer on IonosCloud.

## Example Usage

```hcl

resource "ionoscloud_datacenter" "datacenter_example" {
  name              = "test_alb"
  location          = "de/txl"
  description       = "datacenter_example"
}

resource "ionoscloud_lan" "lan_example_1" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id 
  public        = false
  name          = "lan_example_1"
}

resource "ionoscloud_lan" "lan_example_2" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id  
  public        = false
  name          = "lan_example_2"
}


resource "ionoscloud_application_loadbalancer" "example" { 
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  name          = "example"
  listener_lan  = ionoscloud_lan.lan_example_1.id
  ips           = [ "10.12.118.224"]
  target_lan    = ionoscloud_lan.lan_example_2.id
  lb_private_ips= [ "10.13.72.225/24"]
}
```

## Argument Reference

The following arguments are supported:

- `datacenter_id` - (Required)[string] The ID of a Virtual Data Center.
- `name` - (Required)[string] Name of the application load balancer.
- `listener_lan` - (Required)[int] Id of the listening LAN. (inbound).
- `ips` - (Optional)[list] Collection of IP addresses of the Application Load Balancer. (inbound and outbound) IP of the listenerLan must be a customer reserved IP for the public load balancer and private IP for the private load balancer.
- `target_lan` - (Required)[int] Id of the balanced private target LAN. (outbound).
- `lb_private_ips` - (Optional)[list] Collection of private IP addresses with subnet mask of the Application Load Balancer. IPs must contain valid subnet mask. If user will not provide any IP then the system will generate one IP with /24 subnet.


## Import

Resource Application Load Balancer can be imported using the `resource id` and `datacenter id`, e.g.

```shell
terraform import ionoscloud_application_loadbalancer.myalb {datacenter uuid}/{applicationLoadBalancer uuid}
```