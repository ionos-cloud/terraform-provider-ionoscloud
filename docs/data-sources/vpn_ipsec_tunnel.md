---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_vpn_ipsec_tunnel"
sidebar_current: "docs-ionoscloud-datasource-vpn-ipsec-tunnel"
description: |-
  Reads IonosCloud VPN IPSec Gateway Tunnel objects.
---

# ionoscloud_vpn_ipsec_tunnel

The **VPN IPSec Gateway Tunnel data source** can be used to search for and return an existing IPSec Gateway Tunnel.
You can provide a string for the name parameter which will be compared with provisioned IPSec Gateway Tunnels.
If a single match is found, it will be returned. If your search results in multiple matches, an error will be returned.
When this happens, please refine your search string so that it is specific enough to return only one result.

## Example Usage

### By ID

```hcl
data "ionoscloud_vpn_ipsec_tunnel" "example" {
  id = "tunnel_id"
  gateway_id = "gateway_id"
  location = "gateway_location"
}
```

### By Name

Needs to have the resource be previously created, or a depends_on clause to ensure that the resource is created before
this data source is called.

```hcl
data "ionoscloud_vpn_ipsec_gateway" "example" {
  name     = "ipsec-tunnel"
  gateway_id = "gateway_id"
  location = "gateway_location"
}
```

## Argument Reference

* `id` - (Optional) ID of an existing IPSec Gateway Tunnel that you want to search for.
* `name` - (Optional) Name of an existing IPSec Gateway Tunnel that you want to search for.
* `gateway_id` - (Required) The ID of the IPSec Gateway that the tunnel belongs to.
* `location` - (Optional) The location of the IPSec Gateway Tunnel.

## Attributes reference

The following attributes are returned by the datasource:

* `id` - The unique ID of the IPSec Gateway Tunnel.
* `name` - The name of the IPSec Gateway Tunnel.
* `description` - The human-readable description of your IPSec Gateway Tunnel.
* `remote_host` - The remote peer host fully qualified domain name or public IPV4 IP to connect to.
* `ike` - Settings for the initial security exchange phase.
    * `diffie_hellman_group` - The Diffie-Hellman Group to use for IPSec Encryption.
    * `encryption_algorithm` - The encryption algorithm to use for IPSec Encryption.
    * `integrity_algorithm` - The integrity algorithm to use for IPSec Encryption.
    * `lifetime` - The phase lifetime in seconds.
* `esp` - Settings for the IPSec SA (ESP) phase.
    * `diffie_hellman_group` - The Diffie-Hellman Group to use for IPSec Encryption.
    * `encryption_algorithm` - The encryption algorithm to use for IPSec Encryption.
    * `integrity_algorithm` - The integrity algorithm to use for IPSec Encryption.
    * `lifetime` - The phase lifetime in seconds.
* `auth` - Properties with all data needed to define IPSec Authentication.
    * `method` - The authentication method to use for IPSec Authentication.
* `cloud_network_cidrs` - The network CIDRs on the "Left" side that are allowed to connect to the IPSec
  tunnel, i.e. the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses.
* `peer_network_cidrs` - The network CIDRs on the "Right" side that are allowed to connect to the IPSec
  tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses.
