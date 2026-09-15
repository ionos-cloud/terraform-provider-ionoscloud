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
				Config: testAccDataSourceTemplateCategoryCoresRam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Basic Cube S"),
					resource.TestCheckResourceAttr(templateName, "category", "Basic Templates"),
					resource.TestCheckResourceAttr(templateName, "cores", "2"),
					resource.TestCheckResourceAttr(templateName, "ram", "4096"),
				),
			},
			{
				Config: testAccDataSourceTemplateNameCategory,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "Basic Cube XL"),
					resource.TestCheckResourceAttr(templateName, "category", "Basic Templates"),
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
			{
				Config:      testAccDataSourceTemplateWrongCategory,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateCategoryMismatch,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateNoCriteriaMultipleError,
				ExpectError: regexp.MustCompile(`more than one template found`),
			},
			{
				Config:      testAccDataSourceTemplatePartialNameMultipleError,
				ExpectError: regexp.MustCompile(`more than one template found`),
			},
			{
				Config: testAccDataSourceTemplateNameStorageTypePerformance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "H200-M"),
					resource.TestCheckResourceAttr(templateName, "cores", "30"),
					resource.TestCheckResourceAttr(templateName, "ram", "546816"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "1536"),
					resource.TestCheckResourceAttr(templateName, "storage_type", "PERFORMANCE"),
					resource.TestCheckResourceAttr(templateName, "category", "GPU Category"),
				),
			},
			{
				Config: testAccDataSourceTemplateNameStorageTypeSsdPremium,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "L40S-M"),
					resource.TestCheckResourceAttr(templateName, "cores", "30"),
					resource.TestCheckResourceAttr(templateName, "ram", "546816"),
					resource.TestCheckResourceAttr(templateName, "storage_size", "600"),
					resource.TestCheckResourceAttr(templateName, "storage_type", "SSD Premium"),
					resource.TestCheckResourceAttr(templateName, "category", "AI Model Inference Host"),
				),
			},
			{
				Config: testAccDataSourceTemplateCategoryStorageTypeCoresRam,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(templateName, "name", "H200-M"),
					resource.TestCheckResourceAttr(templateName, "storage_type", "PERFORMANCE"),
					resource.TestCheckResourceAttr(templateName, "category", "GPU Category"),
				),
			},
			{
				Config:      testAccDataSourceTemplateWrongStorageTypeError,
				ExpectError: regexp.MustCompile(`no template found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceTemplateNameStorageTypeMismatchError,
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

const testAccDataSourceTemplateNoCriteriaMultipleError = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
}`

const testAccDataSourceTemplatePartialNameMultipleError = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name = "Basic Cube"
}`

const testAccDataSourceTemplateCategoryCoresRam = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	category = "Basic Templates"
	cores    = 2
	ram      = 4096
}`

const testAccDataSourceTemplateNameCategory = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name     = "Basic Cube XL"
	category = "Basic Templates"
}`

const testAccDataSourceTemplateWrongCategory = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	category = "NON_EXISTENT_CATEGORY"
}`

const testAccDataSourceTemplateCategoryMismatch = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name     = "Basic Cube S"
	category = "NON_EXISTENT_CATEGORY"
}`

// H200-M exists as two GPU templates sharing the same name/cores/ram/storage_size,
// distinguished only by storage_type ("PERFORMANCE" and "SSD Premium").
const testAccDataSourceTemplateNameStorageTypePerformance = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name         = "H200-M"
	storage_type = "PERFORMANCE"
}`

const testAccDataSourceTemplateNameStorageTypeSsdPremium = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name         = "L40S-M"
	storage_type = "SSD Premium"
}`

const testAccDataSourceTemplateCategoryStorageTypeCoresRam = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	category     = "GPU Category"
	storage_type = "PERFORMANCE"
	cores        = 30
	ram          = 546816
}`

const testAccDataSourceTemplateWrongStorageTypeError = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	storage_type = "NON_EXISTENT_STORAGE_TYPE"
}`

// "Basic Cube S" has no storage_type set at all, so filtering it by any storage_type
// value returns zero results.
const testAccDataSourceTemplateNameStorageTypeMismatchError = `
data ` + constant.TemplateResource + ` ` + constant.TemplateTestResource + ` {
	name         = "Basic Cube S"
	storage_type = "HDD"
}`
