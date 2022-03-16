//go:build compute || all || s3key

package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAccS3KeyBasic(t *testing.T) {
	var s3Key ionoscloud.S3Key

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccChecks3KeyDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccChecks3KeyConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccChecks3KeyExists(S3KeyResource+"."+S3KeyTestResource, &s3Key),
					resource.TestCheckResourceAttrSet(S3KeyResource+"."+S3KeyTestResource, "secret_key"),
					resource.TestCheckResourceAttr(S3KeyResource+"."+S3KeyTestResource, "active", "true"),
				),
			},
			{
				Config: testAccDataSourceS3KeyMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(S3KeyResource+"."+S3KeyTestResource, "id"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "id", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "id"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "secret", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "secret"),
					resource.TestCheckResourceAttrPair(S3KeyResource+"."+S3KeyTestResource, "active", DataSource+"."+S3KeyResource+"."+S3KeyDataSourceById, "active"),
				),
			},
			//{
			//	Config: testAccChecks3KeyConfigUpdate,
			//	Check: resource.ComposeTestCheckFunc(
			//		testAccChecks3KeyExists(S3KeyResource+"."+S3KeyTestResource, &s3Key),
			//		resource.TestCheckResourceAttrSet(S3KeyResource+"."+S3KeyTestResource, "secret_key"),
			//		resource.TestCheckResourceAttr(S3KeyResource+"."+S3KeyTestResource, "active", "true"),
			//	),
			//},
		},
	})
}

func testAccChecks3KeyDestroyCheck(s *terraform.State) error {

	client := testAccProvider.Meta().(SdkBundle).CloudApiClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != S3KeyResource {
			continue
		}

		userId := rs.Primary.Attributes["user_id"]
		_, apiResponse, err := client.UserS3KeysApi.UmUsersS3keysFindByKeyId(context.TODO(), userId, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if apiResponse == nil || apiResponse.Response != nil && apiResponse.StatusCode != 404 {
				return fmt.Errorf("an error occurred while fetching s3 key %s: %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("s3 Key still exists %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccChecks3KeyExists(n string, s3Key *ionoscloud.S3Key) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		client := testAccProvider.Meta().(SdkBundle).CloudApiClient

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
			return fmt.Errorf("error occured while fetching S3 Key: %s", rs.Primary.ID)
		}

		if *foundS3Key.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		s3Key = &foundS3Key

		return nil
	}
}

var testAccChecks3KeyConfigBasic = `
resource ` + UserResource + ` "example" {
  first_name = "terraform"
  last_name = "test"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
}

resource ` + S3KeyResource + ` ` + S3KeyTestResource + ` {
  user_id    = ` + UserResource + `.example.id
  active     = true
}`

// this step is commented since the current behaviour of s3 keys is that when you create an s3 key with active set on false
// it is set to true by the API, so an update from false to true can not be done

//var testAccChecks3KeyConfigUpdate = `
//resource ` + UserResource + ` "example" {
//  first_name = "terraform"
//  last_name = "test"
//  email = "` + GenerateEmail() + `"
//  password = "abc123-321CBA"
//  administrator = false
//  force_sec_auth= false
//}
//
//resource ` + S3KeyResource + ` ` + S3KeyTestResource + ` {
//  user_id    = ` + UserResource + `.example.id
//  active     = true
//}`
var testAccDataSourceS3KeyMatchId = testAccChecks3KeyConfigBasic + `
data ` + S3KeyResource + ` ` + S3KeyDataSourceById + ` {
user_id    	= ` + UserResource + `.example.id
id			= ` + S3KeyResource + `.` + S3KeyTestResource + `.id
}
`
