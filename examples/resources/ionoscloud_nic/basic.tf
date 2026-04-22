resource "ionoscloud_datacenter" "example" {
  name                = "Datacenter Example"
  location            = "us/las"
  description         = "Datacenter Description"
  sec_auth_protection = false
}

resource "ionoscloud_ipblock" "example" {
  location            = ionoscloud_datacenter.example.location
  size                = 2
  name                = "IP Block Example"
}

resource "ionoscloud_lan" "example"{
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "Lan"
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
  name                  = "NIC"
  dhcp                  = true
  firewall_active       = true
  firewall_type         = "INGRESS"
  ips                   = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
