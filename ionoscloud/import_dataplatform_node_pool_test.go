//go:build all || dataplatform

package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataplatformNodePoolImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDataplatformNodePoolDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataplatformNodePoolConfigBasic,
			},
			{
				ResourceName:            constant.ResourceNameDataplatformNodePool,
				ImportStateIdFunc:       testAccDataplatformNodePoolImportStateID,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"maintenance_window.0.time"},
			},
		},
	})
}

func testAccDataplatformNodePoolImportStateID(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DataplatformNodePoolResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["cluster_id"], rs.Primary.ID)
	}
	return importID, nil
}
