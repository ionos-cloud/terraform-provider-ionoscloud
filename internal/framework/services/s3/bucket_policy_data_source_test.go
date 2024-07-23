//go:build all || s3
// +build all s3

package s3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccBucketPolicyDataSource(t *testing.T) {
	rName := "tf-acctest-test-bucket-policy"
	name := "ionoscloud_bucket_policy.test"

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.TestAccProtoV5ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBucketPolicyDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "bucket", rName),
					testAccCheckBucketPolicyData,
				),
			},
		},
	})
}

func testAccBucketPolicyDataSourceConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketPolicyConfig_basic(bucketName), `
data "ionoscloud_bucket_policy" "test" {
 bucket = ionoscloud_bucket_policy.test.bucket
}
`)
}
