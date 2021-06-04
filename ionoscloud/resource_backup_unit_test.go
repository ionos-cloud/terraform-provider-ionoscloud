package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccbackupUnit_Basic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit
	backupUnitName := "example"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckbackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckbackupUnitConfigBasic, backupUnitName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "name", backupUnitName),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example@ionoscloud.com"),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "password", "DemoPassword123$"),
				),
			},
			{
				Config: testAccCheckbackupUnitConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckbackupUnitExists("ionoscloud_backup_unit.example", &backupUnit),
					resource.TestCheckResourceAttr("ionoscloud_backup_unit.example", "email", "example-updated@ionoscloud.com"),
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

		if err != nil {
			if apiResponse != nil && apiResponse.StatusCode != 404 {
				payload := fmt.Sprintf("API response: %s", string(apiResponse.Payload))
				return fmt.Errorf("backup unit still exists %s - an error occurred while checking it %s %s", rs.Primary.ID, err, payload)
			}
		} else {
			return fmt.Errorf("backup unit still exists %s", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckbackupUnitExists(n string, backupUnit *ionoscloud.BackupUnit) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*ionoscloud.APIClient)

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

		foundBackupUnit, apiResponse, err := client.BackupUnitApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			payload := ""
			if apiResponse != nil {
				payload = fmt.Sprintf("API response: %s", string(apiResponse.Payload))
			}
			return fmt.Errorf("error occured while fetching backup unit: %s %s", rs.Primary.ID, payload)
		}
		if *foundBackupUnit.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		backupUnit = &foundBackupUnit

		return nil
	}
}

const testAccCheckbackupUnitConfigBasic = `
resource "ionoscloud_backup_unit" "example" {
	name        = "%s"
	password    = "DemoPassword123$"
	email       = "example@ionoscloud.com"
}
`

const testAccCheckbackupUnitConfigUpdate = `
resource "ionoscloud_backup_unit" "example" {
	name        = "example"
	email       = "example-updated@ionoscloud.com"
	password    = "DemoPassword1234$"
}
`
