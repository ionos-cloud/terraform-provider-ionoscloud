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
* `name` - (Optional) Name of an existing network load balancer forwarding rule that you want to search for.
* `id` - (Optional) ID of the network load balancer forwarding rule you want to search for.

Both `datacenter_id` and `natgateway_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `type`
* `protocol`
* `source_subnet`
* `public_ip`
* `target_subnet`
* `target_port_range` - list of
    * `start`
    * `end`
    