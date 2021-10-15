package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBackupUnitBasic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit
	backupUnitName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckBackupUnitConfigBasic, backupUnitName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "name", backupUnitName),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example@profitbricks.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword123$"),
				),
			},
			{
				Config: testAccCheckBackupUnitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "name", "example"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example-updated@ionoscloud.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword1234$Updated"),
				),
			},
		},
	})
}

func testAccCheckBackupUnitDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "ionoscloud_backup_unit" {
			continue
		}

		_, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking for the destruction of backup unit %s: %s",
					rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("backup unit %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckBackupUnitExists(n string, backupUnit *ionoscloud.BackupUnit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundBackupUnit, _, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching backup unit: %s", rs.Primary.ID)
		}
		if *foundBackupUnit.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		backupUnit = &foundBackupUnit

		return nil
	}
}

const testAccCheckBackupUnitConfigBasic = `
resource "ionoscloud_backup_unit" "example" {
	name        = "%s"
	password    = "DemoPassword123$"
	email       = "example@profitbricks.com"
}
`

const testAccCheckBackupUnitConfigUpdate = `
resource "ionoscloud_backup_unit" "example" {
	name        = "example"
	email       = "example-updated@ionoscloud.com"
	password    = "DemoPassword1234$Updated"
}
`
