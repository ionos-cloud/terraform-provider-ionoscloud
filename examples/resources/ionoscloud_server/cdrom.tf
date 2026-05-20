resource "ionoscloud_datacenter" "cdrom" {
  name = "CDROM Test"
  location = "de/txl"
  description = "CDROM image test"
  sec_auth_protection = false
}

resource "ionoscloud_lan" "public" {
  datacenter_id = ionoscloud_datacenter.cdrom.id
  public = true
  name = "Uplink"
}

data "ionoscloud_image" "cdrom" {
  image_alias = "ubuntu:latest_iso"
  type        = "CDROM"
  location    = "de/txl"
  cloud_init  = "NONE"
}

resource "ionoscloud_server" "test" {
  datacenter_id  = ionoscloud_datacenter.cdrom.id
  name           = "ubuntu_latest_from_cdrom"
  cores          = 1
  ram            = 1024
  cpu_family     = ionoscloud_datacenter.cdrom.cpu_architecture[0].cpu_family
  type           = "ENTERPRISE"
  volume {
    name         = "hdd0"
    disk_type    = "HDD"
    size         = 50
    licence_type = "OTHER"
  }
  nic {
    lan    = 1
    dhcp   = true
    firewall_active = false
  }
}
