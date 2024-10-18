//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketObjectLockConfigurationResourceBasic(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_object_lock_configuration.test"

	resource.Test(t, resource.TestCase{
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

func TestAccBucketObjectLockConfigurationResourceYears(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_object_lock_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketObjectLockConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketObjectLockConfigurationConfig_years(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "object_lock_enabled", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.mode", "GOVERNANCE"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.years", "1"),
				),
			},
		},
	})
}

func TestAccBucketObjectLockConfigurationResourceConflict(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketObjectLockConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccBucketObjectLockConfigurationConfig_conflict(bucketName),
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
		},
	})
}

func TestAccBucketObjectLockConfigurationResourceUpdate(t *testing.T) {
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_object_lock_configuration.test"

	resource.Test(t, resource.TestCase{
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
				Config: testAccBucketObjectLockConfigurationConfig_update(bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "object_lock_enabled", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.mode", "COMPLIANCE"),
					resource.TestCheckResourceAttr(name, "rule.default_retention.days", "2"),
				),
			},
		},
	})
}

func testAccBucketObjectLockConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
  object_lock_enabled = true
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

func testAccBucketObjectLockConfigurationConfig_update(bucketName string) string {
	return utils.ConfigCompose(testAccBucketObjectLockConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_object_lock_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  object_lock_enabled = "Enabled"
  rule {
	default_retention {
	  mode = "COMPLIANCE"
	  days = 2
  	}
  }
}
`))
}

func testAccBucketObjectLockConfigurationConfig_conflict(bucketName string) string {
	return utils.ConfigCompose(testAccBucketObjectLockConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_object_lock_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  object_lock_enabled = "Enabled"
  rule {
	default_retention {
	  mode = "COMPLIANCE"
	  days = 2
      years = 1
  	}
  }
}
`))
}

func testAccBucketObjectLockConfigurationConfig_years(bucketName string) string {
	return utils.ConfigCompose(testAccBucketObjectLockConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_object_lock_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  object_lock_enabled = "Enabled"
  rule {
	default_retention {
	  mode = "GOVERNANCE"
	  years = 1
  	}
  }
}
`))
}

func testAccCheckBucketObjectLockConfigurationDestroy(s *terraform.State) error {
	client, err := acctest.ObjectStorageClient()
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
