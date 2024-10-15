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
	name := "data.ionoscloud_s3_accesskey.testres"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.TestAccProtoV6ProviderFactories,
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAccesskeyDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "description", "desc"),
					resource.TestCheckResourceAttrSet(name, "id"),
					resource.TestCheckResourceAttrSet(name, "accesskey"),
					resource.TestCheckResourceAttrSet(name, "canonical_user_id"),
					resource.TestCheckResourceAttrSet(name, "contract_user_id"),
				),
			},
		},
	})
}

func testAccAccesskeyDataSourceConfigBasic() string {
	return utils.ConfigCompose(testAccAccesskeyConfigDescription("desc"), `
data "ionoscloud_s3_accesskey" "testres" {
	id = ionoscloud_s3_accesskey.test.id
}
`)
}
