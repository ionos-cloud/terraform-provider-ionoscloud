//go:build compute || all

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccNSGRuleImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNSGRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNSGFirewallRulesBasic,
			},
			{
				ResourceName:      constant.NSGFirewallRuleResource + "." + constant.NSGFirewallRuleTestResource + "_1",
				ImportStateIdFunc: testAccNSGRuleImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      constant.NSGFirewallRuleResource + "." + constant.NSGFirewallRuleTestResource + "_2",
				ImportStateIdFunc: testAccNSGRuleImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNSGRuleImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NSGFirewallRuleResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["nsg_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
