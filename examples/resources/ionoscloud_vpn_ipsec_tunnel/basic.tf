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
