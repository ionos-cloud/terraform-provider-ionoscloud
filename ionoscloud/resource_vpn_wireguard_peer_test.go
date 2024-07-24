//go:build vpn || all || wireguard

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

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
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0", "1.2.3.4/32"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key", "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI="),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host", "1.2.3.4"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port", "51820"),
				),
			},
			{
				Config: WireguardPeerDataSourceMatchById,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description"),
				),
			},
			{
				Config: WireguardPeerDataSourceMatchByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", constant.DataSource+"."+constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description"),
				),
			},
			{
				Config:      WireguardPeerDataSourceInvalidBothIDAndName,
				ExpectError: regexp.MustCompile("ID and name cannot be both specified at the same time"),
			},
			{
				Config:      WireguardPeerDataSourceInvalidNoIDNoName,
				ExpectError: regexp.MustCompile("please provide either the wireguard peer ID or name"),
			},
			{
				Config:      WireguardPeerDataSourceWrongNameError,
				ExpectError: regexp.MustCompile("no vpn wireguard peer found with the specified name"),
			},
			{
				Config: WireguardPeerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, nameAttribute, constant.WireGuardPeerTestResource+"1"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "description", "description1"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "allowed_ips.0", "1.2.3.5/32"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "public_key", "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI=1"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.host", "1.2.3.5"),
					resource.TestCheckResourceAttr(constant.WireGuardPeerResource+"."+constant.WireGuardPeerTestResource, "endpoint.0.port", "51821"),
				),
			},
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

const WireguardPeerConfig = wireguardGatewayConfig + `
resource` + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + `{
  name = "` + constant.WireGuardPeerTestResource + `"
  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
  description = "description"
  allowed_ips = [ "1.2.3.4/32" ]
  public_key = "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI="
  endpoint {
    host = "1.2.3.4"
    port = 51820
  }
}`

const WireguardPeerDataSourceMatchById = WireguardPeerConfig + `
` + constant.DataSource + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + ` {
	  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
	  id = ` + constant.WireGuardPeerResource + `.` + constant.WireGuardPeerTestResource + `.id
}`

const WireguardPeerDataSourceMatchByName = WireguardPeerConfig + `
` + constant.DataSource + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + ` {
	  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
	  name = ` + constant.WireGuardPeerResource + `.` + constant.WireGuardPeerTestResource + `.name
}`

const WireguardPeerDataSourceInvalidBothIDAndName = WireguardPeerConfig + `
` + constant.DataSource + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + ` {
	  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
	  id = ` + constant.WireGuardPeerResource + `.` + constant.WireGuardPeerTestResource + `.id
	  name = ` + constant.WireGuardPeerResource + `.` + constant.WireGuardPeerTestResource + `.name
}`

const WireguardPeerDataSourceInvalidNoIDNoName = WireguardPeerConfig + `	
` + constant.DataSource + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + ` {
  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
}`

const WireguardPeerDataSourceWrongNameError = WireguardPeerConfig + `
` + constant.DataSource + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + ` {
	  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
	  name = "wrong-name"
}`

const WireguardPeerConfigUpdate = wireguardGatewayConfig + `
resource` + ` ` + constant.WireGuardPeerResource + ` ` + constant.WireGuardPeerTestResource + `{
  name = "` + constant.WireGuardPeerTestResource + `1"
  gateway_id = ` + constant.WireGuardGatewayResource + `.` + constant.WireGuardGatewayTestResource + `.id
  description = "description1"
  allowed_ips = [ "1.2.3.5/32" ]
  public_key = "no8iaSEoqfbI6PVYsdEiUU5efYdtKX8VAhKity19MWI=1"
  endpoint {
    host = "1.2.3.5"
    port = 51821
  }
}`
