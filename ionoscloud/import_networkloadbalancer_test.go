//go:build nlb
// +build nlb

package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkLoadBalancerImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkLoadBalancerConfigBasic,
			},

			{
				ResourceName:      networkLoadBalancerResource,
				ImportStateIdFunc: testAccNetworkLoadBalancerImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkLoadBalancerImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != NetworkLoadBalancerResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
