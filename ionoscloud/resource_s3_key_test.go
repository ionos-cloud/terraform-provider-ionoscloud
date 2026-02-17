//go:build compute || all || s3key

package ionoscloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccKeyBasic(t *testing.T) {
	var s3Key ionoscloud.S3Key

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccChecksKeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccChecks3KeyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(constant.S3KeyResource+"."+constant.S3KeyTestResource, &s3Key),
					resource.TestCheckResourceAttrSet(constant.S3KeyResource+"."+constant.S3KeyTestResource, "secret_key"),
					resource.TestCheckResourceAttr(constant.S3KeyResource+"."+constant.S3KeyTestResource, "active", "true"),
				),
			},
			{
				Config: testAccDataSourceS3KeyMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(constant.S3KeyResource+"."+constant.S3KeyTestResource, "id"),
					resource.TestCheckResourceAttrPair(constant.S3KeyResource+"."+constant.S3KeyTestResource, "id", constant.DataSource+"."+constant.S3KeyResource+"."+constant.S3KeyDataSourceById, "id"),
					resource.TestCheckResourceAttrPair(constant.S3KeyResource+"."+constant.S3KeyTestResource, "secret", constant.DataSource+"."+constant.S3KeyResource+"."+constant.S3KeyDataSourceById, "secret"),
					resource.TestCheckResourceAttrPair(constant.S3KeyResource+"."+constant.S3KeyTestResource, "active", constant.DataSource+"."+constant.S3KeyResource+"."+constant.S3KeyDataSourceById, "active"),
				),
			},
			{
				Config: testAccChecks3KeyConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeyExists(constant.S3KeyResource+"."+constant.S3KeyTestResource, &s3Key),
					resource.TestCheckResourceAttrSet(constant.S3KeyResource+"."+constant.S3KeyTestResource, "secret_key"),
					resource.TestCheckResourceAttr(constant.S3KeyResource+"."+constant.S3KeyTestResource, "active", "false"),
				),
			},
		},
	})
}

func testAccChecksKeyDestroyCheck(s *terraform.State) error {

	client := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.S3KeyResource {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while fetching Object Storage key %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("Object Storage Key still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckKeyExists(n string, s3Key *ionoscloud.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		client := testAccProvider.Meta().(bundleclient.SdkBundle).NewCloudAPIClient("")

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		userId := rs.Primary.Attributes["user_id"]
		foundS3Key, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching Object Storage Key: %s", rs.Primary.ID)
		}

		if *foundS3Key.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		s3Key = &foundS3Key

		return nil
	}
}

var testAccChecks3KeyConfigBasic = `
resource ` + constant.UserResource + ` "example" {
  first_name 	 = "terraform"
  last_name 	 = "test"
  email 		 = "` + utils.GenerateEmail() + `"
  password 		 = "abc123-321CBA"
  administrator  = false
  force_sec_auth = false
  active		 = true
  group_ids      = [ ionoscloud_group.s3group.id]
}

resource "ionoscloud_group" "s3group" {
  name          = "S3 test Group"
  s3_privilege  = true
}

resource ` + constant.S3KeyResource + ` ` + constant.S3KeyTestResource + ` {
  user_id    = ` + constant.UserResource + `.example.id
  active     = true
}`

var testAccChecks3KeyConfigUpdate = `
resource ` + constant.UserResource + ` "example" {
  first_name 	 = "terraform"
  last_name 	 = "test"
  email 		 = "` + utils.GenerateEmail() + `"
  password		 = "abc123-321CBA"
  administrator  = false
  force_sec_auth = false
  active         = true
  group_ids 	 = [ ionoscloud_group.s3group.id]
}

resource "ionoscloud_group" "s3group" {
  name         = "S3 test Group"
  s3_privilege = true
}

resource ` + constant.S3KeyResource + ` ` + constant.S3KeyTestResource + ` {
  user_id = ` + constant.UserResource + `.example.id
  active  = false
}`
var testAccDataSourceS3KeyMatchId = testAccChecks3KeyConfigBasic + `
data ` + constant.S3KeyResource + ` ` + constant.S3KeyDataSourceById + ` {
  user_id = ` + constant.UserResource + `.example.id
  id      = ` + constant.S3KeyResource + `.` + constant.S3KeyTestResource + `.id
}
`
