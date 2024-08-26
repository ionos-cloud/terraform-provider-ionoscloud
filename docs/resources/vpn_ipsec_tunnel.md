---
subcategory: "VPN"
layout: "ionoscloud"
page_title: "IonosCloud: ionoscloud_vpn_ipsec_tunnel"
sidebar_current: "docs-resource-vpn-ipsec-tunnel"
description: |-
  IPSec Gateway Tunnel
---

# ionoscloud_vpn_ipsec_tunnel

An IPSec Gateway Tunnel resource manages the creation, management, and deletion of VPN IPSec Gateway Tunnels within the
IONOS Cloud infrastructure. This resource facilitates the creation of VPN IPSec Gateway Tunnels, enabling secure
connections between your network resources.

## Usage example

```hcl
# Basic example

resource "ionoscloud_datacenter" "test_datacenter" {
  name = "test_vpn_gateway_basic"
  location = "de/fra"
}

resource "ionoscloud_lan" "test_lan" {
  name = "test_lan_basic"
  public = false
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
}

resource "ionoscloud_ipblock" "test_ipblock" {
  name = "test_ipblock_basic"
  location = "de/fra"
  size = 1
}

resource "ionoscloud_vpn_ipsec_gateway" "example" {
  name = "ipsec_gateway_basic"
  location = "de/fra"
  gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
  version = "IKEv2"
  description = "This gateway connects site A to VDC X."

  connections {
    datacenter_id = ionoscloud_datacenter.test_datacenter.id
    lan_id = ionoscloud_lan.test_lan.id
    ipv4_cidr = "192.168.100.10/24"
  }
}

resource "ionoscloud_vpn_ipsec_tunnel" "example" {
    location = "es/vit"
    gateway_id = ionoscloud_vpn_ipsec_gateway.example.id
    
    name = "example-tunnel"
    remote_host = "vpn.mycompany.com"
    description = "Allows local subnet X to connect to virtual network Y."
    
    auth {
        method = "PSK"
        psk_key = "X2wosbaw74M8hQGbK3jCCaEusR6CCFRa"
    }
    
    ike {
        diffie_hellman_group = "16-MODP4096"
        encryption_algorithm = "AES256"
        integrity_algorithm = "SHA256"
        lifetime             = 86400
    }
    
    esp {
        diffie_hellman_group = "16-MODP4096"
        encryption_algorithm = "AES256"
        integrity_algorithm = "SHA256"
        lifetime             = 3600
    }
    
    cloud_network_cidrs = [
        "0.0.0.0/0"
    ]
    
    peer_network_cidrs = [
        "1.2.3.4/32"
    ]
}
```

```hcl
# Complete example

resource "ionoscloud_datacenter" "test_datacenter" {
  name = "vpn_gateway_test"
  location = "de/fra"
}

resource "ionoscloud_lan" "test_lan" {
  name = "test_lan"
  public = false
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  ipv6_cidr_block = local.lan_ipv6_cidr_block
}

resource "ionoscloud_ipblock" "test_ipblock" {
  name = "test_ipblock"
  location = "de/fra"
  size = 1
}

resource "ionoscloud_server" "test_server" {
  name = "test_server"
  datacenter_id = ionoscloud_datacenter.test_datacenter.id
  cores = 1
  ram = 2048
  image_name = "ubuntu:latest"
  image_password = random_password.server_image_password.result

  nic {
    lan = ionoscloud_lan.test_lan.id
    name = "test_nic"
    dhcp = true
    dhcpv6 = false
    ipv6_cidr_block = local.ipv6_cidr_block
    firewall_active   = false
  }

  volume {
    name         = "test_volume"
    disk_type    = "HDD"
    size         = 10
    licence_type = "OTHER"
  }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}

locals {
  lan_ipv6_cidr_block_parts = split("/", ionoscloud_datacenter.test_datacenter.ipv6_cidr_block)
  lan_ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "64")

  ipv4_cidr_block = format("%s/%s", ionoscloud_server.test_server.nic[0].ips[0], "24")
  ipv6_cidr_block = format("%s/%s", local.lan_ipv6_cidr_block_parts[0], "80")
}

resource "ionoscloud_vpn_ipsec_gateway" "example" {
  name = "ipsec-gateway"
  location = "de/fra"
  gateway_ip = ionoscloud_ipblock.test_ipblock.ips[0]
  version = "IKEv2"
  description = "This gateway connects site A to VDC X."

  connections {
    datacenter_id = ionoscloud_datacenter.test_datacenter.id
    lan_id = ionoscloud_lan.test_lan.id
    ipv4_cidr = local.ipv4_cidr_block
    ipv6_cidr = local.ipv6_cidr_block
  }
}

resource "ionoscloud_vpn_ipsec_tunnel" "example" {
    location = "es/vit"
    gateway_id = ionoscloud_vpn_ipsec_gateway.example.id
    
    name = "example-tunnel"
    remote_host = "vpn.mycompany.com"
    description = "Allows local subnet X to connect to virtual network Y."
    
    auth {
        method = "PSK"
        psk_key = "X2wosbaw74M8hQGbK3jCCaEusR6CCFRa"
    }
    
    ike {
        diffie_hellman_group = "16-MODP4096"
        encryption_algorithm = "AES256"
        integrity_algorithm = "SHA256"
        lifetime             = 86400
    }
    
    esp {
        diffie_hellman_group = "16-MODP4096"
        encryption_algorithm = "AES256"
        integrity_algorithm = "SHA256"
        lifetime             = 3600
    }
    
    cloud_network_cidrs = [
        "0.0.0.0/0"
    ]
    
    peer_network_cidrs = [
        "1.2.3.4/32"
    ]
}
```

