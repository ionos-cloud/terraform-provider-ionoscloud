//go:build all || natgateway
// +build all natgateway

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGatewayRuleImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigBasic, NatGatewayRuleTestResource),
			},

			{
				ResourceName:      resourceNatGatewayRuleResource,
				ImportStateIdFunc: testAccNatGatewayRuleImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNatGatewayRuleImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NatGatewayRuleResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
