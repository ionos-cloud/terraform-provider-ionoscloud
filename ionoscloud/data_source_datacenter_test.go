package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDatacenter_matching(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testaccdatasourcedatacenterMatching,
			},
			{
				Config: testaccdatasourcedatacenterMatchingwithdatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_datacenter.foobar", "name", "test_name"),
					resource.TestCheckResourceAttr("data.ionoscloud_datacenter.foobar", "location", "us/las"),
				),
			},
		},
	})

}

const testaccdatasourcedatacenterMatching = `
resource "ionoscloud_datacenter" "foobar" {
    name       = "test_name"
    location = "us/las"
}
`

const testaccdatasourcedatacenterMatchingwithdatasource = `
resource "ionoscloud_datacenter" "foobar" {
    name       = "test_name"
    location = "us/las"
}

data "ionoscloud_datacenter" "foobar" {
    name = "${ionoscloud_datacenter.foobar.name}"
    location = "us/las"
}`
