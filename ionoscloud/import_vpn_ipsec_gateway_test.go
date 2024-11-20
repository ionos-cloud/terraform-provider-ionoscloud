//go:build vpn || all || ipsec

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccIPSecGatewayImport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		ExternalProviders: randomProviderVersion343(),
		CheckDestroy:      testCheckIPSecGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: configIPSecGatewayBasic(GatewayResourceName, GatewayAttributeNameValue),
			},
			{
				ResourceName:      constant.IPSecGatewayResource + "." + GatewayResourceName,
				ImportStateIdFunc: testAccIPSecGatewayImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIPSecGatewayImportStateID(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.IPSecGatewayResource {
			continue
		}

		importID = fmt.Sprintf("%s:%s", rs.Primary.Attributes["location"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
