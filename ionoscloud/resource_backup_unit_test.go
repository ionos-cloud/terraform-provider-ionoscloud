package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccbackupUnit_Basic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit
	backupUnitName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckbackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckbackupUnitConfigBasic, backupUnitName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "name", backupUnitName),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example@profitbricks.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword123$"),
				),
			},
			{
				Config: testAccCheckbackupUnitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example-updated@profitbricks.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword1234$"),
				),
			},
		},
	})
}

func testAccCheckbackupUnitDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "ionoscloud_backup_unit" {
			continue
		}

		_, apiResponse, err := client.BackupUnitApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()

		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("backup unit still exists %s %s", rs.Primary.ID, string(apiResponse.Payload))
			}
		} else {
			return fmt.Errorf("Unable to fetch backup unit %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckbackupUnitExists(n string, backupUnit *ionoscloud.BackupUnit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundBackupUnit, _, err := client.BackupUnitApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching backup unit: %s", rs.Primary.ID)
		}
		if *foundBackupUnit.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		backupUnit = &foundBackupUnit

		return nil
	}
}

const testAccCheckbackupUnitConfigBasic = `
resource "ionoscloud_backup_unit" "example" {
	name        = "%s"
	password    = "DemoPassword123$"
	email       = "example@profitbricks.com"
}
`

const testAccCheckbackupUnitConfigUpdate = `
resource "ionoscloud_backup_unit" "example" {
	name        = "example"
	email       = "example-updated@profitbricks.com"
	password    = "DemoPassword1234$"
}
`
