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
				Config: testAccDataSourceBackupUnitCreateResources,
			},
			{
				Config: testAccDataSourceBackupUnitMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_id", "name", "ionoscloud_backup_unit.test_ds_backup_unit", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_id", "email", "ionoscloud_backup_unit.test_ds_backup_unit", "email"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_id", "login", "ionoscloud_backup_unit.test_ds_backup_unit", "login"),
				),
			},
			{
				Config: testAccDataSourceBackupUnitMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_name", "name", "ionoscloud_backup_unit.test_ds_backup_unit", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_name", "email", "ionoscloud_backup_unit.test_ds_backup_unit", "email"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_backup_unit.test_backup_unit_name", "login", "ionoscloud_backup_unit.test_ds_backup_unit", "login"),
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

data "ionoscloud_backup_unit" "test_backup_unit_id" {
  id			= ionoscloud_backup_unit.test_ds_backup_unit.id
}
`

const testAccDataSourceBackupUnitMatchName = `
resource "ionoscloud_backup_unit" "test_ds_backup_unit" {
	name        = "test ds backup unit"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}

data "ionoscloud_backup_unit" "test_backup_unit_name" {
  name			= "test ds backup unit"
}
`
