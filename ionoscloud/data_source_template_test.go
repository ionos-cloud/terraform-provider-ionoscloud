//go:build compute || all || template

package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const templateName = DataSource + "." + TemplateResource + "." + TemplateTestResource

func TestAccDataSourceTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "CUBES S"),
					resource.TestCheckResourceAttr(templateName, "cores", "1"),
					resource.TestCheckResourceAttr(templateName, "ram", "2048"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "50"),
				),
			},
			{
				Config: testAccDataSourceTemplateCores,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "CUBES XL"),
					resource.TestCheckResourceAttr(templateName, "cores", "6"),
					resource.TestCheckResourceAttr(templateName, "ram", "16384"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "320"),
				),
			},
			{
				Config: testAccDataSourceTemplateRam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "CUBES 3XL"),
					resource.TestCheckResourceAttr(templateName, "cores", "12"),
					resource.TestCheckResourceAttr(templateName, "ram", "49152"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "960"),
				),
			},
			{
				Config: testAccDataSourceTemplateStorageSize,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "CUBES M"),
					resource.TestCheckResourceAttr(templateName, "cores", "2"),
					resource.TestCheckResourceAttr(templateName, "ram", "4096"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "80"),
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
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	name = "CUBES S"
}`

const testAccDataSourceTemplateCores = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	cores = 6
}`

const testAccDataSourceTemplateRam = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	ram = 49152
}`

const testAccDataSourceTemplateStorageSize = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	storage_size = 80
}`

const testAccDataSourceTemplateStorageWrongNameError = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	name		 = "CUBES S"
	cores		 = 6
	ram			 = 16384
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongCores = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 1
	ram			 = 16384
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongRam = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 6
	ram			 = 2048
	storage_size = 320
}`

const testAccDataSourceTemplateStorageWrongStorage = `
data ` + TemplateResource + ` ` + TemplateTestResource + ` {
	name		 = "CUBES XL"
	cores		 = 6
	ram			 = 16384
	storage_size = 50
}`
