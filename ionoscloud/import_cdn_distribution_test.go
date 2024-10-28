//go:build cdn || all || distribution

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCDNDistributionImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckCDNDistributionDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCDNDistributionConfigBasicImport,
			},
			{
				ResourceName:      constant.CDNDistributionResource + "." + constant.CDNDistributionTestResource,
				ImportStateIdFunc: testAccCDNDistributionImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCDNDistributionImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.CDNDistributionResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
