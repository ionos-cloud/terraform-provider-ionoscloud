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

func TestAccBucketSSEConfigurationResourceBasic(t *testing.T) {
	ctx := context.Background()
	bucketName := acctest.GenerateRandomResourceName(bucketPrefix)
	name := "ionoscloud_s3_bucket_server_side_encryption_configuration.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketSSEConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketSSEConfigurationConfig_basic(bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSSEConfigurationExists(ctx, name),
					resource.TestCheckResourceAttr(name, "bucket", bucketName),
					resource.TestCheckResourceAttr(name, "rule.0.apply_server_side_encryption_by_default.sse_algorithm", "AES256"),
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

func testAccBucketSSEConfigurationConfig_base(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
  region = "eu-central-3"
}
`, bucketName)
}

func testAccBucketSSEConfigurationConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketSSEConfigurationConfig_base(bucketName), fmt.Sprintf(`
resource "ionoscloud_s3_bucket_server_side_encryption_configuration" "test" {
  bucket = ionoscloud_s3_bucket.test.name
  rule {
	apply_server_side_encryption_by_default {
	  sse_algorithm = "AES256"
	}
  }
}
`))
}

func testAccCheckBucketSSEConfigurationDestroy(s *terraform.State) error {
	client, err := acctest.ObjectStorageClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_server_side_encryption_configuration" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.EncryptionApi.GetBucketEncryption(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket server_side_encryption: %s", err)
			}

			return fmt.Errorf("bucket server_side_encryption still exists")
		}
	}

	return nil
}

func testAccCheckSSEConfigurationExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		client, err := acctest.ObjectStorageClient()
		if err != nil {
			return err
		}

		_, _, err = client.EncryptionApi.GetBucketEncryption(ctx, rs.Primary.Attributes["bucket"]).Execute()
		return err
	}
}