## Argument reference

* `name` - (Required)[string] The name of the IPSec Gateway Tunnel.
* `location` - (Required)[string] The location of the IPSec Gateway Tunnel. Supported locations: de/fra, de/txl, es/vit,
  gb/lhr, us/ewr, us/las, us/mci, fr/par
* `gateway_id` - (Required)[string] The ID of the IPSec Gateway that the tunnel belongs to.
* `description` - (Optional)[string] The human-readable description of your IPSec Gateway Tunnel.
* `remote_host` - (Required)[string] The remote peer host fully qualified domain name or public IPV4 IP to connect to.
* `ike` - (Required)[list] Settings for the initial security exchange phase. Minimum items: 1. Maximum items: 1.
    * `diffie_hellman_group` - (Optional)[string] The Diffie-Hellman Group to use for IPSec Encryption. Possible
      values: `15-MODP3072`, `16-MODP4096`, `19-ECP256`, `20-ECP384`, `21-ECP521`, `28-ECP256BP`, `29-ECP384BP`, `30-ECP512BP`.
      Default value: `16-MODP4096`.
    * `encryption_algorithm` - (Optional)[string] The encryption algorithm to use for IPSec Encryption. Possible
      values: `AES128`, `AES256`, `AES128-CTR`, `AES256-CTR`, `AES128-GCM-16`, `AES256-GCM-16`, `AES128-GCM-12`, `AES256-GCM-12`, `AES128-CCM-12`,
      `AES256-CCM-12`. Default value: `AES256`.
    * `integrity_algorithm` - (Optional)[string] The integrity algorithm to use for IPSec Encryption. Possible
      values: `SHA256`, `SHA384`, `SHA512`, `AES-XCBC`. Default value: `SHA256`.
    * `lifetime` - (Optional)[string] The phase lifetime in seconds. Minimum value: `3600`. Maximum value: `86400`.
      Default value: `86400`.
* `esp` - (Required)[list] Settings for the IPSec SA (ESP) phase. Minimum items: 1. Maximum items: 1.
    * `diffie_hellman_group` - (Optional)[string] The Diffie-Hellman Group to use for IPSec Encryption. Possible
      values: `15-MODP3072`, `16-MODP4096`, `19-ECP256`, `20-ECP384`, `21-ECP521`, `28-ECP256BP`, `29-ECP384BP`, `30-ECP512BP`.
      Default value: `16-MODP4096`.
    * `encryption_algorithm` - (Optional)[string] The encryption algorithm to use for IPSec Encryption. Possible
      values: `AES128`, `AES256`, `AES128-CTR`, `AES256-CTR`, `AES128-GCM-16`, `AES256-GCM-16`, `AES128-GCM-12`, `AES256-GCM-12`, `AES128-CCM-12`,
      `AES256-CCM-12`. Default value: `AES256`.
    * `integrity_algorithm` - (Optional)[string] The integrity algorithm to use for IPSec Encryption. Possible
      values: `SHA256`, `SHA384`, `SHA512`, `AES-XCBC`. Default value: `SHA256`.
    * `lifetime` - (Optional)[string] The phase lifetime in seconds. Minimum value: `3600`. Maximum value: `86400`.
      Default value: `86400`.
* `auth` - (Required)[string] Properties with all data needed to define IPSec Authentication. Minimum items: 1. Maximum
  items: 1.
    * `method` - (Optional)[string] The authentication method to use for IPSec Authentication. Possible values: `PSK`.
      Default value: `PSK`.
    * `psk_key` - (Optional)[string] The pre-shared key to use for IPSec Authentication. **Note**: Required if method is
      PSK.
* `cloud_network_cidrs` - (Required)[list] The network CIDRs on the "Left" side that are allowed to connect to the IPSec
  tunnel, i.e. the CIDRs within your IONOS Cloud LAN. Specify "0.0.0.0/0" or "::/0" for all addresses. Minimum items: 1.
  Maximum items: 20.
* `peer_network_cidrs` - (Required)[list] The network CIDRs on the "Right" side that are allowed to connect to the IPSec
  tunnel. Specify "0.0.0.0/0" or "::/0" for all addresses. Minimum items: 1. Maximum items: 20.

## Import

The resource can be imported using the `location`, `gateway_id` and `tunnel_id`, for example:

```
terraform import ionoscloud_vpn_ipsec_tunnel.example {location}:{gateway_id}:{tunnel_id}
```
