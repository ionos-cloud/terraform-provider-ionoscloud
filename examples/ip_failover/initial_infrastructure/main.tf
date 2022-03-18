terraform {
  required_version = "> 0.12.0"
  required_providers {
    ionoscloud = {
      source = "ionos-cloud/ionoscloud"
      version = "6.2.0"
    }
  }
}

provider "ionoscloud" {}

data "ionoscloud_image" "example" {
  type                  = "HDD"
  cloud_init            = "V1"
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
  size                  = 2
  name                  = "IP Block Example"
}

resource "ionoscloud_server" "example_A" {
  name                  = "Server A"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
  availability_zone     = "ZONE_1"
  cpu_family            = "AMD_OPTERON"
  image_name            = data.ionoscloud_image.example.id
  image_password        = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name                = "system"
    size                = 14
    disk_type           = "SSD"
  }
  nic {
    name                = "NIC A"
    lan                 = ionoscloud_lan.example.id
    dhcp                = true
    firewall_active     = true
    ips                 = [ ionoscloud_ipblock.example.ips[0] ]
  }
}


resource "ionoscloud_server" "example_B" {
  name                  = "Server B"
  datacenter_id         = ionoscloud_datacenter.example.id
  cores                 = 1
  ram                   = 1024
  availability_zone     = "ZONE_1"
  cpu_family            = "AMD_OPTERON"
  image_name            = data.ionoscloud_image.example.id
  image_password        = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name                = "system"
    size                = 14
    disk_type           = "SSD"
  }
  nic {
    name                = "NIC B"
    lan                 = ionoscloud_lan.example.id
    dhcp                = true
    firewall_active     = true
    ips                 = [ ionoscloud_ipblock.example.ips[1] ]
  }
}


resource "ionoscloud_ipfailover" "example" {
  depends_on            = [ ionoscloud_lan.example ]
  datacenter_id         = ionoscloud_datacenter.example.id
  lan_id                = ionoscloud_lan.example.id
  ip                    = ionoscloud_ipblock.example.ips[0]
  nicuuid               = ionoscloud_server.example_A.primary_nic
}