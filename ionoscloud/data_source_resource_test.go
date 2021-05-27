package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_resource.res", "resource_type", "datacenter"),
				),
			},
		},
	})

}

const testAccDataSourceResource_basic = `
resource "ionoscloud_datacenter" "foobar" {
  name       = "test_name"
  location = "us/las"
}

data "ionoscloud_resource" "res" {
  resource_type = "datacenter"
  resource_id="${ionoscloud_datacenter.foobar.id}"
}`
