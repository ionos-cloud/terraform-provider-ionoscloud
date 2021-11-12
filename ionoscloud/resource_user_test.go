package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccUserBasic(t *testing.T) {
	var user ionoscloud.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckUserDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserExists(UserResource+"."+UserTestResource, &user),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "first_name", UserTestResource),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "last_name", UserTestResource),
					resource.TestCheckResourceAttrSet(UserResource+"."+UserTestResource, "email"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "administrator", "true"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "force_sec_auth", "true"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "active", "true"),
				),
			},
			{
				Config: testAccCheckUserConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "first_name", UpdatedResources),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "last_name", UpdatedResources),
					resource.TestCheckResourceAttrSet(UserResource+"."+UserTestResource, "email"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "administrator", "false"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "force_sec_auth", "false"),
					resource.TestCheckResourceAttr(UserResource+"."+UserTestResource, "active", "false")),
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
		if rs.Type != UserResource {
			continue
		}
		_, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
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

		foundUser, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

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

var testAccCheckUserConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "` + UserTestResource + `"
  last_name = "` + UserTestResource + `"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = true
  force_sec_auth= true
  active  = true
}`

var testAccCheckUserConfigUpdate = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "` + UpdatedResources + `"
  last_name = "` + UpdatedResources + `"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBAupdated"
  administrator = false
  force_sec_auth= false
  active  = false
}`
