//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketWebsiteConfigurationResourceBasic(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_website_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebsiteConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "index_document.suffix", "index.html"),
					resource.TestCheckResourceAttr(name, "error_document.key", "error.html"),
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

func TestAccBucketWebsiteConfigurationResourceIndex(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_website_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_index(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebsiteConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "index_document.suffix", "index.html"),
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

func TestAccBucketWebsiteConfigurationResourceRedirectAll(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_website_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketWebsiteConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketWebsiteConfigurationConfig_redirectAll(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWebsiteConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "redirect_all_requests_to.host_name", "example.com"),
					resource.TestCheckResourceAttr(name, "redirect_all_requests_to.protocol", "https"),
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

func testAccBucketWebsiteConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketWebsiteConfigurationConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketWebsiteConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_website_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  index_document {
    suffix = "index.html"
  }

  error_document {
    key = "error.html"
  }
}
`))
}

func testAccBucketWebsiteConfigurationConfig_index(bucketName string) string {
	return utils.ConfigCompose(testAccBucketWebsiteConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_website_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  index_document {
    suffix = "index.html"
  }
}
`))
}

func testAccBucketWebsiteConfigurationConfig_redirectAll(bucketName string) string {
	return utils.ConfigCompose(testAccBucketWebsiteConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_website_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  redirect_all_requests_to {
	host_name = "example.com"
    protocol = "https"
  }
}
`))
}

func testAccCheckBucketWebsiteConfigurationDestroy(s *terraform.State) error {
	client, err := acctest.ObjectStorageClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_website_configuration" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.WebsiteApi.GetBucketWebsite(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket website: %s", err)
			}

			return fmt.Errorf("bucket website still exists")
		}
	}

	return nil
}

func testAccCheckWebsiteConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client, err := acctest.ObjectStorageClient()
		if err != nil {
			return err
		}

		_, _, err = client.WebsiteApi.GetBucketWebsite(ctx, rs.Primary.Attributes["bucket"]).Execute()
		return err
	}
}
