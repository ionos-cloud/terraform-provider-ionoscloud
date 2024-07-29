//go:build apigateway || all || gateway

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApiGateway_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApiGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApiGatewayConfig_basic,
			},
			{
				ResourceName:      "ionoscloud_apigateway.example",
				ImportStateIdFunc: testAccApiGatewayImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApiGatewayImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_apigateway" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
