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

func TestAccBucketLifecycleConfigurationResourceBasic(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_lifecycle_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "rule.0.id", "Logs delete"),
					resource.TestCheckResourceAttr(name, "rule.0.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.0.prefix", "/logs"),
					resource.TestCheckResourceAttr(name, "rule.0.expiration.days", "90"),
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

func TestAccBucketLifecycleConfigurationResourceNoncurrent(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_lifecycle_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_noncurrent(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "rule.0.id", "Logs delete"),
					resource.TestCheckResourceAttr(name, "rule.0.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.0.prefix", "/logs"),
					resource.TestCheckResourceAttr(name, "rule.0.noncurrent_version_expiration.noncurrent_days", "90"),
				),
			},
		},
	})
}

func TestAccBucketLifecycleConfigurationResourceAbort(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_lifecycle_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_abort(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "rule.0.id", "Logs delete"),
					resource.TestCheckResourceAttr(name, "rule.0.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.0.prefix", "/logs"),
					resource.TestCheckResourceAttr(name, "rule.0.abort_incomplete_multipart_upload.days_after_initiation", "1"),
				),
			},
		},
	})
}

func TestAccBucketLifecycleConfigurationResourceMultipleRules(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_lifecycle_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketLifecycleConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketLifecycleConfigurationConfig_multipleRules(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLifecycleConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "rule.0.id", "Logs delete"),
					resource.TestCheckResourceAttr(name, "rule.0.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.0.prefix", "/logs"),
					resource.TestCheckResourceAttr(name, "rule.0.abort_incomplete_multipart_upload.days_after_initiation", "1"),
					resource.TestCheckResourceAttr(name, "rule.1.id", "Rule 2"),
					resource.TestCheckResourceAttr(name, "rule.1.status", "Enabled"),
					resource.TestCheckResourceAttr(name, "rule.1.prefix", "/logs"),
					resource.TestCheckResourceAttr(name, "rule.1.expiration.days", "90"),
				),
			},
		},
	})

}

func testAccBucketLifecycleConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketLifecycleConfigurationConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketLifecycleConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_lifecycle_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  rule {

	id = "Logs delete"
	status = "Enabled"
	
	prefix = "/logs"

	expiration {
      days = 90
    }
  }
}
`))
}

func testAccBucketLifecycleConfigurationConfig_noncurrent(bucketName string) string {
	return utils.ConfigCompose(testAccBucketLifecycleConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_lifecycle_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  rule {

	id = "Logs delete"
	status = "Enabled"
	
	prefix = "/logs"

	expiration {
      days = 90
    }

    noncurrent_version_expiration {
	  noncurrent_days = 90
	}
  }
}
`))
}

func testAccBucketLifecycleConfigurationConfig_abort(bucketName string) string {
	return utils.ConfigCompose(testAccBucketLifecycleConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_lifecycle_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  rule {

	id = "Logs delete"
	status = "Enabled"
	
	prefix = "/logs"
	
	abort_incomplete_multipart_upload {
	  days_after_initiation = 1
	}
  }
}
`))
}

func testAccBucketLifecycleConfigurationConfig_multipleRules(bucketName string) string {
	return utils.ConfigCompose(testAccBucketLifecycleConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_lifecycle_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  rule {

	id = "Logs delete"
	status = "Enabled"
	
	prefix = "/logs"
	
	abort_incomplete_multipart_upload {
	  days_after_initiation = 1
	}
  }

  rule {
    id = "Rule 2"
	status = "Enabled"
	prefix = "/logs"
	expiration {
      days = 90
    }
  }
}
`))
}

func testAccCheckBucketLifecycleConfigurationDestroy(s *terraform.State) error {
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_lifecycle_configuration" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.LifecycleApi.GetBucketLifecycle(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket lifecycle: %s", err)
			}

			return fmt.Errorf("bucket lifecycle still exists")
		}
	}

	return nil
}

func testAccCheckLifecycleConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client, err := acctest.S3Client()
		if err != nil {
			return err
		}

		_, _, err = client.LifecycleApi.GetBucketLifecycle(ctx, rs.Primary.Attributes["bucket"]).Execute()
		return err
	}
}
