package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceTemplate_matching(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceTemplate_matchingWithDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "name", "BETA CUBES S"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "cores", "1"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "ram", "2048"),
				),
			},
		},
	})

}

const testAccDataSourceTemplate_matchingWithDataSource = `
data "ionoscloud_template" "template" {
	name = "BETA CUBES S"
	cores = 1
	ram			= 2048
	storage_size = 50
}`
