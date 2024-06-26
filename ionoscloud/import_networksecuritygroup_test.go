package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccNetworkSecurityGroupImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ExternalProviders: randomProviderVersion343(),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkSecurityGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNetworkSecurityGroupConfigBasic,
			},

			{
				ResourceName:            constant.NetworkSecurityGroupResource + "." + constant.NetworkSecurityGroupTestResource,
				ImportStateIdFunc:       testAccNetworkSecurityGroupImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rule_ids"},
			},
		},
	})
}

func testAccNetworkSecurityGroupImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.NetworkSecurityGroupResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
