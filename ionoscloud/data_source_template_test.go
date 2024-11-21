//go:build compute || all || template

package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

const templateName = constant.DataSource + "." + constant.TemplateResource + "." + constant.TemplateTestResource

func TestAccDataSourceTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Basic Cube S"),
					resource.TestCheckResourceAttr(templateName, "cores", "2"),
					resource.TestCheckResourceAttr(templateName, "ram", "4096"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "120"),
				),
			},
			{
				Config: testAccDataSourceTemplateCoresRam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Basic Cube XL"),
				),
			},
			{
				Config: testAccDataSourceTemplateRamStorage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Memory Cube S"),
				),
			},
			{
				Config: testAccDataSourceTemplateStorageSizeRam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Basic Cube M"),
					resource.TestCheckResourceAttr(templateName, "cores", "4"),
					resource.TestCheckResourceAttr(templateName, "ram", "8192"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "240"),
				),
			},
			{
				Config:      testAccDataSourceTemplateStorageWrongNameError,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateStorageWrongCores,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateStorageWrongRam,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateStorageWrongStorage,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
		},
	})

}

const testAccDataSourceTemplateName = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name = "Basic Cube S"
}`

const testAccDataSourceTemplateCoresRam = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	cores = 16
	ram = 32768
}`

const testAccDataSourceTemplateRamStorage = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	ram = 8192
	storage_size = 120
}`

const testAccDataSourceTemplateStorageSizeRam = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	storage_size = 240
	ram = 8192
}`

const testAccDataSourceTemplateStorageWrongNameError = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name		 = "CUBES S"
	cores		 = 6
	ram			 = 16384
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongCores = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 1
	ram			 = 16384
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongRam = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 6
	ram			 = 2048
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongStorage = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 6
	ram			 = 16384
	storage_size = 50
}`
