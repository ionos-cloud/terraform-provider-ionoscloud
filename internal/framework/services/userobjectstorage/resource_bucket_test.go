//go:build all || userobjectstorage

package userobjectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

const bucketPrefix = "acctest-tf-user-bucket-"

func TestAccUserObjectStorageBucketResource(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_user_object_storage_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckUserBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBucketConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUserBucketExists(context.Background(), resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "region", "de"),
					resource.TestCheckResourceAttr(resourceName, "id", rName),
					resource.TestCheckResourceAttr(resourceName, "object_lock_enabled", "false"),
				),
			},
			{
				ResourceName:                         resourceName,
				ImportStateId:                        rName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "name",
				ImportStateVerifyIgnore:              []string{"force_destroy", "object_lock_enabled"},
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccCheckUserBucketExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		client := acctest.NewTestBundleClientFromEnv().UserS3Client.GetBaseClient()
		_, err := client.BucketsApi.HeadBucket(ctx, rs.Primary.Attributes["name"]).Execute()
		return err
	}
}

func testAccCheckUserBucketDestroy(s *terraform.State) error {
	client := acctest.NewTestBundleClientFromEnv().UserS3Client.GetBaseClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_user_object_storage_bucket" {
			continue
		}
		if rs.Primary.Attributes["name"] == "" {
			continue
		}
		apiResponse, err := client.BucketsApi.HeadBucket(
			context.Background(), rs.Primary.Attributes["name"],
		).Execute()
		if apiResponse.HttpNotFound() {
			return nil
		}
		if err != nil {
			return err
		}
		return fmt.Errorf("user object storage bucket still exists: %s", rs.Primary.Attributes["name"])
	}

	return nil
}

func TestAccUserObjectStorageBucketDataSource(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	resourceName := "ionoscloud_user_object_storage_bucket.test"
	dataSourceName := "data.ionoscloud_user_object_storage_bucket.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckUserBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserBucketDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "region", resourceName, "region"),
				),
			},
		},
	})
}

func testAccUserBucketConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_user_object_storage_bucket" "test" {
  name = %[1]q
}
`, bucketName)
}

func testAccUserBucketDataSourceConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccUserBucketConfig_basic(bucketName), `
data "ionoscloud_user_object_storage_bucket" "test" {
  name   = ionoscloud_user_object_storage_bucket.test.name
  region = ionoscloud_user_object_storage_bucket.test.region
}
`)
}
