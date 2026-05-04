//go:build compute || all || user

package ionoscloud

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

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
				Config: testAccCheckNewUserGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccRemoveUserFromGroup(constant.GroupResource+".group1", constant.UserResource+"."+constant.NewUserResource)),
			},
		},
	})
}

func TestUserWriteOnlyPassword(t *testing.T) {
	var user ionoscloud.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders: randomProviderVersion343(),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.11.0"))),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserConfigBasicWriteOnly,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(constant.UserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "password_wo_version", "1"),
					resource.TestCheckNoResourceAttr(constant.UserResource+"."+constant.UserTestResource, "password"),
					resource.TestCheckNoResourceAttr(constant.UserResource+"."+constant.UserTestResource, "password_wo"),
				),
			},
			{
				Config: testAccCheckUserConfigBasicWriteOnlyUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(constant.UserResource+"."+constant.UserTestResource, &user),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "password_wo_version", "2"),
					resource.TestCheckNoResourceAttr(constant.UserResource+"."+constant.UserTestResource, "password"),
				),
			},
		},
	})
}

func testAccCheckUserDestroyCheck(s *terraform.State) error {
	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
	if err != nil {
		return err
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

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
		if err != nil {
			return err
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

		client, err := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClientWithFailover(ctx)
		if err != nil {
			return err
		}

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

var testAccCheckUserConfigBasicWriteOnly = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "` + constant.UserTestResource + `"
  last_name = "` + constant.UserTestResource + `"
  email = "` + utils.GenerateEmail() + `"
  password_wo =  ` + constant.RandomPassword + `.user_password.result
  password_wo_version = 1
  administrator = true
  force_sec_auth= true
  active  = true
}
`

var testAccCheckUserConfigBasicWriteOnlyUpdated = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name = "` + constant.UserTestResource + `"
  last_name = "` + constant.UserTestResource + `"
  email = "` + utils.GenerateEmail() + `"
  password_wo =  ` + constant.RandomPassword + `.user_password.result
  password_wo_version = 2
  administrator = true
  force_sec_auth= true
  active  = true
}
`

// TestAccUserFailoverCreatesOnSecondEndpoint verifies that when the first configured endpoint
// returns a retryable HTTP 503 error, the failover round-tripper retries the request against
// the second endpoint (the real IONOS API) and the user is successfully created there.
//
// The test starts a local HTTP server that always returns 503, writes a temporary file config
// that lists the mock server first and the real IONOS API second, and configures the provider
// via IONOS_CONFIG_FILE to use that config with the roundRobin failover strategy.
func TestAccUserFailoverCreatesOnSecondEndpoint(t *testing.T) {
	if os.Getenv("IONOS_API_URL") != "" {
		t.Skip("IONOS_API_URL is set; this test requires control over endpoint configuration")
	}

	// Count how many times the mock server is called so we can assert failover occurred.
	var mockCallCount int32

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&mockCallCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer mockServer.Close()

	// File config: mock server (returns 503) first, real IONOS API second.
	//
	// Both endpoints must include the full API base path (/cloudapi/v6) because:
	//   - The ionoscloud SDK builds request URLs via simple concatenation:
	//     localVarPath = localBasePath + "/um/users"
	//   - The failover round-tripper rewrites only scheme+host, preserving the path.
	// So the SDK builds "http://mock:PORT/cloudapi/v6/um/users"; failover rewrites
	// scheme+host to produce "https://api.ionos.com/cloudapi/v6/um/users". ✓
	//
	// POST is added to retryableMethods so that user-create requests also participate in failover.
	configContent := fmt.Sprintf(`version: 1.0
environments:
  - name: default
    products:
      - name: cloud
        endpoints:
          - name: %s/cloudapi/v6
          - name: https://api.ionos.com/cloudapi/v6
failover:
  strategy: roundRobin
  failoverOnStatusCodes:
    - 503
  retryableMethods:
    - GET
    - POST
    - PUT
    - DELETE
    - HEAD
    - OPTIONS
  maxRetries: 1
`, mockServer.URL)

	tmpFile := createTempConfigFile(t, configContent)
	defer os.Remove(tmpFile)

	t.Setenv("IONOS_CONFIG_FILE", tmpFile)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserFailoverConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "first_name", constant.UserTestResource),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "last_name", constant.UserTestResource),
					resource.TestCheckResourceAttrSet(constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "administrator", "false"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "true"),
					// Verify the mock server was hit, proving requests went through the failover path.
					func(_ *terraform.State) error {
						if atomic.LoadInt32(&mockCallCount) == 0 {
							return fmt.Errorf("mock server was never called; failover round-tripper did not route through the first endpoint")
						}
						return nil
					},
				),
			},
		},
	})
}

