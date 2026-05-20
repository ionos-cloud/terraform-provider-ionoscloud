resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "IPv6 Enabled LAN"
  ipv6_cidr_block   = cidrsubnet(ionoscloud_datacenter.example.ipv6_cidr_block,8,2)
}

resource "ionoscloud_server" "example" {
  name                  = "Server Example"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
  image_name            = "Ubuntu-20.04"
  image_password        = random_password.server_image_password.result
  volume {
    name                = "system"
    size                = 14
    disk_type           = "SSD"
  }
  nic {
    lan                 = "1"
    dhcp                = true
    firewall_active     = true
  }
}

resource "ionoscloud_nic" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  server_id             = ionoscloud_server.example.id
  lan                   = ionoscloud_lan.example.id
  name                  = "IPv6 Enabled NIC"
  dhcp                  = true
  firewall_active       = true
  firewall_type         = "INGRESS"
  dhcpv6                = false
  ipv6_cidr_block       = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14)
  ipv6_ips              = [
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),10),
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),20),
                              cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,14),30)
                          ]
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
