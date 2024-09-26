//go:build all || s3management
// +build all s3management

package s3management_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/acctest"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
)

func TestAccS3AccesskeyDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyDataSourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_s3_accesskey.testres", "description", "desc"),
				),
			},
		},
	})
}

func testAccAccesskeyDataSourceConfig_basic() string {
	return utils.ConfigCompose(testAccAccesskeyConfig_description("desc"), `
data "ionoscloud_s3_accesskey" "testres" {
	id = ionoscloud_s3_accesskey.test.id
}
`)
}
