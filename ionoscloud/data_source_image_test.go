package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceImage_basic(t *testing.T) {
	r, _ := regexp.Compile("Ubuntu-16.04")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceImage_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "location", "us/las"),
					resource.TestMatchResourceAttr("data.ionoscloud_image.img", "name", r),
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "type", "HDD"),
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "cloud_init", "NONE"),
				),
			},
		},
	})

}

const testAccDataSourceImage_basic = `
	data "ionoscloud_image" "img" {
	  name = "Ubuntu"
	  type = "HDD"
	  version = "16"
	  location = "us/las"
	  cloud_init = "NONE"
	}
`
