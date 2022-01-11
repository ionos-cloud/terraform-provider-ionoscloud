package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceImageBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testaccdatasourceimageBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "cloud_init", "NONE"),
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "location", "de/fkb"),
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "name", "ubuntu-18.04.3-live-server-amd64.iso"),
					resource.TestCheckResourceAttr("data.ionoscloud_image.img", "type", "CDROM"),
				),
			},
		},
	})

}

const testaccdatasourceimageBasic = `
	data "ionoscloud_image" "img" {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`
