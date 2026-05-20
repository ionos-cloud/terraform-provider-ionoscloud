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
    location = "de/fra"
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
