//go:build compute || all || backup

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBackupUnit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBackupUnitConfigBasic,
			},
			{
				Config: testAccDataSourceBackupUnitMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceById, "name", BackupUnitResource+"."+BackupUnitTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceById, "email", BackupUnitResource+"."+BackupUnitTestResource, "email"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceById, "login", BackupUnitResource+"."+BackupUnitTestResource, "login"),
				),
			},
			{
				Config: testAccDataSourceBackupUnitMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "name", BackupUnitResource+"."+BackupUnitTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "email", BackupUnitResource+"."+BackupUnitTestResource, "email"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "login", BackupUnitResource+"."+BackupUnitTestResource, "login"),
				),
			},
		},
	})
}

const testAccDataSourceBackupUnitMatchId = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceById + ` {
  id			= ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
}
`

const testAccDataSourceBackupUnitMatchName = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "` + BackupUnitTestResource + `"
}
`
