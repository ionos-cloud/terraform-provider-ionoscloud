//go:build compute || all || s3key

package ionoscloud

import (
	"context"
	"fmt"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccS3KeyBasic(t *testing.T) {
	var s3Key ionoscloud.S3Key

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccChecks3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccChecks3KeyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccChecks3KeyExists(constant.S3KeyResource+"."+constant.S3KeyTestResource, &s3Key),
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
			// {
			//	Config: testAccChecks3KeyConfigUpdate,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccChecks3KeyExists(S3KeyResource+"."+S3KeyTestResource, &s3Key),
			//		resource.TestCheckResourceAttrSet(S3KeyResource+"."+S3KeyTestResource, "secret_key"),
			//		resource.TestCheckResourceAttr(S3KeyResource+"."+S3KeyTestResource, "active", "true"),
			//	),
			// },
		},
	})
}

func testAccChecks3KeyDestroyCheck(s *terraform.State) error {

	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.S3KeyResource {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while fetching s3 key %s: %w", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("s3 Key still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccChecks3KeyExists(n string, s3Key *ionoscloud.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occurred while fetching S3 Key: %s", rs.Primary.ID)
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
  first_name = "terraform"
  last_name = "test"
  email = "` + utils.GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource ` + constant.S3KeyResource + ` ` + constant.S3KeyTestResource + ` {
  user_id    = ` + constant.UserResource + `.example.id
  active     = true
}`

// this step is commented since the current behaviour of s3 keys is that when you create an s3 key with active set on false
// it is set to true by the API, so an update from false to true can not be done

// var testAccChecks3KeyConfigUpdate = `
//
//	resource ` + UserResource + ` "example" {
//	 first_name = "terraform"
//	 last_name = "test"
//	 email = "` + utils.GenerateEmail() + `"
//	 password = "abc123-321CBA"
//	 administrator = false
//	 force_sec_auth= false
//	}
//
//	resource ` + S3KeyResource + ` ` + S3KeyTestResource + ` {
//	 user_id    = ` + UserResource + `.example.id
//	 active     = true
//	}`
var testAccDataSourceS3KeyMatchId = testAccChecks3KeyConfigBasic + `
data ` + constant.S3KeyResource + ` ` + constant.S3KeyDataSourceById + ` {
user_id    	= ` + constant.UserResource + `.example.id
id			= ` + constant.S3KeyResource + `.` + constant.S3KeyTestResource + `.id
}
`
