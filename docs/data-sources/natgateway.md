---
subcategory: "NAT Gateway"
layout: "ionoscloud"
page_title: "IonosCloud : ionoscloud_natgateway"
sidebar_current: "docs-ionoscloud-datasource-natgateway"
description: |-
  Get information on a Nat Gateway
---

# ionoscloud_natgateway

The **NAT gateway data source** can be used to search for and return existing NAT Gateways.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search and make sure that your resources have unique names.

## Example Usage
### By ID
```hcl
data "ionoscloud_natgateway" "example" {
  datacenter_id = <datacenter_id>
  id			= <nat_gateway_id>
}
```

### By Name
```hcl
data "ionoscloud_natgateway" "example" {
  datacenter_id = <datacenter_id>
  name			= "NAT Gateway Example"
}
```

### By Name with Partial Match
```hcl
data "ionoscloud_natgateway" "example" {
  datacenter_id = <datacenter_id>
  name			= "Example"
  partial_match	= true
}
```

## Argument Reference

* `datacenter_id` - (Required) Datacenter's UUID.
* `id` - (Optional) ID of the network load balancer forwarding rule you want to search for.
* `name` - (Optional) Name of an existing network load balancer forwarding rule that you want to search for. Search by name is case-insensitive. The whole resource name is required if `partial_match` parameter is not set to true..
* `partial_match` - (Optional) Whether partial matching is allowed or not when using name argument. Default value is false.

`datacenter_id` and either `id` or `name` must be provided. If none, or both of `name` and `id` are provided, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Id of that natgateway
* `name` - Name of that natgateway
* `public_ips` - Collection of public IP addresses of the NAT gateway. Should be customer reserved IP addresses in that location
* `lans` - Collection of LANs connected to the NAT gateway. IPs must contain valid subnet mask. If user will not provide any IP then system will generate an IP with /24 subnet.
    * `id` - Id for the LAN connected to the NAT gateway
    * `gateway_ips` - Collection of gateway IP addresses of the NAT gateway. Will be auto-generated if not provided. Should ideally be an IP belonging to the same subnet as the LAN
