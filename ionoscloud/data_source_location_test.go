package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceLocation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceLocation_basic,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("data.ionoscloud_location.loc", "id", "de/fkb"),
					resource.TestCheckResourceAttr("data.ionoscloud_location.loc", "name", "karlsruhe"),
				),
			},
		},
	})

}

const testAccDataSourceLocation_basic = `
	data "ionoscloud_location" "loc" {
	  name = "karlsruhe"
	  feature = "SSD"
	}
	`
