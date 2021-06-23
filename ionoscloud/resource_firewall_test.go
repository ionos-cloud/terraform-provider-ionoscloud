package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewall_Basic(t *testing.T) {
	var firewall ionoscloud.FirewallRule
	firewallName := "firewall"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckFirewallDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckfirewallconfigBasic, firewallName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallExists("ionoscloud_firewall.webserver_http", &firewall),
					testAccCheckFirewallAttributes("ionoscloud_firewall.webserver_http", firewallName),
					resource.TestCheckResourceAttr("ionoscloud_firewall.webserver_http", "name", firewallName),
				),
			},
			{
				Config: testacccheckfirewallconfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallAttributes("ionoscloud_firewall.webserver_http", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_firewall.webserver_http", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckFirewallDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_firewall" {
			continue
		}

		_, apiRsponse, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
			rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()

		if err != nil {
			if apiRsponse != nil && apiRsponse.StatusCode != 404 {
				return fmt.Errorf("error occurent at checking deletion of firewall %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("firewall still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckFirewallAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckFirewallAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckFirewallExists(n string, firewall *ionoscloud.FirewallRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckFirewallExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundServer, _, err := client.FirewallRulesApi.DatacentersServersNicsFirewallrulesFindById(ctx, rs.Primary.Attributes["datacenter_id"],
			rs.Primary.Attributes["server_id"], rs.Primary.Attributes["nic_id"], rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching Firewall rule: %s", rs.Primary.ID)
		}
		if *foundServer.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		firewall = &foundServer

		return nil
	}
}

const testacccheckfirewallconfigBasic = `
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
  image_name = "Ubuntu-20.04-LTS-server-2021-06-01"
  image_password = "K3tTj8G14a3EgKyNeeiY"
  volume {
    name = "system"
    size = 14
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
  name = "%s"
  port_range_start = 80
  port_range_end = 80
}
`

const testacccheckfirewallconfigUpdate = `
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
  image_name = "Ubuntu-20.04-LTS-server-2021-06-01"
  image_password = "test1234"
  volume {
    name = "system"
    size = 14
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
  name = "updated"
  port_range_start = 80
  port_range_end = 80
}
`
