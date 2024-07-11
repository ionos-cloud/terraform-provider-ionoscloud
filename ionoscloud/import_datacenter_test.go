//go:build compute || all || datacenter

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataCenterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckDatacenterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatacenterConfigBasic,
			},

			{
				ResourceName:      constant.DatacenterResource + "." + constant.DatacenterTestResource,
				ImportStateIdFunc: testAccDatacenterImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDatacenterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.DatacenterResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