// TestAccUserFailoverNetworkError verifies failover when a network-level error occurs (e.g., connection refused).
func TestAccUserFailoverNetworkError(t *testing.T) {
	if os.Getenv("IONOS_API_URL") != "" {
		t.Skip("IONOS_API_URL is set; this test requires control over endpoint configuration")
	}

	// Use an address that is likely to refuse connection
	badAddr := "http://127.0.0.1:1"

	configContent := fmt.Sprintf(`version: 1.0
environments:
  - name: default
    products:
      - name: cloud
        endpoints:
          - name: %s/cloudapi/v6
          - name: https://api.ionos.com/cloudapi/v6
failover:
  strategy: roundRobin
  retryableMethods:
    - GET
    - POST
    - DELETE
    - PUT
  maxRetries: 1
`, badAddr)

	tmpFile := createTempConfigFile(t, configContent)
	defer os.Remove(tmpFile)
	t.Setenv("IONOS_CONFIG_FILE", tmpFile)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ExternalProviders:        randomProviderVersion343(),
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserFailoverConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "true"),
				),
			},
			{
				Config: testAccCheckUserFailoverDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "first_name", constant.UserResource+"."+constant.UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "last_name", constant.UserResource+"."+constant.UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "email", constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "administrator", constant.UserResource+"."+constant.UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "active", constant.UserResource+"."+constant.UserTestResource, "active"),
				),
			},
			{
				Config: testAccCheckUserFailoverDataSourceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "first_name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "last_name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "administrator", "true"),
					resource.TestCheckResourceAttr(constant.UserResource+"."+constant.UserTestResource, "active", "false"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "first_name", constant.UserResource+"."+constant.UserTestResource, "first_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "last_name", constant.UserResource+"."+constant.UserTestResource, "last_name"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "email", constant.UserResource+"."+constant.UserTestResource, "email"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "administrator", constant.UserResource+"."+constant.UserTestResource, "administrator"),
					resource.TestCheckResourceAttrPair(constant.DataSource+"."+constant.UserResource+"."+constant.UserDataSourceById, "active", constant.UserResource+"."+constant.UserTestResource, "active"),
				),
			},
		},
	})
}

// Helper functions and shared logic for failover tests

func createTempConfigFile(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "ionos-failover-test-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp config file: %s", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write config file: %s", err)
	}
	tmpFile.Close()
	return tmpFile.Name()
}

func verifyMockCallCount(t *testing.T, count *int32) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if atomic.LoadInt32(count) == 0 {
			return fmt.Errorf("mock server was never called")
		}
		return nil
	}
}

var testAccCheckUserFailoverConfig = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name     = "` + constant.UserTestResource + `"
  last_name      = "` + constant.UserTestResource + `"
  email          = "` + utils.GenerateEmail() + `"
  password       = ` + constant.RandomPassword + `.user_password.result
  administrator  = false
  force_sec_auth = false
  active         = true
}
`

var testAccCheckUserFailoverDataSourceConfig = testAccCheckUserFailoverConfig + `
data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
  id = ` + constant.UserResource + `.` + constant.UserTestResource + `.id
}
`

var testAccCheckUserFailoverConfigUpdate = `
resource ` + constant.RandomPassword + ` "user_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource ` + constant.UserResource + ` ` + constant.UserTestResource + ` {
  first_name     = "` + constant.UpdatedResources + `"
  last_name      = "` + constant.UpdatedResources + `"
  email          = "` + utils.GenerateEmail() + `"
  password       = ` + constant.RandomPassword + `.user_password.result
  administrator  = true
  force_sec_auth = false
  active         = false
}
`

var testAccCheckUserFailoverDataSourceConfigUpdate = testAccCheckUserFailoverConfigUpdate + `
data ` + constant.UserResource + ` ` + constant.UserDataSourceById + ` {
  id = ` + constant.UserResource + `.` + constant.UserTestResource + `.id
}
`
