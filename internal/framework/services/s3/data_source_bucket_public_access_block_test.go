//go:build all || s3
// +build all s3

package s3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccBucketPublicAccessBlockDataSource(t *testing.T) {
	rName := "tf-acctest-test-bucket-policy"
	name := "ionoscloud_s3_bucket_public_access_block.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPublicAccessBlockDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					resource.TestCheckResourceAttr(name, "ignore_public_acls", "false"),
					resource.TestCheckResourceAttr(name, "restrict_public_buckets", "false"),
					resource.TestCheckResourceAttr(name, "block_public_policy", "true"),
					resource.TestCheckResourceAttr(name, "block_public_acls", "true"),
				),
			},
		},
	})
}

func testAccBucketPublicAccessBlockDataSourceConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketPublicAccessBlockConfig_basic(bucketName), `
data "ionoscloud_s3_bucket_public_access_block" "test" {
 bucket = ionoscloud_s3_bucket_public_access_block.test.bucket
}
`)
}
