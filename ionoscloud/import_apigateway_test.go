//go:build apigateway || all || gateway

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccAPIGateway_import(t *testing.T) {
	resource.Test(
		t, resource.TestCase{
			PreCheck:                 func() { testAccPreCheck(t) },
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
			CheckDestroy:             testAccCheckAPIGatewayDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccCheckAPIGatewayConfig_basic,
				},
				{
					ResourceName:      "ionoscloud_apigateway.example",
					ImportStateIdFunc: testAccAPIGatewayImportStateId,
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		},
	)
}

func testAccAPIGatewayImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_apigateway" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
