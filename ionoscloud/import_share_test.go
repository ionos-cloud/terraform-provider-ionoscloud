//go:build compute || all || share

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccShareImportBasic(t *testing.T) {
	resourceName := "share"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckShareDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckShareConfigBasic),
			},

			{
				ResourceName:      fmt.Sprintf("%s.%s", constant.ShareResource, resourceName),
				ImportStateIdFunc: testAccShareImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccShareImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.ShareResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["resource_id"])
	}

	return importID, nil
}
