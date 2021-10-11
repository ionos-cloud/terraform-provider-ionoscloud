package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceServerCreateResources,
			},
			{
				Config: testAccDataSourceServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "name", "ionoscloud_server."+ServerResourceName, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "cores", "ionoscloud_server."+ServerResourceName, "cores"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "ram", "ionoscloud_server."+ServerResourceName, "ram"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "availability_zone", "ionoscloud_server."+ServerResourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "cpu_family", "ionoscloud_server."+ServerResourceName, "cpu_family"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "type", "ionoscloud_server."+ServerResourceName, "type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "volume.0.name", "ionoscloud_server."+ServerResourceName, "volume.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "volume.0.size", "ionoscloud_server."+ServerResourceName, "volume.0.size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "volume.0.disk_type", "ionoscloud_server."+ServerResourceName, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "volume.0.bus", "ionoscloud_server."+ServerResourceName, "volume.0.bus"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "volume.0.availability_zone", "ionoscloud_server."+ServerResourceName, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.lan", "ionoscloud_server."+ServerResourceName, "nic.0.lan"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.name", "ionoscloud_server."+ServerResourceName, "nic.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.dhcp", "ionoscloud_server."+ServerResourceName, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall_active", "ionoscloud_server."+ServerResourceName, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall_type", "ionoscloud_server."+ServerResourceName, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.ips.0", "ionoscloud_server."+ServerResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.ips.1", "ionoscloud_server."+ServerResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.protocol", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.name", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.port_range_start", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.port_range_end", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.source_mac", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.source_ip", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.target_ip", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceById, "nic.0.firewall.0.type", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "name", "ionoscloud_server."+ServerResourceName, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "cores", "ionoscloud_server."+ServerResourceName, "cores"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "ram", "ionoscloud_server."+ServerResourceName, "ram"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "availability_zone", "ionoscloud_server."+ServerResourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "cpu_family", "ionoscloud_server."+ServerResourceName, "cpu_family"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "image_password", "ionoscloud_server."+ServerResourceName, "image_password"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "type", "ionoscloud_server."+ServerResourceName, "type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "volume.0.name", "ionoscloud_server."+ServerResourceName, "volume.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "volume.0.size", "ionoscloud_server."+ServerResourceName, "volume.0.size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "volume.0.disk_type", "ionoscloud_server."+ServerResourceName, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "volume.0.bus", "ionoscloud_server."+ServerResourceName, "volume.0.bus"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "volume.0.availability_zone", "ionoscloud_server."+ServerResourceName, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.lan", "ionoscloud_server."+ServerResourceName, "nic.0.lan"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.name", "ionoscloud_server."+ServerResourceName, "nic.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.dhcp", "ionoscloud_server."+ServerResourceName, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall_active", "ionoscloud_server."+ServerResourceName, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall_type", "ionoscloud_server."+ServerResourceName, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.ips.0", "ionoscloud_server."+ServerResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.ips.1", "ionoscloud_server."+ServerResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.protocol", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.name", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.port_range_start", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.port_range_end", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.source_mac", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.source_ip", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.target_ip", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerDataSourceByName, "nic.0.firewall.0.type", "ionoscloud_server."+ServerResourceName, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: "/* intentionally left blank to ensure proper datasource removal from the plan */",
			},
		},
	})
}

const testAccDataSourceServerCreateResources = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerResourceName + ` {
  name = "` + ServerResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerResourceName + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}
`

const testAccDataSourceServerMatchId = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerResourceName + ` {
  name = "` + ServerResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerResourceName + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}

data "ionoscloud_server" ` + ServerDataSourceById + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  id			= ionoscloud_server.` + ServerResourceName + `.id
}
`

const testAccDataSourceServerMatchName = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "server-test"
	location = "us/las"
}

resource "ionoscloud_backup_unit" "example" {
	name        = "serverTest"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

resource "ionoscloud_ipblock" "webserver_ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 4
  name = "webserver_ipblock"
}

resource "ionoscloud_lan" "webserver_lan" {
  datacenter_id = ionoscloud_datacenter.foobar.id
  public = true
  name = "public"
}

resource "ionoscloud_server" ` + ServerResourceName + ` {
  name = "` + ServerResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerResourceName + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}

data "ionoscloud_server" ` + ServerDataSourceByName + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  name			= "` + ServerResourceName + `"
}
`
