//go:build compute || all || template

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
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "name", "CUBES L"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "cores", "4"),
					resource.TestCheckResourceAttr("data.ionoscloud_template.template", "ram", "8192"),
				),
			},
		},
	})

}

const testaccdatasourcetemplateMatchingwithdatasource = `
data "ionoscloud_template" "template" {
	cores = 4
}`
