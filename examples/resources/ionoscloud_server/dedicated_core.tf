data "ionoscloud_image" "example" {
    type                  = "HDD"
    cloud_init            = "V1"
    image_alias           = "ubuntu:latest"
    location              = "us/las"
}

resource "ionoscloud_datacenter" "example" {
    name                  = "Datacenter Example"
    location              = "us/las"
    description           = "Datacenter Description"
    sec_auth_protection   = false
}

resource "ionoscloud_lan" "example" {
    datacenter_id         = ionoscloud_datacenter.example.id
    public                = true
    name                  = "Lan Example"
}

resource "ionoscloud_ipblock" "example" {
    location              = ionoscloud_datacenter.example.location
    size                  = 4
    name                  = "IP Block Example"
}

resource "ionoscloud_server" "example" {
    name                  = "Server Example"
    datacenter_id         = ionoscloud_datacenter.example.id
    cores                 = 1
    ram                   = 1024
    image_name            = data.ionoscloud_image.example.name
    image_password        = random_password.server_image_password.result
    type                  = "ENTERPRISE"
    volume {
        name              = "system"
        size              = 5
        disk_type         = "SSD Standard"
        user_data         = "foo"
        bus               = "VIRTIO"
        availability_zone = "ZONE_1"
    }
    nic {
        lan               = ionoscloud_lan.example.id
        name              = "system"
        dhcp              = true
        firewall_active   = true
        firewall_type     = "BIDIRECTIONAL"
        ips               = [ ionoscloud_ipblock.example.ips[0], ionoscloud_ipblock.example.ips[1] ]
        firewall {
          protocol          = "TCP"
          name              = "SSH"
          port_range_start  = 22
          port_range_end    = 22
          source_mac        = "00:0a:95:9d:68:17"
          source_ip         = ionoscloud_ipblock.example.ips[2]
          target_ip         = ionoscloud_ipblock.example.ips[3]
          type              = "EGRESS"
        }
    }
    label {
        key = "labelkey1"
        value = "labelvalue1"
    }
    label {
        key = "labelkey2"
        value = "labelvalue2"
    }
}
resource "random_password" "server_image_password" {
  length           = 16
  special          = false
}
