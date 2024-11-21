//go:build all || objectstoragemanagement
// +build all objectstoragemanagement

package objectstoragemanagement_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
)

func TestAccS3RegionDataSource(t *testing.T) {
	name := "data.ionoscloud_object_storage_region.testreg"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRegionDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "id", "de"),
					resource.TestCheckResourceAttr(name, "endpoint", "s3.eu-central-1.ionoscloud.com"),
					resource.TestCheckResourceAttr(name, "website", "s3-website.de-central.profitbricks.com"),
					resource.TestCheckResourceAttr(name, "storage_classes.0", "standard"),
					resource.TestCheckResourceAttr(name, "capability.iam", "false"),
				),
			},
		},
	})
}

func testAccRegionDataSourceConfigBasic() string {
	return `
data "ionoscloud_object_storage_region" "testreg" {
	id = "de"
}
`
}
