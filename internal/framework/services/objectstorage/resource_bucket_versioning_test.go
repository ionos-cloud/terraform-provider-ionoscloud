//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketVersioningResource(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_versioning.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketVersioningDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketVersioningConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "versioning_configuration.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "versioning_configuration.mfa_delete", "Disabled"),
				),
			},
			{
				ResourceName:                         name,
				ImportStateId:                        bucketName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "bucket",
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccBucketVersioningConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketVersioningConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketVersioningConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_versioning" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  versioning_configuration {
	status = "Enabled"
  }
}
`))
}

func testAccCheckBucketVersioningDestroy(s *terraform.State) error {
	client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_versioning" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.VersioningApi.GetBucketVersioning(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket versioning: %s", err)
			}

			return fmt.Errorf("bucket versioning still exists")
		}
	}

	return nil
}
