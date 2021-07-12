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
resource "ionoscloud_application_loadbalancer" "example" { 
  datacenter_id = ionoscloud_datacenter.example.id
  name          = "example"
  listener_lan  = ionoscloud_lan.example1.id
  ips           = [ "81.173.1.2",
                    "22.231.2.2",
                    "22.231.2.3"
                  ]
  target_lan    = ionoscloud_lan.example2.id
  lb_private_ips= [ "81.173.1.5/24",
                    "22.231.2.5/24"
                  ]
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
