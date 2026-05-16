resource "ionoscloud_server_boot_device_selection" "example"{
  datacenter_id  = ionoscloud_datacenter.example.id
  server_id      = ionoscloud_server.example.id
  boot_device_id = ionoscloud_volume.example.id
}

resource "ionoscloud_server" "example" {
  name              = "Server Example"
  availability_zone = "ZONE_2"
  image_name        = "ubuntu:latest"
  cores             = 2
  ram               = 2048
  image_password    = random_password.server_image_password.result
  datacenter_id     = ionoscloud_datacenter.example.id
  volume {
    name = "Inline Updated"
    size = 20
    disk_type = "SSD Standard"
    bus = "VIRTIO"
    availability_zone = "AUTO"
  }
  nic {
    lan             = ionoscloud_lan.example.id
    name            = "Nic Example"
    dhcp            = true
    firewall_active = true
  }
}

resource "ionoscloud_volume" "example" {
  server_id = ionoscloud_server.example.id
  datacenter_id     = ionoscloud_datacenter.example.id
  name = "External 1"
  size = 10
  disk_type = "HDD"
  availability_zone = "AUTO"
  image_name = "debian:latest"
  image_password = random_password.server_image_password.result
}
