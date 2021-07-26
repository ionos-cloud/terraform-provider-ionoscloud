package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"

	"math/rand"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUser_Basic(t *testing.T) {
	var user ionoscloud.User
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	email := strconv.Itoa(r1.Intn(100000)) + "terraform_test" + strconv.Itoa(r1.Intn(100000)) + "@go.com"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testacccheckuserconfigBasic, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists("ionoscloud_user.resource_user", &user),
					testAccCheckUserAttributes("ionoscloud_user.resource_user", "resource_user"),
					resource.TestCheckResourceAttr("ionoscloud_user.resource_user", "first_name", "resource_user"),
					resource.TestCheckResourceAttr("ionoscloud_user.resource_user", "active", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testacccheckuserconfigUpdate, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserAttributes("ionoscloud_user.resource_user", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_user.resource_user", "first_name", "updated"),
					resource.TestCheckResourceAttr("ionoscloud_user.resource_user", "active", "false"),
				),
			},
		},
	})
}

func testAccCheckUserDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_user" {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			if apiResponse == nil || apiResponse.StatusCode != 404 {
				return fmt.Errorf("user still exists %s - an error occurred while checking it %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("user still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckUserAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckUserAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["first_name"] != name {
			return fmt.Errorf("bad first_name: %s", rs.Primary.Attributes["first_name"])
		}

		return nil
	}
}

func testAccCheckUserExists(n string, user *ionoscloud.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(SdkBundle).CloudApiClient
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

		foundUser, _, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("error occured while fetching User: %s %s", rs.Primary.ID, err)
		}
		if *foundUser.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		user = &foundUser

		return nil
	}
}

const testacccheckuserconfigBasic = `
resource "ionoscloud_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
}

resource "ionoscloud_user" "resource_user" {
  first_name = "resource_user"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active  = true
}`

const testacccheckuserconfigUpdate = `
resource "ionoscloud_user" "resource_user" {
  first_name = "updated"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource "ionoscloud_group" "group" {
  name = "terraform user group"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = false
  user_id="${ionoscloud_user.resource_user.id}"
}
`
