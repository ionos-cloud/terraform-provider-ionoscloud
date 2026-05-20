data "ionoscloud_template" "example" {
  name            = "Basic Cube XS"
}
resource "ionoscloud_datacenter" "example" {
	name            = "Datacenter Example"
	location        = "de/txl"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "de/txl"
  size = 4
  name = "webserver_ipblock"
}
resource "ionoscloud_lan" "example" {
  datacenter_id     = ionoscloud_datacenter.example.id
  public            = true
  name              = "Lan Example"
  ipv6_cidr_block = cidrsubnet(ionoscloud_datacenter.example.ipv6_cidr_block,8,10)
}
resource "ionoscloud_cube_server" "example" {
  name              = "Server Example"
  availability_zone = "AUTO"
  image_name        = "ubuntu:latest"
  template_uuid     = data.ionoscloud_template.example.id
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name            = "Volume Example"
    licence_type    = "LINUX"
    disk_type       = "DAS"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    ips             = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1]]

    dhcpv6          = false
    ipv6_cidr_block = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5)
    ipv6_ips        = [
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),1),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),2),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,5),3)
                      ]

    firewall_active = true
  }
}

resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
