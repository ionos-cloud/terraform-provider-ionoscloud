//go:build all || s3management
// +build all s3management

package s3management_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccS3RegionDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ionoscloud_s3_region.testreg", "id", "ionoscloud_s3_region.testreg", "de"),
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
