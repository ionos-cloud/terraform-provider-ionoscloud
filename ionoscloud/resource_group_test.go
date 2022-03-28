//go:build compute || all || group

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGroupBasic(t *testing.T) {
	var group ionoscloud.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckGroupDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGroupConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(GroupResource+"."+GroupTestResource, &group),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "name", GroupTestResource),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_datacenter", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_snapshot", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "reserve_ip", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "access_activity_log", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_pcc", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "s3_privilege", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_backup_unit", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_internet_access", "true"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_k8s_cluster", "true"),
					testNotEmptySlice(GroupResource, "users")),
			},
			{
				Config: testAccDataSourceGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "name", GroupResource+"."+GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_datacenter", GroupResource+"."+GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_snapshot", GroupResource+"."+GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "reserve_ip", GroupResource+"."+GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "access_activity_log", GroupResource+"."+GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_pcc", GroupResource+"."+GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "s3_privilege", GroupResource+"."+GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_backup_unit", GroupResource+"."+GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_internet_access", GroupResource+"."+GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_k8s_cluster", GroupResource+"."+GroupTestResource, "create_k8s_cluster"),
					testNotEmptySlice(DataSource+"."+GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "name", GroupResource+"."+GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_datacenter", GroupResource+"."+GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_snapshot", GroupResource+"."+GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "reserve_ip", GroupResource+"."+GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "access_activity_log", GroupResource+"."+GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_pcc", GroupResource+"."+GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "s3_privilege", GroupResource+"."+GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_backup_unit", GroupResource+"."+GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_internet_access", GroupResource+"."+GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_k8s_cluster", GroupResource+"."+GroupTestResource, "create_k8s_cluster"),
					testNotEmptySlice(DataSource+"."+GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "name", GroupResource+"."+GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_datacenter", GroupResource+"."+GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_snapshot", GroupResource+"."+GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "reserve_ip", GroupResource+"."+GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "access_activity_log", GroupResource+"."+GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_pcc", GroupResource+"."+GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "s3_privilege", GroupResource+"."+GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_backup_unit", GroupResource+"."+GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_internet_access", GroupResource+"."+GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_k8s_cluster", GroupResource+"."+GroupTestResource, "create_k8s_cluster"),
					testNotEmptySlice(DataSource+"."+GroupResource, "users"),
				),
			},
			{
				Config:      testAccDataSourceGroupMultipleResultsError,
				ExpectError: regexp.MustCompile("more than one group found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceGroupWrongNameError,
				ExpectError: regexp.MustCompile("no group found with the specified name"),
			},
			{
				Config: testAccCheckGroupConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(GroupResource+"."+GroupTestResource, &group),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_datacenter", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_snapshot", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "reserve_ip", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "access_activity_log", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_pcc", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "s3_privilege", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_backup_unit", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_internet_access", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_k8s_cluster", "false"),
					resource.TestCheckResourceAttrPair(GroupResource+".test_user_id", "users.0.id", UserResource+"."+UserTestResource+"3", "id"),
					testNotEmptySlice(GroupResource, "users")),
			},
			{
				Config: testAccCheckGroupUpdateMigrateToUserIds,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGroupExists(GroupResource+"."+GroupTestResource, &group),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "name", UpdatedResources),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_datacenter", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_snapshot", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "reserve_ip", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "access_activity_log", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_pcc", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "s3_privilege", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_backup_unit", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_internet_access", "false"),
					resource.TestCheckResourceAttr(GroupResource+"."+GroupTestResource, "create_k8s_cluster", "false"),
					resource.TestCheckResourceAttrPair(GroupResource+".test_user_id", "users.0.id", UserResource+"."+UserTestResource+"3", "id"),
					testNotEmptySlice(GroupResource, "users")),
			},
			{
				Config:      testAccCheckGroupBothUserArgumentsError,
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

func testAccCheckGroupDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {

		if rs.Type != GroupResource {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while checking the destruction of group %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("group %s still exists", rs.Primary.ID)
		}

	}

	return nil
}

func testAccCheckGroupExists(n string, group *ionoscloud.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckGroupExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundgroup, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occured while fetching Group: %s", rs.Primary.ID)
		}
		if *foundgroup.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		group = &foundgroup

		return nil
	}
}

var testAccCheckGroupCreateUsers = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + UserResource + ` ` + UserTestResource + `2 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

`

var testAccCheckGroupConfigBasic = testAccCheckGroupCreateUsers + `
resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  user_ids = [` + UserResource + `.` + UserTestResource + `.id, ` + UserResource + `.` + UserTestResource + `2.id]
}
`

var testAccDataSourceGroupMatchId = testAccCheckGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceById + ` {
  id			= ` + GroupResource + `.` + GroupTestResource + `.id
}
`

var testAccDataSourceGroupMatchName = testAccCheckGroupConfigBasic + `
resource ` + GroupResource + ` ` + GroupTestResource + `similar {
  name = "similar` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}
data ` + GroupResource + ` ` + GroupDataSourceByName + ` {
  name			= "` + GroupTestResource + `"
}
`

var testAccDataSourceGroupMultipleResultsError = testAccCheckGroupConfigBasic + `
resource ` + GroupResource + ` ` + GroupTestResource + `_multiple_results {
  name = "` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
}

data ` + GroupResource + ` ` + GroupDataSourceByName + ` {
  name			= "` + GroupTestResource + `"
}
`

var testAccDataSourceGroupWrongNameError = testAccCheckGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceByName + ` {
  name			= "wrong_name"
}
`

var testAccCheckGroupConfigUpdate = testAccCheckGroupCreateUsers + `
resource ` + UserResource + ` ` + UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` "test_user_id" {
  name = "` + GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_id = ` + UserResource + `.` + UserTestResource + `3.id
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + UserResource + `.` + UserTestResource + `.id, ` + UserResource + `.` + UserTestResource + `3.id]
}
`

var testAccCheckGroupUpdateMigrateToUserIds = testAccCheckGroupCreateUsers + `
resource ` + UserResource + ` ` + UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` "test_user_id" {
  name = "` + GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + UserResource + `.` + UserTestResource + `3.id]
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + UserResource + `.` + UserTestResource + `.id, ` + UserResource + `.` + UserTestResource + `3.id]
}
`

var testAccCheckGroupBothUserArgumentsError = testAccCheckGroupCreateUsers + `
resource ` + UserResource + ` ` + UserTestResource + `3 {
  first_name = "user"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` "test_user_id" {
  name = "` + GroupTestResource + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + UserResource + `.` + UserTestResource + `3.id]
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + UpdatedResources + `"
  create_datacenter = false
  create_snapshot = false
  reserve_ip = false
  access_activity_log = false
  create_pcc = false
  s3_privilege = false
  create_backup_unit = false
  create_internet_access = false
  create_k8s_cluster = false
  user_ids = [` + UserResource + `.` + UserTestResource + `.id, ` + UserResource + `.` + UserTestResource + `3.id]
  user_id = ` + UserResource + `.` + UserTestResource + `.id
}
`
