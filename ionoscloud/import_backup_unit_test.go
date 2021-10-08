package ionoscloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackupUnitImportBasic(t *testing.T) {
	resourceName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckBackupUnitConfigBasic, resourceName),
			},
			{
				ResourceName:            fmt.Sprintf("ionoscloud_backup_unit.%s", resourceName),
				ImportStateIdFunc:       testAccBackupUnitImportStateId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccBackupUnitImportStateId(s *terraform.State) (string, error) {
	importID := ""

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_backup_unit" {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
