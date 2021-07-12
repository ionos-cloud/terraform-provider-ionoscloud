---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_networkloadbalancer"
sidebar_current: "docs-ionoscloud-datasource-networkloadbalancer"
description: |-
Get information on a Network Load Balancer
---

# ionoscloud_networkloadbalancer

The network load balancer data source can be used to search for and return existing network load balancers.

## Example Usage

```hcl
data "ionoscloud_networkloadbalancer" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name			= "example_"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `name` - (Optional) Name of an existing network load balancer that you want to search for.
* `id` - (Optional) ID of the network load balancer you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `listener_lan`
* `target_lan`
* `ips`
* `lb_private_ips` 

