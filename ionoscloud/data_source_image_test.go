//go:build compute || all || image

package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const imageTestName = constant.DataSource + "." + constant.ImageResource + "." + constant.ImageTestResource

func TestAccDataSourceImageBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceImageAliasLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(imageTestName, "cloud_init", "V1"),
					resource.TestCheckResourceAttr(imageTestName, "location", "de/txl"),
					resource.TestCheckResourceAttr(imageTestName, "type", "HDD"),
					resource.TestCheckResourceAttrSet(imageTestName, "expose_serial"),
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
					resource.TestCheckResourceAttrSet(imageTestName, "expose_serial"),
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
  image_alias           = "ubuntu:latest"
  location              = "de/txl"
}`

const testDataSourceImageAliasMultipleError = `data ` + constant.ImageResource + ` ` + constant.ImageTestResource + ` {
  image_alias           = "ubuntu:22.04_iso"
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
