//go:build all || backup_unit

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBackupUnitImportBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBackupUnitConfigBasic,
			},
			{
				ResourceName:            BackupUnitResource + "." + BackupUnitTestResource,
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
		if rs.Type != BackupUnitResource {
			continue
		}

		importID = rs.Primary.Attributes["id"]
	}

	return importID, nil
}
