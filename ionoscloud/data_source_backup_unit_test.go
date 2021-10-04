package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBackupUnit_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBackupUnitCreateResources,
			},
			{
				Config: testAccDataSourceBackupUnitMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_backup_unit.test_backup_unit", "name", "test ds backup unit"),
				),
			},
		},
	})
}

func TestAccDataSourceBackupUnit_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBackupUnitCreateResources,
			},
			{
				Config: testAccDataSourceBackupUnitMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_backup_unit.test_backup_unit", "name", "test ds backup unit"),
				),
			},
		},
	})

}

const testAccDataSourceBackupUnitCreateResources = `
resource "ionoscloud_backup_unit" "test_ds_backup_unit" {
	name        = "test ds backup unit"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
`

const testAccDataSourceBackupUnitMatchId = `
resource "ionoscloud_backup_unit" "test_ds_backup_unit" {
	name        = "test ds backup unit"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
data "ionoscloud_backup_unit" "test_backup_unit" {
  id			= ionoscloud_backup_unit.test_ds_backup_unit.id
}
`

const testAccDataSourceBackupUnitMatchName = `
resource "ionoscloud_backup_unit" "test_ds_backup_unit" {
	name        = "test ds backup unit"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
data "ionoscloud_backup_unit" "test_backup_unit" {
  name			= "test ds backup unit"
}
`
