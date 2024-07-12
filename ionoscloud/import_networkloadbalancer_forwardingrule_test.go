//go:build nlb
// +build nlb

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkLoadBalancerForwardingRuleImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckNetworkLoadBalancerForwardingRuleDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerForwardingRuleConfigBasic,
			},

			{
				ResourceName:      networkLoadBalancerForwardingRuleResource,
				ImportStateIdFunc: testAccNetworkLoadBalancerForwardingRuleImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkLoadBalancerForwardingRuleImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NetworkLoadBalancerForwardingRuleResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["networkloadbalancer_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
