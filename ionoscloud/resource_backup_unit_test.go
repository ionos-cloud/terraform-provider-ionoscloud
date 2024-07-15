//go:build all || backup

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBackupUnitBasic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckBackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBackupUnitConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "name", constant.BackupUnitTestResource),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "email", "example@ionoscloud.com"),
					resource.TestCheckResourceAttrPair(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "password", constant.RandomPassword+".backup_unit_password", "result"),
				),
			},
			{
				Config: testAccDataSourceBackupUnitMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceById, "name", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceById, "email", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceById, "login", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "login"),
				),
			},
			{
				Config: testAccDataSourceBackupUnitMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceByName, "name", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceByName, "email", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.BackupUnitResource+"."+constant.BackupUnitDataSourceByName, "login", constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "login"),
				),
			},
			{
				Config:      testAccDataSourceBackupUnitMatchWrongNameError,
				ExpectError: regexp.MustCompile("no backup unit found with the specified name"),
			},
			{
				Config: testAccCheckBackupUnitConfigUpdatePassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "name", constant.BackupUnitTestResource),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "email", "example@ionoscloud.com"),
					resource.TestCheckResourceAttrPair(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "password", constant.RandomPassword+".backup_unit_password_updated", "result"),
				),
			},
			{
				Config: testAccCheckBackupUnitConfigUpdateEmail,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "name", constant.BackupUnitTestResource),
					resource.TestCheckResourceAttr(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "email", "example-updated@ionoscloud.com"),
					resource.TestCheckResourceAttrPair(constant.BackupUnitResource+"."+constant.BackupUnitTestResource, "password", constant.RandomPassword+".backup_unit_password_updated", "result"),
				),
			},
		},
	})
}

func testAccCheckBackupUnitDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {

		if rs.Type != constant.BackupUnitResource {
			continue
		}

		_, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking for the destruction of backup unit %s: %w",
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
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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

		foundBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching backup unit: %s", rs.Primary.ID)
		}
		if *foundBackupUnit.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}
		backupUnit = &foundBackupUnit

		return nil
	}
}

const testAccCheckBackupUnitConfigUpdatePassword = `
resource ` + constant.BackupUnitResource + ` ` + constant.BackupUnitTestResource + ` {
	name        = "` + constant.BackupUnitTestResource + `"
	password    = ` + constant.RandomPassword + `.backup_unit_password_updated.result
	email       = "example@ionoscloud.com"
}
resource ` + constant.RandomPassword + ` "backup_unit_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`

const testAccCheckBackupUnitConfigUpdateEmail = `
resource ` + constant.BackupUnitResource + ` ` + constant.BackupUnitTestResource + ` {
	name        = "` + constant.BackupUnitTestResource + `"
	password    = ` + constant.RandomPassword + `.backup_unit_password_updated.result
	email       = "example-updated@ionoscloud.com"
}
resource ` + constant.RandomPassword + ` "backup_unit_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}
`
const testAccDataSourceBackupUnitMatchId = testAccCheckBackupUnitConfigBasic + `
data ` + constant.BackupUnitResource + ` ` + constant.BackupUnitDataSourceById + ` {
  id			= ` + constant.BackupUnitResource + `.` + constant.BackupUnitTestResource + `.id
}
`

const testAccDataSourceBackupUnitMatchName = testAccCheckBackupUnitConfigBasic + `
resource ` + constant.BackupUnitResource + ` ` + constant.BackupUnitTestResource + `similar {
	name        = "similar` + constant.BackupUnitTestResource + `"
	password    = ` + constant.RandomPassword + `.backup_unit_password_updated.result
	email       = "example-updated@ionoscloud.com"
}
resource ` + constant.RandomPassword + ` "backup_unit_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

data ` + constant.BackupUnitResource + ` ` + constant.BackupUnitDataSourceByName + ` {
  name			= "` + constant.BackupUnitTestResource + `"
}
`

const testAccDataSourceBackupUnitMatchWrongNameError = testAccCheckBackupUnitConfigBasic + `
data ` + constant.BackupUnitResource + ` ` + constant.BackupUnitDataSourceByName + ` {
  name			= "wrong_name"
}
`
