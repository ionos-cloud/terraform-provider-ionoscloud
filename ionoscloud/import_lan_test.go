//go:build compute || all || lan

package ionoscloud

import (
	"fmt"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLanImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckLanDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLanConfigBasic,
			},

			{
				ResourceName:      constant.LanResource + "." + constant.LanTestResource,
				ImportStateIdFunc: testAccLanImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccLanImportStateId(s *terraform.State) (string, error) {
	var importID = ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.LanResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["datacenter_id"], rs.Primary.Attributes["id"])
	}

	return importID, nil
}
