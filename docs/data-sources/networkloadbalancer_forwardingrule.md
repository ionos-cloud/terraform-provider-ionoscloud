---
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_networkloadbalancer_forwardingrule"
sidebar_current: "docs-ionoscloud-datasource-networkloadbalancer_forwardingrule"
description: |-
Get information on a Network Loadbalancer Forwarding Rule
---

# ionoscloud_networkloadbalancer_forwardingrule

The network loadbalancer forwarding rule data source can be used to search for and return existing network forwarding rules.

## Example Usage

```hcl
data "ionoscloud_networkloadbalancer_forwardingrule" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  networkloadbalancer_id  = ionoscloud_networkloadbalancer.example.id
  name			= "example_"
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `networkloadbalancer_id` - (Required) Network Loadbalancer's UUID.
* `name` - (Optional) Name of an existing network loadbalancer forwarding rule that you want to search for.
* `id` - (Optional) ID of the network loadbalancer forwarding rule you want to search for.

Both `datacenter_id` and `networkloadbalancer_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id`
* `name`
* `algorithm`
* `protocol`
* `listener_ip`
* `listener_port`
* `health_check` - list of
    * `client_timeout`
    * `check_timeout`
    * `connect_timeout`
    * `target_timeout`
    * `retries`
* `targets`
    * `ip`
    * `port`
    * `weight`
    * `health_check`
        * `check`
        * `check_interval`
        * `maintenance`

