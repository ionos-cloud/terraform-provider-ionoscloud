---
subcategory: "NAT Gateway"
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
* `name` - (Optional) Name of an existing network load balancer forwarding rule that you want to search for.
* `id` - (Optional) ID of the network load balancer forwarding rule you want to search for.

`datacenter_id` and either `name` or `id` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of that natgateway
* `name` - Name of that natgateway
* `public_ips` - Collection of public IP addresses of the NAT gateway. Should be customer reserved IP addresses in that location
* `lans` - Collection of LANs connected to the NAT gateway. IPs must contain valid subnet mask. If user will not provide any IP then system will generate an IP with /24 subnet.
    * `id` - Id for the LAN connected to the NAT gateway
    * `gateway_ips` - Collection of gateway IP addresses of the NAT gateway. Will be auto-generated if not provided. Should ideally be an IP belonging to the same subnet as the LAN
