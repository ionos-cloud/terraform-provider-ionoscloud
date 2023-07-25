//go:build all || natgateway
// +build all natgateway

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGatewayImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayConfigBasic, constant.NatGatewayTestResource),
			},

			{
				ResourceName:      resourceNatGatewayResource,
				ImportStateIdFunc: testAccNatGatewayImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNatGatewayImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NatGatewayResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
