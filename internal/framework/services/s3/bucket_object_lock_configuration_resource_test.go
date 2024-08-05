//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketObjectLockConfigurationResource(t *testing.T) {
	bucketName := "acctest-tf-bucket-object-lock-configuration"
	name := "ionoscloud_s3_bucket_object_lock_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketObjectLockConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectLockConfigurationConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "object_lock_enabled", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.days", "1"),
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

func testAccBucketObjectLockConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketObjectLockConfigurationConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketObjectLockConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_object_lock_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  object_lock_enabled = "Enabled"
  rule {
	default_retention {
	  mode = "GOVERNANCE"
	  days = 1
  	}
  }
}
`))
}

func testAccCheckBucketObjectLockConfigurationDestroy(s *terraform.State) error {
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_object_lock_configuration" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.ObjectLockApi.GetObjectLockConfiguration(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket object_lock_configuration: %s", err)
			}

			return fmt.Errorf("bucket object_lock_configuration still exists")
		}
	}

	return nil
}
