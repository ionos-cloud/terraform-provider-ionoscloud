//go:build all || objectstorage
// +build all objectstorage

package objectstorage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccBucketPublicAccessBlockResource(t *testing.T) {
	rName := "acctest-tf-bucket"
	name := "ionoscloud_s3_bucket_public_access_block.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		CheckDestroy: testAccCheckBucketPublicAccessBlockDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					resource.TestCheckResourceAttr(name, "ignore_public_acls", "false"),
					resource.TestCheckResourceAttr(name, "restrict_public_buckets", "false"),
					resource.TestCheckResourceAttr(name, "block_public_policy", "true"),
					resource.TestCheckResourceAttr(name, "block_public_acls", "true"),
				),
			},
			{
				Config: testAccBucketPublicAccessBlockConfig_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					resource.TestCheckResourceAttr(name, "ignore_public_acls", "true"),
					resource.TestCheckResourceAttr(name, "restrict_public_buckets", "true"),
					resource.TestCheckResourceAttr(name, "block_public_policy", "false"),
					resource.TestCheckResourceAttr(name, "block_public_acls", "false"),
				),
			},
			{
				ResourceName:                         name,
				ImportStateId:                        rName,
				ImportState:                          true,
				ImportStateVerifyIdentifierAttribute: "bucket",
				ImportStateVerify:                    true,
			},
		},
	})
}

func testAccCheckBucketPublicAccessBlockDestroy(s *terraform.State) error {
	client, err := acctest.ObjectStorageClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_s3_bucket_public_access_block" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, apiResponse, err := client.PublicAccessBlockApi.GetPublicAccessBlock(context.Background(), rs.Primary.Attributes["bucket"]).Execute()
			if apiResponse.HttpNotFound() {
				return nil
			}

			if err != nil {
				return fmt.Errorf("error checking for bucket public access block: %s", err)
			}

			return fmt.Errorf("bucket public access block still exists")
		}
	}

	return nil
}

func testAccBucketPublicAccessBlockConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_bucket_public_access_block" "test"{
    bucket = ionoscloud_s3_bucket.test.name
    ignore_public_acls = false
    restrict_public_buckets = false
    block_public_policy = true
    block_public_acls = true
}
`, bucketName)
}

func testAccBucketPublicAccessBlockConfig_update(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_s3_bucket" "test" {
  name = %[1]q
}

resource "ionoscloud_s3_bucket_public_access_block" "test"{
    bucket = ionoscloud_s3_bucket.test.name
    ignore_public_acls = true
    restrict_public_buckets = true
    block_public_policy = false
    block_public_acls = false
}
`, bucketName)
}
