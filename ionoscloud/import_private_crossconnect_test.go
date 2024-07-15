//go:build compute || all || pcc

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPrivateCrossConnectImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckPrivateCrossConnectDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckPrivateCrossConnectConfigBasic,
			},
			{
				ResourceName:      constant.PCCResource + "." + constant.PCCTestResource,
				ImportStateIdFunc: testAccPCCImportStateID,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPCCImportStateID(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.PCCResource {
			continue
		}

		importID = rs.Primary.ID
	}

	return importID, nil
}
