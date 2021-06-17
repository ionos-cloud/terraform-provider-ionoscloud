package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceTemplate_matching(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testaccdatasourcetemplateMatchingwithdatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "name", "BETA CUBES S"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "cores", "1"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "ram", "2048"),
				),
			},
		},
	})

}

const testaccdatasourcetemplateMatchingwithdatasource = `
data "ionoscloud_template" "template" {
	name = "BETA CUBES S"
	cores = 1
	ram			= 2048
	storage_size = 50
}`
