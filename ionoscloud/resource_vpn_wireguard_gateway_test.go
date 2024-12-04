//go:build vpn || all || wireguard

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWireguardGateway(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testWireguardGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: wireguardGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, nameAttribute, constant.WireGuardGatewayTestResource),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", "description"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", "192.168.1.100/24"),
					// todo ipv6 does not work yet
					// resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.IpBlockResource+"."+constant.IpBlockTestResource, "ips.0"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.lan_id", "ionoscloud_lan.lan_example", "id"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.ipv4_cidr", "192.168.1.108/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week", "Monday"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time", "09:00:00"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier", "STANDARD"),
				),
			},
			{
				Config: WireguardGwDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.datacenter_id", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.datacenter_id"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.lan_id", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.lan_id"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.ipv4_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.ipv4_cidr"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr"),
					//resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier"),
				),
			},
			{
				Config: WireguardGWDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr"),
					//resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier"),
				),
			},
			{
				Config:      WireguardGWDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      WireguardGWDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the WireGuard Gateway ID or name"),
			},
			{
				Config:      WireguardGWDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no VPN WireGuard Gateway found with the specified name"),
			},
			{
				Config: WireguardGWConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, nameAttribute, constant.WireGuardGatewayTestResource+"1"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", "description1"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", "192.168.1.101/24"),
					//resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.IpBlockResource+"."+constant.IpBlockTestResource, "ips.0"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.datacenter_example", "id"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.lan_id", "ionoscloud_lan.lan_example", "id"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.1.lan_id", "ionoscloud_lan.lan_example2", "id"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.ipv4_cidr", "192.168.1.109/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.1.ipv4_cidr", "192.168.1.110/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.day_of_the_week", "Tuesday"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "maintenance_window.0.time", "13:00:00"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "tier", "STANDARD_HA"),
				),
			},
		},
	})
}

func testWireguardGatewayDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).VPNClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.WireGuardGatewayResource {
			continue
		}
		ID := rs.Primary.ID
		location := rs.Primary.Attributes["location"]
		_, apiResponse, err := client.GetWireguardGatewayByID(ctx, ID, location)
		if err != nil {
			if !apiResponse.HttpNotFound() {
				return fmt.Errorf("an error occurred while checking the destruction of wireguard gateway with ID: %s, error: %w", ID, err)
			}
		} else {
			return fmt.Errorf("wireguard gateawy with ID: %s still exists", ID)
		}
	}
	return nil
}

const WireguardGwDataSourceMatchById = wireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  location = "de/fra"
  id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
}
`

const WireguardGWDataSourceMatchByName = wireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  location = "de/fra"
  name = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.name
}
`

const WireguardGWDataSourceInvalidBothIDAndName = wireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  location = "de/fra"
  id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
  name = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.name
}
`

const WireguardGWDataSourceInvalidNoIDNoName = wireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + ` {
  location = "de/fra"
}
`

const WireguardGWDataSourceWrongNameError = wireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + ` {
  location = "de/fra"
  name = "nonexistent"
}
`

const WireguardGWConfigUpdate = `
resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "de/fra"
}
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "de/fra"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource "ionoscloud_lan" "lan_example2" {
  name = "lan_example2"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource` + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = "` + constant.WireGuardGatewayTestResource + `1"
  location = "de/fra"
  description = "description1"
  gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  interface_ipv4_cidr =  "192.168.1.101/24"
  private_key = "0HpE4BNwGHabeaC4aY/GFxB6fBSc0d49Db0qAzRVSVc="
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.109/24"
  }
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example2.id
    ipv4_cidr       =  "192.168.1.110/24"
  }
  maintenance_window {
    day_of_the_week       = "Tuesday"
    time                  = "13:00:00"
  }
  tier = "STANDARD_HA"
}
`
