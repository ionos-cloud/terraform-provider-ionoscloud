package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("data.ionoscloud_location.loc", "id", "de/fra"),
					resource.TestCheckResourceAttr("data.ionoscloud_location.loc", "name", "frankfurt"),
				),
			},
		},
	})

}

const testAccDataSourceLocation_basic = `
	data "ionoscloud_location" "loc" {
	  name = "frankfurt"
	  feature = "SSD"
	}
	`
