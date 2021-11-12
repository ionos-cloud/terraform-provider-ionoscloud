package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccShareImportBasic(t *testing.T) {
	resourceName := "share"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckShareDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckShareConfigBasic),
			},

			{
				ResourceName:      fmt.Sprintf("%s.%s", ShareResource, resourceName),
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
		if rs.Type != ShareResource {
			continue
		}

		importID = fmt.Sprintf("%s/%s", rs.Primary.Attributes["group_id"], rs.Primary.Attributes["resource_id"])
	}

	return importID, nil
}
