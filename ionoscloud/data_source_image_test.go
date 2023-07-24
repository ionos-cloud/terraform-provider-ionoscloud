//go:build compute || all || image

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const imageTestName = constant.DataSource + "." + constant.ImageResource + "." + constant.ImageTestResource

func TestAccDataSourceImageBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageAliasLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "V1"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/txl"),
					resource.TestCheckResourceAttr(imageTestName, "name", "CentOS-7-GenericCloud-2211"),
					resource.TestCheckResourceAttr(imageTestName, "type", "HDD"),
				),
			},
			{
				Config:      testDataSourceImageAliasMultipleError,
				ExpectError: regexp.MustCompile("more than one image found, enable debug to learn more"),
			},
			{
				Config:      testAccDataSourceWrongAliasError,
				ExpectError: regexp.MustCompile("no image found with the specified criteria"),
			},
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

const testDataSourceImageAliasLocation = `data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
  image_alias           = "centos:latest"
  location              = "de/txl"
}`

const testDataSourceImageAliasMultipleError = `data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
  image_alias           = "centos:latest"
}`

const testAccDataSourceWrongAliasError = `data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
  image_alias           = "doesNotExist"
  location              = "de/txl"
}`

const testAccDataSourceImageBasic = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "22.04-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongNameError = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "wrong_name"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongType = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "ubuntu"
	  type = "wrong_type"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongVersion = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "wrong_version"
	  location = "de/fkb"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongLocation = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "wrong_location"
	  cloud_init = "NONE"
	}
`

const testAccDataSourceImageWrongCloudInit = `
	data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
	  name = "ubuntu"
	  type = "CDROM"
	  version = "18.04.3-live-server-amd64.iso"
	  location = "de/fkb"
	  cloud_init = "wrong_cloud_init"
	}
`
