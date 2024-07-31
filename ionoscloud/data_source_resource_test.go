//go:build compute || all || resource

package ionoscloud

import (
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const resourceName = constant.DataSource + "." + constant.ResourceResource + "." + constant.ResourceTestResource

func TestAccResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "resource_type", "datacenter"),
				),
			},
		},
	})

}

const testAccDataSourceResourceBasic = `
resource ` + constant.DatacenterResource + ` "foobar" {
  name       = "test_name"
  location = "us/las"
}

data ` + constant.ResourceResource + ` ` + constant.ResourceTestResource + ` {
  resource_type = "datacenter"
  resource_id= ` + constant.DatacenterResource + `.foobar.id
}`
