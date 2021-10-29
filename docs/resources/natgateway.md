---
layout: "ionoscloud"
page_title: "IonosCloud: natgateway"
sidebar_current: "docs-resource-natgateway"
description: |-
  Creates and manages Nat Gateway objects.
---

# ionoscloud_natgateway

Manages a Nat Gateway on IonosCloud.

## Example Usage

```hcl
resource "ionoscloud_natgateway" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  name          = "example"
  public_ips    =  ["${ionoscloud_ipblock.example.ips[0]}", "${ionoscloud_ipblock.example.ips[1]}"]
  lans {
     id          = ionoscloud_lan.example.id
     gateway_ips = [ "10.11.2.5/32"]
  }
}
```

## Argument reference

- `name` - (Required)[string] Name of the NAT gateway.
- `public_ips` - (Required)[list]Collection of public IP addresses of the NAT gateway. Should be customer reserved IP addresses in that location.
- `lans` - (Required)[list] A list of Local Area Networks the node pool should be part of.
  - `id` - (Required)[int] Id for the LAN connected to the NAT gateway.
  - `gateway_ips` - (Optional)[list] Collection of gateway IP addresses of the NAT gateway. Will be auto-generated if not provided. Should ideally be an IP belonging to the same subnet as the LAN.
- `datacenter_id` - (Required)[string] A Datacenter's UUID.


## Import

A Nat Gateway resource can be imported using its `resource id` and the `datacenter id`, e.g.

```shell
terraform import ionoscloud_natgateway.my_natgateway {datacenter uuid}/{nat gateway uuid}
```
