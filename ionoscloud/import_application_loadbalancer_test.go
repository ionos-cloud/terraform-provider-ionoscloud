//go:build temporarilyDisabled

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccApplicationLoadBalancerImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckApplicationLoadBalancerDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckApplicationLoadBalancerConfigBasic,
			},

			{
				ResourceName:      resourceNameAlb,
				ImportStateIdFunc: testAccApplicationLoadBalancerImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccApplicationLoadBalancerImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ALBResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
