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

func TestAccBucketCORSConfigurationResourceBasic(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_cors_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketCORSConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketCORSConfigurationConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCORSConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.0", "POST"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.1", "PUT"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_origins.0", "https://s3-website-test.hashicorp.com"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.expose_headers.0", "ETag"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.max_age_seconds", "3000"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.id", "1234"),
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

func TestAccBucketCORSConfigurationResourceMultiple(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_cors_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketCORSConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketCORSConfigurationConfig_multiple(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCORSConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_headers.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.0", "POST"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_methods.1", "PUT"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.allowed_origins.0", "https://s3-website-test.hashicorp.com"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.expose_headers.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.expose_headers.0", "ETag"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.max_age_seconds", "3000"),
					resource.TestCheckResourceAttr(name, "cors_rule.0.id", "1234"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_headers.#", "0"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_methods.#", "2"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_methods.0", "DELETE"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_methods.1", "GET"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_origins.#", "1"),
					resource.TestCheckResourceAttr(name, "cors_rule.1.allowed_origins.0", "https://s3.ionoscloud.com"),
				),
			},
		},
	})
}

func testAccBucketCORSConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketCORSConfigurationConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketCORSConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_cors_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
    id = 1234
  }
}
`))
}

func testAccBucketCORSConfigurationConfig_multiple(bucketName string) string {
	return utils.ConfigCompose(testAccBucketCORSConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_cors_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
    id = 1234
  }

 cors_rule {
    allowed_methods = ["DELETE", "GET"]
    allowed_origins = ["https://s3.ionoscloud.com"]
 }
}
`))
}

func testAccCheckBucketCORSConfigurationDestroy(s *terraform.State) error {
	client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_cors_configuration" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.CORSApi.GetBucketCors(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket cors: %w", err)
			}

			return fmt.Errorf("bucket cors still exists")
		}
	}

	return nil
}

func testAccCheckCORSConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client := acctest.NewTestBundleClientFromEnv().S3Client.GetBaseClient()
		_, _, err := client.CORSApi.GetBucketCors(ctx, rs.Primary.Attributes["bucket"]).Execute()
		return err
	}
}
