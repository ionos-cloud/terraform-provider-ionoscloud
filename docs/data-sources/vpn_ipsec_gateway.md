---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_vpn_ipsec_gateway"
sidebar_current: "docs-ionoscloud-datasource-vpn-ipsec-gateway"
description: |-
  Reads IonosCloud VPN IPSec Gateway objects.
---

# ionoscloud_apigateway_route

The **VPN IPSec Gateway data source** can be used to search for and return an existing IPSec Gateway.
You can provide a string for the name parameter which will be compared with provisioned IPSec Gateways.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```hcl
data "ionoscloud_vpn_ipsec_gateway" "example" {
  location = "de/fra"
  id = <gateway_id>
  location = <gateway_location>
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```hcl
data "ionoscloud_vpn_ipsec_gateway" "example" {
  location = "de/fra"
  name     = "ipsec-gateway"
  location = <gateway_location>
}
```

## Argument Reference

* `id` - (Optional) ID of an existing IPSec Gateway that you want to search for.
* `name` - (Optional) Name of an existing IPSec Gateway that you want to search for.
* `location` - (Required) The location of the IPSec Gateway.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - The unique ID of the IPSec Gateway.
* `name` - The name of the IPSec Gateway.
* `gateway_ip` - Public IP address to be assigned to the gateway.
* `description` - (Optional)[string] The human-readable description of the IPSec Gateway.
* `connections` - The network connection for your gateway.
    * `datacenter_id` - The datacenter to connect your VPN Gateway to.
    * `lan_id` - The numeric LAN ID to connect your VPN Gateway to.
    * `ipv4_cidr` - Describes the private ipv4 subnet in your LAN that should be accessible by the
      VPN Gateway.
    * `ipv6_cidr` - Describes the ipv6 subnet in your LAN that should be accessible by the VPN Gateway.
* `version` - The IKE version that is permitted for the VPN tunnels.
* `maintenance_window` - A weekly 4 hour-long window, during which maintenance might occur.
  * `time` - Start of the maintenance window in UTC time.
  * `day_of_the_week` - The name of the week day.
* `tier` - Gateway performance options.
