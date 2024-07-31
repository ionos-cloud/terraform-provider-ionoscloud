//go:build cdn || all || distribution

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCdnDistributionImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCdnDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCdnDistributionConfigBasicImport,
			},
			{
				ResourceName:      constant.CdnDistributionResource + "." + constant.CdnDistributionTestResource,
				ImportStateIdFunc: testAccCdnDistributionImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCdnDistributionImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.CdnDistributionResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
