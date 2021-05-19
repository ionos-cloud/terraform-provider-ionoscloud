package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccS3Key_Basic(t *testing.T) {
	var s3Key ionoscloud.S3Key
	s3KeyName := "example"
	email := fmt.Sprintf("terraform_test-%d@mailinator.com",time.Now().Unix())

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccChecks3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccChecks3KeyConfigBasic, email, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccChecks3KeyExists("ionoscloud_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttr("ionoscloud_s3_key.example", "active", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccChecks3KeyConfigUpdate, email, s3KeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccChecks3KeyExists("ionoscloud_s3_key.example", &s3Key),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "secret_key"),
					resource.TestCheckResourceAttrSet("ionoscloud_s3_key.example", "active"),
				),
			},
		},
	})
}

func testAccChecks3KeyDestroyCheck(s *terraform.State) error {
	//client := testAccProvider.Meta().(*ionoscloud.APIClient)
	client := testAccProvider.Meta().(SdkBundle).Client

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_key" {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		_, apiResponse, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()

		if apiError, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse.Response.StatusCode != 404 {
				return fmt.Errorf("S3 Key still exists %s %s", rs.Primary.ID, apiError)
			}
		} else {
			return fmt.Errorf("Unable to fetch S3 Key %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccChecks3KeyExists(n string, s3Key *ionoscloud.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		//client := testAccProvider.Meta().(*ionoscloud.APIClient)
		client := testAccProvider.Meta().(SdkBundle).Client

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		userId := rs.Primary.Attributes["user_id"]
		foundS3Key, _, err := client.UserManagementApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()

		if err != nil {
			return fmt.Errorf("Error occured while fetching S3 Key: %s", rs.Primary.ID)
		}

		if *foundS3Key.Id != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		s3Key = &foundS3Key

		return nil
	}
}

const testAccChecks3KeyConfigBasic = `
resource "ionoscloud_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "%s" {
  user_id    = ionoscloud_user.example.id
  active     = true
}`

const testAccChecks3KeyConfigUpdate = `
resource "ionoscloud_user" "example" {
  first_name = "terraform"
  last_name = "test"
  email = "%s"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource "ionoscloud_s3_key" "%s" {
  user_id    = ionoscloud_user.example.id
  active     = false
}`
