//go:build vpn || all || wireguard

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
				Config: WireguardGatewayConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, nameAttribute, constant.WireGuardGatewayTestResource),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", "description"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", "192.168.1.100/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.IpBlockResource+"."+constant.IpBlockTestResource, "ips.0"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.datacenter_id", "ionoscloud_datacenter.datacenter_example", "id"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.lan_id", "ionoscloud_lan.lan_example.id"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "connections.0.ipv4_cidr", "192.168.1.108/24"),
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
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr"),
				),
			},
			{
				Config: WireguardGWDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr"),
				),
			},
			{
				Config:      WireguardGWDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      WireguardGWDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the wireguard gateway ID or name"),
			},
			{
				Config:      WireguardGWDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no vpn wireguard gateway found with the specified name"),
			},
			{
				Config: WireguardGWConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, nameAttribute, constant.WireGuardGatewayTestResource+"1"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "description", "description1"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv4_cidr", "192.168.1.101/24"),
					resource.TestCheckResourceAttr(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
					resource.TestCheckResourceAttrPair(constant.WireGuardGatewayResource+"."+constant.WireGuardGatewayTestResource, "gateway_ip", constant.IpBlockResource+"."+constant.IpBlockTestResource, "ips.0"),
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
		_, apiResponse, err := client.GetWireguardGatewayByID(ctx, ID)
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

//func testAccWireguardGatewayExistenceCheck(path string, pipeline *logging.Pipeline) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		client := testAccProvider.Meta().(services.SdkBundle).VPNClient
//		rs, ok := s.RootModule().Resources[path]
//
//		if !ok {
//			return fmt.Errorf("not found: %s", path)
//		}
//		if rs.Primary.ID == "" {
//			return fmt.Errorf("no ID is set for the Logging pipeline")
//		}
//		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
//		defer cancel()
//		ID := rs.Primary.ID
//		pipelineResponse, _, err := client.GetWireguardGatewayByID(ctx, ID)
//		if err != nil {
//			return fmt.Errorf("an error occurred while fetching Logging pipeline with ID: %s, error: %w", ID, err)
//		}
//		pipeline = &pipelineResponse
//		return nil
//	}
//}

const WireguardGatewayConfig = `
resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "es/vit"
}
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "es/vit"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource` + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = "` + constant.WireGuardGatewayTestResource + `"
  description = "description"
  private_key = "private"
  gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  interface_ipv4_cidr =  "192.168.1.100/24"
  interface_ipv6_cidr = "2001:0db8:85a3::/24"
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.108/24"
  }
}`

const WireguardGwDataSourceMatchById = WireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
}
`

const WireguardGWDataSourceMatchByName = WireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.name
}
`

const WireguardGWDataSourceInvalidBothIDAndName = WireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
	id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
	name = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.name
}
`

const WireguardGWDataSourceInvalidNoIDNoName = WireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + ` {
}
`

const WireguardGWDataSourceWrongNameError = WireguardGatewayConfig + `
` + constant.DataSource + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + ` {
  name = "nonexistent"
}
`

const WireguardGWConfigUpdate = `
resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "es/vit"
}
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "es/vit"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource` + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = "` + constant.WireGuardGatewayTestResource + `"
  description = "description"
  private_key = "private"resource "ionoscloud_datacenter" "datacenter_example" {
  name = "datacenter_example"
  location = "es/vit"
}
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "es/vit"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}

resource "ionoscloud_lan" "lan_example" {
  name = "lan_example"
  datacenter_id = ionoscloud_datacenter.datacenter_example.id
}

resource` + ` ` + constant.WireGuardGatewayResource + ` ` + constant.WireGuardGatewayTestResource + `{
  name = "` + constant.WireGuardGatewayTestResource + `1"
  description = "description1"
  gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  interface_ipv4_cidr =  "192.168.1.100/24"
  interface_ipv6_cidr = "2001:0db8:85a3::/24"
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.108/24"
  }
}
  gateway_ip = ` + constant.IpBlockResource + `.` + constant.IpBlockTestResource + `.ips[0]
  interface_ipv4_cidr =  "192.168.1.101/24"
  interface_ipv6_cidr = "2001:0db8:85a4::/24"
  connections   {
    datacenter_id   =  ionoscloud_datacenter.datacenter_example.id
    lan_id          =  ionoscloud_lan.lan_example.id
    ipv4_cidr       =  "192.168.1.109/24"
  }
}
`
