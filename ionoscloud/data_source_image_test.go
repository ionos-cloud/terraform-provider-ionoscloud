//go:build compute || all || image

package ionoscloud

import (
	"regexp"
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
					resource.TestCheckResourceAttr(imageTestName, "name", "ubuntu-22.04-live-server-amd64.iso"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
				),
			},
			{
				Config:      testAccDataSourceImageWrongNameError,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceImageWrongType,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceImageWrongLocation,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
			{
				Config:      testAccDataSourceImageWrongCloudInit,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
			{
				Config: testAccDataSourceImageBasicPartialName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "NONE"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/txl"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
				),
			},
			{
				Config:      testAccDataSourceImageBasicWrongPartialName,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
		},
	})

}

const testAccDataSourceImageBasic = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu-22.04-live-server-amd64.iso"
	  type = "CDROM"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongNameError = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "wrong_name"
	  type = "CDROM"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongType = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu-22.04-live-server-amd64.iso"
	  type = "wrong_type"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongVersion = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu-22.04-live-server-amd64.iso"
	  type = "CDROM"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongLocation = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu-22.04-live-server-amd64.iso"
	  type = "CDROM"
	  location = "wrong_location"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongCloudInit = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu-22.04-live-server-amd64.iso"
	  type = "CDROM"
	  location = "de/fkb"
	  cloud_init = "wrong_cloud_init"
	}
`
const testAccDataSourceImageBasicPartialName = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "CentOS-7-x86_64-NetInstall"
      location = "de/txl"
      partial_match = true
}
`
const testAccDataSourceImageBasicWrongPartialName = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "wrong_name"
      location = "de/txl"
      partial_match = true
}
`
