package ionoscloud

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroup_ImportBasic(t *testing.T) {
	resourceName := "resource_group"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckgroupconfigBasic, resourceName),
			},

			{
				ResourceName:      fmt.Sprintf("ionoscloud_group.%s", resourceName),
				ImportStateIdFunc: testAccGroupImportStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_group" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
