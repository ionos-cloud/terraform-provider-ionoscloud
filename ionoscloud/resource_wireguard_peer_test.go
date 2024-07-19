//go:build vpn || all || wireguard

package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccWireguardPeer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testWireguardPeerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: WireguardPeerConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, nameAttribute, constant.WireGuardPeerTestResource),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", "description"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0", "1.2.3.4/24"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key", "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI="),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host", "1.2.3.4"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port", "51820"),
				),
			},
			//{
			//	Config: WireguardGwDataSourceMatchById,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.datacenter_id", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.datacenter_id"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.lan_id", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.lan_id"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.ipv4_cidr", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "connections.0.ipv4_cidr"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "gateway_ip", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "gateway_ip"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv4_cidr", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv4_cidr"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv6_cidr"),
			//	),
			//},
			//{
			//	Config: WireguardGWDataSourceMatchByName,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "gateway_ip", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "gateway_ip"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv4_cidr", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv4_cidr"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv6_cidr", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv6_cidr"),
			//	),
			//},
			//{
			//	Config:      WireguardGWDataSourceInvalidBothIDAndName,
			//	ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			//},
			//{
			//	Config:      WireguardGWDataSourceInvalidNoIDNoName,
			//	ExpectError: regexp.MustCompile("please provide either the wireguard gateway ID or name"),
			//},
			//{
			//	Config:      WireguardGWDataSourceWrongNameError,
			//	ExpectError: regexp.MustCompile("no vpn wireguard gateway found with the specified name"),
			//},
			//{
			//	Config: WireguardGWConfigUpdate,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, nameAttribute, constant.WireGuardPeerTestResource+"1"),
			//		resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", "description1"),
			//		resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv4_cidr", "192.168.1.101/24"),
			//		resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "interface_ipv6_cidr", "2001:0db8:85a3::/24"),
			//		resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "gateway_ip", constant.IpBlockResource+"."+constant.IpBlockTestResource, "ips.0"),
			//	),
			//},
		},
	})
}

func testWireguardPeerDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).VPNClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.WireGuardPeerResource {
			continue
		}
		ID := rs.Primary.ID
		gatewayID := rs.Primary.Attributes["gateway_id"]
		_, apiResponse, err := client.GetWireguardPeerByID(ctx, gatewayID, ID)
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

const WireguardPeerConfig = WireguardGatewayConfig + `
resource` + ` ` + constant.WireGuardPeerTestResource + ` ` + constant.WireGuardPeerTestResource + `{
  name = "` + constant.WireGuardPeerTestResource + `"
  gateway_ID = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
  description = "description"
  allowed_ips = [ "1.2.3.4/32" ]
  "publicKey": "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI="
  "endpoint": {
    "host": "1.2.3.4",
    "port": 51820
  }
}`
