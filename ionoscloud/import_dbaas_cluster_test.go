package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDbaasPgSqlClusterImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDbaasPgSqlClusterDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasic,
			},
			{
				ResourceName:            DBaaSClusterResource + "." + DBaaSClusterTestResource,
				ImportStateIdFunc:       testAccDbaasPgSqlClusterImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credentials"},
			},
		},
	})
}

func testAccDbaasPgSqlClusterImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != DBaaSClusterResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
