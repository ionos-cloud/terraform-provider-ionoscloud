---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_natgateway"
sidebar_current: "docs-ionoscloud-datasource-natgateway"
description: |-
Get information on a Nat Gateway
---

# ionoscloud_natgateway

The nat gateway data source can be used to search for and return existing natgateways.

## Example Usage

```hcl
data "ionoscloud_natgateway" "natgateway_example" {
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
  name			= "example_"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing network loadbalancer forwarding rule that you want to search for.
* `id` - (Optional) ID of the network loadbalancer forwarding rule you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `public_ips`
* `lans` - list of
    * `id`
    * `gateway_ips`
