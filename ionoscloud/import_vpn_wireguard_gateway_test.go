//go:build vpn || all || wireguard

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccWireguardGatewayImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testWireguardGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: wireguardGatewayConfig,
			},
			{
				ResourceName:            constant.WireGuardGatewayResource + "." + constant.WireGuardGatewayTestResource,
				ImportState:             true,
				ImportStateIdFunc:       testAccWireguardGatewayImportStateID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func testAccWireguardGatewayImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.WireGuardGatewayResource {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
