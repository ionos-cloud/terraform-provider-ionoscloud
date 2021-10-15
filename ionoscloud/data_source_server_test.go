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
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "name", "ionoscloud_server."+ServerTestResourceName, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "cores", "ionoscloud_server."+ServerTestResourceName, "cores"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "ram", "ionoscloud_server."+ServerTestResourceName, "ram"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "availability_zone", "ionoscloud_server."+ServerTestResourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "cpu_family", "ionoscloud_server."+ServerTestResourceName, "cpu_family"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "type", "ionoscloud_server."+ServerTestResourceName, "type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "volumes.0.name", "ionoscloud_server."+ServerTestResourceName, "volume.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "volumes.0.size", "ionoscloud_server."+ServerTestResourceName, "volume.0.size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "volumes.0.type", "ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "volumes.0.bus", "ionoscloud_server."+ServerTestResourceName, "volume.0.bus"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "volumes.0.availability_zone", "ionoscloud_server."+ServerTestResourceName, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.lan", "ionoscloud_server."+ServerTestResourceName, "nic.0.lan"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.name", "ionoscloud_server."+ServerTestResourceName, "nic.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.dhcp", "ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_active", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_type", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.ips.0", "ionoscloud_server."+ServerTestResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.ips.1", "ionoscloud_server."+ServerTestResourceName, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.protocol", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.name", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.port_range_start", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.port_range_end", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.source_mac", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.source_ip", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.target_ip", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceById, "nics.0.firewall_rules.0.type", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "name", "ionoscloud_server."+ServerTestResourceName, "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "cores", "ionoscloud_server."+ServerTestResourceName, "cores"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "ram", "ionoscloud_server."+ServerTestResourceName, "ram"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "availability_zone", "ionoscloud_server."+ServerTestResourceName, "availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "cpu_family", "ionoscloud_server."+ServerTestResourceName, "cpu_family"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "type", "ionoscloud_server."+ServerTestResourceName, "type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "volumes.0.name", "ionoscloud_server."+ServerTestResourceName, "volume.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "volumes.0.size", "ionoscloud_server."+ServerTestResourceName, "volume.0.size"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "volumes.0.type", "ionoscloud_server."+ServerTestResourceName, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "volumes.0.bus", "ionoscloud_server."+ServerTestResourceName, "volume.0.bus"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "volumes.0.availability_zone", "ionoscloud_server."+ServerTestResourceName, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.lan", "ionoscloud_server."+ServerTestResourceName, "nic.0.lan"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.name", "ionoscloud_server."+ServerTestResourceName, "nic.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.dhcp", "ionoscloud_server."+ServerTestResourceName, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_active", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_type", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.ips.0", "ionoscloud_server."+ServerTestResourceName, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.ips.1", "ionoscloud_server."+ServerTestResourceName, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.protocol", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.name", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.port_range_start", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.port_range_end", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.source_mac", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.source_ip", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.target_ip", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_server."+ServerTestDataSourceByName, "nics.0.firewall_rules.0.type", "ionoscloud_server."+ServerTestResourceName, "nic.0.firewall.0.type"),
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

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerTestResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerTestResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerTestResourceName + `"
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

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerTestResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerTestResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerTestResourceName + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}

data "ionoscloud_server" ` + ServerTestDataSourceById + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  id			= ionoscloud_server.` + ServerTestResourceName + `.id
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

resource "ionoscloud_server" ` + ServerTestResourceName + ` {
  name = "` + ServerTestResourceName + `"
  datacenter_id = ionoscloud_datacenter.foobar.id
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name ="Debian-10-cloud-init.qcow2"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  type = "ENTERPRISE"
  volume {
    name = "` + ServerTestResourceName + `"
    size = 5
    disk_type = "SSD Standard"
	backup_unit_id = ionoscloud_backup_unit.example.id
    user_data = "foo"
    bus = "VIRTIO"
    availability_zone = "ZONE_1"
}
  nic {
    lan = ionoscloud_lan.webserver_lan.id
    name = "` + ServerTestResourceName + `"
    dhcp = true
    firewall_active = true
	firewall_type = "BIDIRECTIONAL"
    ips            = [ ionoscloud_ipblock.webserver_ipblock.ips[0], ionoscloud_ipblock.webserver_ipblock.ips[1] ]
    firewall {
      protocol = "TCP"
      name = "` + ServerTestResourceName + `"
      port_range_start = 22
      port_range_end = 22
	  source_mac = "00:0a:95:9d:68:17"
	  source_ip = ionoscloud_ipblock.webserver_ipblock.ips[2]
	  target_ip = ionoscloud_ipblock.webserver_ipblock.ips[3]
	  type = "EGRESS"
    }
  }
}

data "ionoscloud_server" ` + ServerTestDataSourceByName + ` {
  datacenter_id = ionoscloud_datacenter.foobar.id
  name			= "` + ServerTestResourceName + `"
}
`
