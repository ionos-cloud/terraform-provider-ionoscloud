package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFirewall_matchId(t *testing.T) {
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
					resource.TestCheckResourceAttr("data.ionoscloud_firewall.test_firewall", "name", "test_ds_firewall_rule"),
				),
			},
		},
	})
}

func TestAccDataSourceFirewall_matchName(t *testing.T) {
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
				Config: testAccDataSourceFirewallMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_firewall.test_firewall", "name", "test_ds_firewall_rule"),
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

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "TCP"
  name = "test_ds_firewall_rule"
  port_range_start = 80
  port_range_end = 80
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

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "TCP"
  name = "test_ds_firewall_rule"
  port_range_start = 80
  port_range_end = 80
}
data "ionoscloud_firewall" "test_firewall" {
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

resource "ionoscloud_firewall" "webserver_http" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  protocol = "TCP"
  name = "test_ds_firewall_rule"
  port_range_start = 80
  port_range_end = 80
}
data "ionoscloud_firewall" "test_firewall" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  nic_id = "${ionoscloud_nic.database_nic.id}"
  name	= "test_ds"
}
`
