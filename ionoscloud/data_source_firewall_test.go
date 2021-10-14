package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFirewall(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewallCreateResources,
			},
			{
				Config: testAccDataSourceFirewallMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "name", "ionoscloud_firewall.webserver_http", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "protocol", "ionoscloud_firewall.webserver_http", "protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "source_mac", "ionoscloud_firewall.webserver_http", "source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "source_ip", "ionoscloud_firewall.webserver_http", "source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "target_ip", "ionoscloud_firewall.webserver_http", "target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "icmp_type", "ionoscloud_firewall.webserver_http", "icmp_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_id", "icmp_code", "ionoscloud_firewall.webserver_http", "icmp_code"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "name", "ionoscloud_firewall.webserver_http", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "protocol", "ionoscloud_firewall.webserver_http", "protocol"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "source_mac", "ionoscloud_firewall.webserver_http", "source_mac"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "source_ip", "ionoscloud_firewall.webserver_http", "source_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "target_ip", "ionoscloud_firewall.webserver_http", "target_ip"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "icmp_type", "ionoscloud_firewall.webserver_http", "icmp_type"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_firewall.test_firewall_name", "icmp_code", "ionoscloud_firewall.webserver_http", "icmp_code"),
				),
			},
		},
	})
}

const testAccDataSourceFirewallCreateResources = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "firewall-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 2
  name = "firewall_ipblock"
}

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "test_datasource"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
}
`

const testAccDataSourceFirewallMatchId = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "firewall-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 2
  name = "firewall_ipblock"
}

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "test_datasource"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
}

data "ionoscloud_firewall" "test_firewall_id" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  id = ionoscloud_firewall.webserver_http.id
}
`

const testAccDataSourceFirewallMatchName = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "firewall-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
	image_name ="ubuntu-16.04"
	image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 5
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
		firewall {
      protocol = "TCP"
      name = "SSH"
      port_range_start = 22
      port_range_end = 22
    }
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = 2
  dhcp = true
  firewall_active = true
  name = "updated"
}

resource "ionoscloud_ipblock" "ipblock" {
  location = ionoscloud_datacenter.foobar.location
  size = 2
  name = "firewall_ipblock"
}

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "ICMP"
  name = "test_datasource"
  source_mac = "00:0a:95:9d:68:16"
  source_ip = ionoscloud_ipblock.ipblock.ips[0]
  target_ip = ionoscloud_ipblock.ipblock.ips[1]
  icmp_type = 1
  icmp_code = 8
}

data "ionoscloud_firewall" "test_firewall_name" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  name	= "test_datasource"
}
`
