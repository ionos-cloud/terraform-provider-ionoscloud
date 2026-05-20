resource "ionoscloud_nic" "example" {
  datacenter_id         = ionoscloud_datacenter.example.id
  server_id             = ionoscloud_server.example.id
  lan                   = ionoscloud_lan.example.id
  name                  = "IPV6 and Flowlog Enabled NIC"
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
  flowlog {
    action    = "ACCEPTED"
    bucket    = "flowlog-bucket"
    direction = "INGRESS"
    name      = "flowlog"
  }
}
