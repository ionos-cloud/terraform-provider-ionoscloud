//go:build compute || all || backup

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBackupUnitBasic(t *testing.T) {
	var backupUnit ionoscloud.BackupUnit

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckBackupUnitDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckBackupUnitConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(BackupUnitResource+"."+BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "name", BackupUnitTestResource),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "email", "example@ionoscloud.com"),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "password", "DemoPassword123$"),
				),
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
				Config: testAccDataSourceBackupUnitPartialMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "name", BackupUnitResource+"."+BackupUnitTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "email", BackupUnitResource+"."+BackupUnitTestResource, "email"),
					resource.TestCheckResourceAttrPair(DataSource+"."+BackupUnitResource+"."+BackupUnitDataSourceByName, "login", BackupUnitResource+"."+BackupUnitTestResource, "login"),
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
			{
				Config:      testAccDataSourceBackupUnitMatchWrongNameError,
				ExpectError: regexp.MustCompile("no backup unit found with the specified name"),
			},
			{
				Config:      testAccDataSourceBackupUnitPartialMatchWrongNameError,
				ExpectError: regexp.MustCompile("no backup unit found with the specified name"),
			},
			{
				Config:      testAccDataSourceBackupUnitMatchPartialMatchDisabledError,
				ExpectError: regexp.MustCompile("no backup unit found with the specified name"),
			},
			{
				Config: testAccCheckBackupUnitConfigUpdatePassword,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(BackupUnitResource+"."+BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "name", BackupUnitTestResource),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "email", "example@ionoscloud.com"),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "password", "DemoPassword1234$Updated"),
				),
			},
			{
				Config: testAccCheckBackupUnitConfigUpdateEmail,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupUnitExists(BackupUnitResource+"."+BackupUnitTestResource, &backupUnit),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "name", BackupUnitTestResource),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "email", "example-updated@ionoscloud.com"),
					resource.TestCheckResourceAttr(BackupUnitResource+"."+BackupUnitTestResource, "password", "DemoPassword1234$Updated"),
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

		if rs.Type != BackupUnitResource {
			continue
		}

		_, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

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

		foundBackupUnit, apiResponse, err := client.BackupUnitsApi.BackupunitsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

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

const testAccCheckBackupUnitConfigUpdatePassword = `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + ` {
	name        = "` + BackupUnitTestResource + `"
	password    = "DemoPassword1234$Updated"
	email       = "example@ionoscloud.com"
}
`

const testAccCheckBackupUnitConfigUpdateEmail = `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + ` {
	name        = "` + BackupUnitTestResource + `"
	password    = "DemoPassword1234$Updated"
	email       = "example-updated@ionoscloud.com"
}
`
const testAccDataSourceBackupUnitMatchId = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceById + ` {
  id			= ` + BackupUnitResource + `.` + BackupUnitTestResource + `.id
}
`

const testAccDataSourceBackupUnitPartialMatchName = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "` + DataSourcePartial + `"
  partial_match = true
}
`

const testAccDataSourceBackupUnitMatchName = testAccCheckBackupUnitConfigBasic + `
resource ` + BackupUnitResource + ` ` + BackupUnitTestResource + `similar {
	name        = "similar` + BackupUnitTestResource + `"
	password    = "DemoPassword1234$Updated"
	email       = "example-updated@ionoscloud.com"
}

data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "` + BackupUnitTestResource + `"
}
`

const testAccDataSourceBackupUnitMatchWrongNameError = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "wrong_name"
}
`

const testAccDataSourceBackupUnitPartialMatchWrongNameError = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "wrong_name"
  partial_match = true
}
`

const testAccDataSourceBackupUnitMatchPartialMatchDisabledError = testAccCheckBackupUnitConfigBasic + `
data ` + BackupUnitResource + ` ` + BackupUnitDataSourceByName + ` {
  name			= "` + DataSourcePartial + `"
}
`
