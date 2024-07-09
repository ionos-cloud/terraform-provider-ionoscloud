//go:build cdn || all || distribution

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDistributionImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDistributionConfigBasic,
			},

			{
				ResourceName:      constant.DistributionResource + "." + constant.DistributionTestResource,
				ImportStateIdFunc: testAccDistributionImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDistributionImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DistributionResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
