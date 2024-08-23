//go:build all || s3
// +build all s3

package s3_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccBucketDataSource(t *testing.T) {
	rName := acctest.GenerateRandomResourceName(bucketPrefix)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBucketDataSourceConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_s3_bucket.test", "name", "ionoscloud_s3_bucket.test", "name"),
					resource.TestCheckResourceAttrPair("data.ionoscloud_s3_bucket.test", "region", "ionoscloud_s3_bucket.test", "region"),
				),
			},
		},
	})
}

func testAccBucketDataSourceConfig_basic(bucketName string) string {
	return utils.ConfigCompose(testAccBucketConfig_basic(bucketName), `
data "ionoscloud_s3_bucket" "test" {
	name = ionoscloud_s3_bucket.test.name
}
`)
}
