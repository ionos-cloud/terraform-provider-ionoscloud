//go:build compute || all || user

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccUserBasic(t *testing.T) {
	var user ionoscloud.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(constant.UserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "first_name", constant.UserTestResource),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "last_name", constant.UserTestResource),
					resource.TestCheckResourceAttrSet(constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "administrator", "true"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "force_sec_auth", "true"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "true"),
					resource.TestCheckResourceAttrPair(constant.UserResource+"."+constant.UserTestResource, "password", constant.RandomPassword+".user_password", "result"),
				),
			},
			{
				Config: testAccDataSourceUserMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "first_name", constant.UserResource+"."+constant.UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "last_name", constant.UserResource+"."+constant.UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "email", constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "administrator", constant.UserResource+"."+constant.UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "force_sec_auth", constant.UserResource+"."+constant.UserTestResource, "force_sec_auth"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "sec_auth_active", constant.UserResource+"."+constant.UserTestResource, "sec_auth_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "s3_canonical_user_id", constant.UserResource+"."+constant.UserTestResource, "s3_canonical_user_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "active", constant.UserResource+"."+constant.UserTestResource, "active"),
				),
			},
			{
				Config: testAccDataSourceUserMatchEmail,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "first_name", constant.UserResource+"."+constant.UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "last_name", constant.UserResource+"."+constant.UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "email", constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "administrator", constant.UserResource+"."+constant.UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "force_sec_auth", constant.UserResource+"."+constant.UserTestResource, "force_sec_auth"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "sec_auth_active", constant.UserResource+"."+constant.UserTestResource, "sec_auth_active"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "s3_canonical_user_id", constant.UserResource+"."+constant.UserTestResource, "s3_canonical_user_id"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceByName, "active", constant.UserResource+"."+constant.UserTestResource, "active"),
				),
			},
			{
				Config:      testAccDataSourceUserWrongEmail,
				ExpectError: regexp.MustCompile(`no user found with the specified criteria: email`),
			},
			{
				Config: testAccCheckUserConfigUpdateForceSec,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(constant.UserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "first_name", constant.UserTestResource),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "last_name", constant.UserTestResource),
					resource.TestCheckResourceAttrSet(constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "administrator", "true"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "force_sec_auth", "false"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "true"),
					resource.TestCheckResourceAttrPair(constant.UserResource+"."+constant.UserTestResource, "password", constant.RandomPassword+".user_password", "result"),
				),
			},
			{
				Config: testAccCheckUserConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "first_name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "last_name", constant.UpdatedResources),
					resource.TestCheckResourceAttrSet(constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "administrator", "false"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "force_sec_auth", "false"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "false"),
					resource.TestCheckResourceAttrPair(constant.UserResource+"."+constant.UserTestResource, "password", constant.RandomPassword+".user_password_updated", "result"),
				),
			},
			{
				Config: testAccCheckUserMultipleGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "group_ids.#", "3"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.*", map[string]string{
						"name": "group1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.*", map[string]string{
						"name": "group2",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.*", map[string]string{
						"name": "group3",
					})),
			},
			{
				Config: testAccCheckUserRemoveAllGroups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "group_ids.#", "0")),
			},
			{
				Config: testAccCheckNewUserGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.NewUserResource, "group_ids.#", "1"),
					resource.TestCheckResourceAttr(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "groups.*", map[string]string{
						"name": "group1",
					})),
			},
			{
				Config:             testAccCheckNewUserGroup,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccRemoveUserFromGroup(constant.GroupResource+".group1", constant.UserResource+"."+constant.NewUserResource)),
			},
		},
	})
}

func testAccCheckUserDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.UserResource {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("user still exists %s - an error occurred while checking it %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("user still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckUserExists(n string, user *ionoscloud.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("testAccCheckUserExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

		if cancel != nil {
			defer cancel()
		}

		foundUser, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching User: %s %s", rs.Primary.ID, err)
		}
		if *foundUser.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		user = &foundUser

		return nil
	}
}

func testAccRemoveUserFromGroup(group, user string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(bundleclient.SdkBundle).CloudApiClient
		gr, ok := s.RootModule().Resources[group]
		if !ok {
			return fmt.Errorf("testAccRemoveUserFromGroup: group not found: %s", group)
		}
		if gr.Primary.ID == "" {
			return fmt.Errorf("missing group id")
		}

		u, ok := s.RootModule().Resources[user]
		if !ok {
			return fmt.Errorf("testAccRemoveUserFromGroup: user not found: %s", user)
		}
		if u.Primary.ID == "" {
			return fmt.Errorf("missing user id")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		defer cancel()

		apiResponse, err := client.UserManagementApi.UmGroupsUsersDelete(ctx, gr.Primary.ID, u.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		return err
	}
}

var testAccCheckUserConfigBasic = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "` + constant.UserTestResource + `"
  last_name = "` + constant.UserTestResource + `"
  email = "` + utils.GenerateEmail() + `"
  password =  ` + constant.RandomPassword + `.user_password.result
  administrator = true
  force_sec_auth= true
  active  = true
}
`

var testAccCheckUserConfigUpdateForceSec = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name = "` + constant.UserTestResource + `"
 last_name = "` + constant.UserTestResource + `"
 email = "` + utils.GenerateEmail() + `"
 password = ` + constant.RandomPassword + `.user_password.result
 administrator = true
 force_sec_auth= false
 active  = true
}`

var testAccCheckUserConfigUpdate = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name = "` + constant.UpdatedResources + `"
 last_name = "` + constant.UpdatedResources + `"
 email = "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 administrator = false
 force_sec_auth= false
 active  = false
}`

var testAccCheckUserMultipleGroups = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name 	= "` + constant.UpdatedResources + `"
 last_name 		= "` + constant.UpdatedResources + `"
 email 			= "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 #password       = ` + constant.RandomPassword + `.user_password.result Updated
 administrator  = false
 force_sec_auth = false
 active  		= false
 group_ids 		= [ ionoscloud_group.group1.id, ionoscloud_group.group2.id, ionoscloud_group.group3.id]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group2" {
  name = "group2"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group3" {
  name = "group3"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
	id = ionoscloud_user.` + constant.UserTestResource + `.id
}
`

var testAccCheckUserMultipleGroups1Element = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name 	= "` + constant.UpdatedResources + `"
 last_name 		= "` + constant.UpdatedResources + `"
 email 			= "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 administrator  = false
 force_sec_auth = false
 active  		= false
 group_ids 		= [ ionoscloud_group.group1.id]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group2" {
  name = "group2"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group3" {
  name = "group3"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
	id = ionoscloud_user.` + constant.UserTestResource + `.id
}
`

var testAccCheckUserRemoveAllGroups = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name 	= "` + constant.UpdatedResources + `"
 last_name 		= "` + constant.UpdatedResources + `"
 email 			= "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 administrator  = false
 force_sec_auth = false
 active  		= false
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group2" {
  name = "group2"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group3" {
  name = "group3"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
	id = ionoscloud_user.` + constant.UserTestResource + `.id
}
`

var testAccCheckUserWrongGroupId = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
 first_name 	= "` + constant.UpdatedResources + `"
 last_name 		= "` + constant.UpdatedResources + `"
 email 			= "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 administrator  = false
 force_sec_auth = false
 active  		= false
 group_ids = ["notAnId"]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group2" {
  name = "group2"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}
resource "ionoscloud_group" "group3" {
  name = "group3"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
	id = ionoscloud_user.` + constant.UserTestResource + `.id
}
`

// Test in which we create the user and in the same time we add the user to a specific group. The
// difference between this test and the other group-related tests is that, for this test, the user
// is new. For the other user-group-related tests, we are operating on a user that already exists.
var testAccCheckNewUserGroup = `
resource ` + constant.RandomPassword + ` "user_password_updated" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.NewUserResource + ` {
 first_name 	= "` + constant.NewUserName + `"
 last_name 		= "` + constant.NewUserName + `"
 email 			= "` + utils.GenerateEmail() + `"
 password 		= ` + constant.RandomPassword + `.user_password_updated.result
 administrator  = false
 force_sec_auth = false
 active  		= false
 group_ids 		= [ ionoscloud_group.group1.id]
}

resource "ionoscloud_group" "group1" {
  name = "group1"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  create_k8s_cluster = true
}

data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
	id = ionoscloud_user.` + constant.NewUserResource + `.id
}
`

var testAccDataSourceUserMatchId = testAccCheckUserConfigBasic + `
data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
  id			= ` + constant.UserResource + `.` + constant.UserTestResource + `.id
}
`

var testAccDataSourceUserMatchEmail = testAccCheckUserConfigBasic + `
data ` + constant.UserResource + ` ` + constant.UserDataSourceByName + ` {
  email			= ` + constant.UserResource + `.` + constant.UserTestResource + `.email
}
`

var testAccDataSourceUserWrongEmail = testAccCheckUserConfigBasic + `
data ` + constant.UserResource + ` ` + constant.UserDataSourceByName + ` {
  email			= "wrong_email"
}
`
