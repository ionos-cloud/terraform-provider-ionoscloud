package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const imageTestName = DataSource + "." + ImageResource + "." + ImageTestResource

func TestAccDataSourceImageBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{

				Config: testAccDataSourceImageBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "NONE"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/fkb"),
					resource.TestCheckResourceAttr(imageTestName, "name", "ubuntu-18.04.3-live-server-amd64.iso"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
				),
			},
		},
	})

}

const testAccDataSourceImageBasic = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`
