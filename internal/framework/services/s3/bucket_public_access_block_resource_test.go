//go:build all || s3
// +build all s3

package s3_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/framework/services/s3"
)

func TestAccBucketPublicAccessBlockResource(t *testing.T) {
	rName := "acctest-tf-bucket"
	name := "ionoscloud_bucket_access_block.test"

	resource.ParallelTest(t, resource.TestCase{
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
	client, err := acctest.S3Client()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ionoscloud_bucket_access_block" {
			continue
		}

		if rs.Primary.Attributes["bucket"] != "" {
			_, err = s3.GetBucketPublicAccessBlock(context.Background(), client, rs.Primary.Attributes["bucket"])
			if err != s3.ErrBucketPublicAccessBlockNotFound {
				return err
			}
		}
	}

	return nil
}

func testAccBucketPublicAccessBlockConfig_basic(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_bucket" "test" {
  bucket = %[1]q
}

resource "ionoscloud_bucket_access_block" "test"{
    bucket = ionoscloud_bucket.test.bucket
    ignore_public_acls = false
    restrict_public_buckets = false
    block_public_policy = true
    block_public_acls = true
}
`, bucketName)
}

func testAccBucketPublicAccessBlockConfig_update(bucketName string) string {
	return fmt.Sprintf(`
resource "ionoscloud_bucket" "test" {
  bucket = %[1]q
}

resource "ionoscloud_bucket_access_block" "test"{
    bucket = ionoscloud_bucket.test.bucket
    ignore_public_acls = true
    restrict_public_buckets = true
    block_public_policy = false
    block_public_acls = false
}
`, bucketName)
}
