---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_natgateway_rule"
sidebar_current: "docs-ionoscloud-datasource-natgateway_rule"
description: |-
Get information on a Nat Gateway Rule
---

# ionoscloud_natgateway_rule

The nat gateway rule data source can be used to search for and return existing natgateway rules.

## Example Usage

```hcl
data "ionoscloud_natgateway_rule" "natgateway_rule_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  natgateway_id = ionoscloud_natgateway.natgateway_example.id
  name			= "example_"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `natgateway_id` - (Required) Nat Gateway's UUID.
* `name` - (Optional) Name of an existing network loadbalancer forwarding rule that you want to search for.
* `id` - (Optional) ID of the network loadbalancer forwarding rule you want to search for.

Both `datacenter_id` and `natgateway_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of the NAT gateway rule
* `name` - Name of the NAT gateway rule
* `type` - ype of the NAT gateway rule.
* `protocol` - Protocol of the NAT gateway rule. Defaults to ALL. If protocol is 'ICMP' then targetPortRange start and end cannot be set.
* `source_subnet` - Source subnet of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on the packets source IP address.
* `public_ip` - Public IP address of the NAT gateway rule. Specifies the address used for masking outgoing packets source address field. Should be one of the customer reserved IP address already configured on the NAT gateway resource
* `target_subnet` - Target or destination subnet of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on the packets destination IP address. If none is provided, rule will match any address.
* `target_port_range` - Target port range of the NAT gateway rule. For SNAT rules it specifies which packets this translation rule applies to based on destination port. If none is provided, rule will match any port
    * `start`
    * `end`
    