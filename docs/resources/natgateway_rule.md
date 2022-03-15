---
subcategory: "NAT Gateway"
layout: "ionoscloud"
page_title: "IonosCloud: natgateway_rule"
sidebar_current: "docs-resource-natgateway_rule"
description: |-
  Creates and manages Nat Gateway Rule objects.
---

# ionoscloud_natgateway_rule

Manages a **Nat Gateway Rule** on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_datacenter" "example" {
	name                    = "Datacenter Example"
	location                = "us/las"
	description             = "Datacenter Description"
	sec_auth_protection     = false
}

resource "ionoscloud_ipblock" "example" {
    location                = "us/las"
    size                    = 2
    name                    = "IP Block Example"
}

resource "ionoscloud_lan" "example" {
    datacenter_id           = ionoscloud_datacenter.example.id
    public                  = true
    name                    = "Lan Example"
}

resource "ionoscloud_natgateway" "example" {
    datacenter_id           = ionoscloud_datacenter.example.id
    name                    = "example"
    public_ips              = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
     lans {
        id                  = ionoscloud_lan.example.id
        gateway_ips         = [ "10.11.2.5"]
     }
}

resource "ionoscloud_natgateway_rule" "example" {
    datacenter_id           = ionoscloud_datacenter.example.id
    natgateway_id           = ionoscloud_natgateway.example.id
    name                    = "example"
    type                    = "SNAT"
    protocol                = "TCP"
    source_subnet           = "10.0.1.0/24"
    public_ip               = ionoscloud_ipblock.example.ips[0]
    target_subnet           = "10.0.1.0/24"
    target_port_range {
        start               = 500
        end                 = 1000
    }
}
```

## Argument reference

- `name` - (Required)[string] Name of the NAT gateway rule.
- `type` - (Optional)[string] Type of the NAT gateway rule.
- `protocol` - (Optional)[string] Protocol of the NAT gateway rule. Defaults to ALL. If protocol is 'ICMP' then targetPortRange start and end cannot be set.
- `source_subnet` - (Required)[string] Source subnet of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on the packets source IP address.
- `public_ip` - (Required)[string] Public IP address of the NAT gateway rule. Specifies the address used for masking outgoing packets source address field. Should be one of the customer reserved IP address already configured on the NAT gateway resource.
- `target_subnet` - (Optional)[string] Target or destination subnet of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on the packets destination IP address. If none is provided, rule will match any address.
- `target_port_range` - (Optional) Target port range of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on destination port. If none is provided, rule will match any port.
    - `start` - (Optional)[int] Target port range start associated with the NAT gateway rule.
    - `end` - (Optional)[int] Target port range end associated with the NAT gateway rule.
- `datacenter_id` - (Required)[string] A Datacenter's UUID.
- `natgateway_id` - (Required)[string] Nat Gateway's UUID.

## Import

A Nat Gateway Rule resource can be imported using its `resource id`, the `datacenter id` and the `natgateway id , e.g.

```shell
terraform import ionoscloud_natgateway_rule.my_natgateway_rule {datacenter uuid}/{nat gateway uuid}/{nat gateway rule uuid}
```
