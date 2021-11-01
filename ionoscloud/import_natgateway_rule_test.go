// +build natgateway

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNatGatewayRuleImportBasic(t *testing.T) {
	natGatewayRuleName := "natGatewayRule"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNatGatewayRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckNatGatewayRuleConfigBasic, natGatewayRuleName),
			},

			{
				ResourceName:      "ionoscloud_natgateway_rule.natgateway_rule",
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
		if rs.Type != "ionoscloud_natgateway_rule" {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["natgateway_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
