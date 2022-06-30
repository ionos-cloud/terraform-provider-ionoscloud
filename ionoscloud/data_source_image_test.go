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
					resource.TestCheckResourceAttr(imageTestName, "name", "ubuntu-18.04.3-live-server-amd64.iso"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
				),
			},
			{
				Config: testAccDataSourceImageBasicId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "id", "f6032393-499f-11ea-94e0-525400f64d8d"),
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "NONE"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/txl"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
				),
			},
			{
				Config: testAccDataSourceImageBasicImageAlias,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "id", "f6032393-499f-11ea-94e0-525400f64d8d"),
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "NONE"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/txl"),
					resource.TestCheckResourceAttr(imageTestName, "type", "CDROM"),
					resource.TestCheckResourceAttr(imageTestName, "image_alias", "centos:7_iso"),
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
				Config:      testAccDataSourceImageWrongVersion,
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

const testAccDataSourceImageBasicId = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  id = "f6032393-499f-11ea-94e0-525400f64d8d"
	}
`

const testAccDataSourceImageBasicImageAlias = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  image_alias = "centos:7_iso"
	  location = "de/txl"
	}
`

const testAccDataSourceImageWrongNameError = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "wrong_name"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongType = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu"
	  type = "wrong_type"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongVersion = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "wrong_version"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongLocation = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "wrong_location"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongCloudInit = `
	data ` + ImageResource + ` ` + ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "wrong_cloud_init"
	}
`
