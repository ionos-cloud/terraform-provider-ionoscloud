resource "ionoscloud_datacenter" "example" {
  name       = "Resource Server Test"
  location = "us/las"
}
resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = "us/las"
  size = 4
  name = "webserver_ipblock"
}
resource "ionoscloud_lan" "example" {
  datacenter_id = ionoscloud_datacenter.example.id
  public = true
  name = "public"
  ipv6_cidr_block = cidrsubnet(ionoscloud_datacenter.example.ipv6_cidr_block,8,10)
}
resource "ionoscloud_server" "example" {
  name = "Resource Server Test"
  datacenter_id = ionoscloud_datacenter.example.id
  cores = 1
  ram = 1024
  image_name ="ubuntu:latest"
  image_password = random_password.server_image_password.result
  type = "ENTERPRISE"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD Standard"
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.example.id
    name = "system"
    dhcp = true
    firewall_active = true
    firewall_type = "BIDIRECTIONAL"
    ips = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]

    dhcpv6 = true
    ipv6_cidr_block = cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,24)
    ipv6_ips        = [
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,24),10),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,24),20),
                        cidrhost(cidrsubnet(ionoscloud_lan.example.ipv6_cidr_block,16,24),30)
                      ]

    firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    source_mac = "00:0a:95:9d:68:17"
    source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
    target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
    type = "EGRESS"
    }

  }
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
