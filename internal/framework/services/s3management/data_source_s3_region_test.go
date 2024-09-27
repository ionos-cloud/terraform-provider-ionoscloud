//go:build all || s3management
// +build all s3management

package s3management_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccS3RegionDataSource(t *testing.T) {
	name := "data.ionoscloud_s3_region.testreg"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_s3_region.testreg", "id", "de"),
					resource.TestCheckResourceAttr("data.ionoscloud_s3_region.testreg", "endpoint", "s3.eu-central-1.ionoscloud.com"),
					resource.TestCheckResourceAttr("data.ionoscloud_s3_region.testreg", "website", "s3-website.de-central.profitbricks.com"),
					resource.TestCheckResourceAttr("data.ionoscloud_s3_region.testreg", "storage_classes.0", "standard"),
					resource.TestCheckResourceAttr("data.ionoscloud_s3_region.testreg", "capability.iam", "false"),
				),
			},
		},
	})
}

func testAccRegionDataSourceConfig_basic() string {
	return `
data "ionoscloud_s3_region" "testreg" {
	id = "de"
}
`
}
