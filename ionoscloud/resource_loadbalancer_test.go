//go:build all || waiting_for_vdc
// +build all waiting_for_vdc

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLoadbalancerBasic(t *testing.T) {
	var loadbalancer ionoscloud.Loadbalancer
	lbName := "loadbalancer"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckLoadbalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckLoadbalancerConfigBasic, lbName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadbalancerExists("ionoscloud_loadbalancer.loadbalancer", &loadbalancer),
					testAccCheckLoadbalancerAttributes("ionoscloud_loadbalancer.loadbalancer", lbName),
					resource.TestCheckResourceAttr("ionoscloud_loadbalancer.loadbalancer", "name", lbName),
				),
			},
			{
				Config: testAccCheckLoadbalancerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadbalancerAttributes("ionoscloud_loadbalancer.loadbalancer", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_loadbalancer.loadbalancer", "name", "updated"),
				),
			},
		},
	})
}

func testAccCheckLoadbalancerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_loadbalancer" {
			continue
		}

		dcId := rs.Primary.Attributes["datacenter_id"]

		_, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking the destruction of load balancer %s: %s",
					rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("load balancer %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLoadbalancerAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckLoadbalancerAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["name"] != name {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["name"])
		}

		return nil
	}
}

func testAccCheckLoadbalancerExists(n string, loadbalancer *ionoscloud.Loadbalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckLoadbalancerExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}
		dcId := rs.Primary.Attributes["datacenter_id"]
		foundLB, apiResponse, err := client.LoadBalancersApi.DatacentersLoadbalancersFindById(ctx, dcId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Loadbalancer: %s", rs.Primary.ID)
		}
		if *foundLB.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		loadbalancer = &foundLB

		return nil
	}
}

const testAccCheckLoadbalancerConfigBasic = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "loadbalancer-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password_updated.result
  volume {
    name = "system"
    size = 14
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

resource "ionoscloud_nic" "database_nic" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = "3"
  dhcp = true
  firewall_active = true
  name = "updated"
  lifecycle {
    ignore_changes = [ lan ]
  }
}

resource "ionoscloud_loadbalancer" "loadbalancer" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  nic_ids = ["${ionoscloud_nic.database_nic.id}"]
  name = "%s"
  dhcp = true
}

resource ` + RandomPassword + ` "server_image_password_updated" {
  length           = 16
  special          = false
}
`

const testAccCheckLoadbalancerConfigUpdate = `
resource "ionoscloud_datacenter" "foobar" {
	name       = "loadbalancer-test"
	location = "us/las"
}

resource "ionoscloud_server" "webserver" {
  name = "webserver"
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  cores = 1
  ram = 1024
  availability_zone = "ZONE_1"
  cpu_family = "AMD_OPTERON"
  image_name = "ubuntu:latest"
  image_password = ` + RandomPassword + `.server_image_password.result
  volume {
    name = "system"
    size = 14
    disk_type = "SSD"
  }
  nic {
    lan = "1"
    dhcp = true
    firewall_active = true
  }
}

` + ServerImagePassword + `

resource "ionoscloud_nic" "database_nic2" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  server_id = "${ionoscloud_server.webserver.id}"
  lan = "3"
  dhcp = true
  firewall_active = true
  name = "updated"
  lifecycle {
    ignore_changes = [ lan ]
  }
}

resource "ionoscloud_loadbalancer" "loadbalancer" {
  datacenter_id = "${ionoscloud_datacenter.foobar.id}"
  nic_ids = ["${ionoscloud_nic.database_nic2.id}"]
  name = "updated"
  dhcp = true
}`
