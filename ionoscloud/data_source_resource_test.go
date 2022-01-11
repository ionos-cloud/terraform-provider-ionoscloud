package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testaccdatasourceresourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_resource.res", "resource_type", "datacenter"),
				),
			},
		},
	})

}

const testaccdatasourceresourceBasic = `
resource "ionoscloud_datacenter" "foobar" {
  name       = "test_name"
  location = "us/las"
}

data "ionoscloud_resource" "res" {
  resource_type = "datacenter"
  resource_id="${ionoscloud_datacenter.foobar.id}"
}`
