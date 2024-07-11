//go:build all || dbaas || mongo

package ionoscloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"
)

func TestAccDataSourceDBaaSMongoTemplate(t *testing.T) {
	checkFunction := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(dataSourceAccess, "id", expectedId),
		resource.TestCheckResourceAttr(dataSourceAccess, "name", expectedName),
		resource.TestCheckResourceAttr(dataSourceAccess, "edition", expectedEdition),
		resource.TestCheckResourceAttr(dataSourceAccess, "cores", expectedCores),
		resource.TestCheckResourceAttr(dataSourceAccess, "ram", expectedRam),
		resource.TestCheckResourceAttr(dataSourceAccess, "storage_size", expectedStorageSize),
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ExternalProviders:        randomProviderVersion343(),
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactoriesInternal(t, &testAccProvider),
		Steps: []resource.TestStep{
			{
				Config:      invalidMissingBothIdAndName,
				ExpectError: regexp.MustCompile("please provide a template ID or name"),
			},
			{
				Config:      invalidProvidingBothIdAndName,
				ExpectError: regexp.MustCompile("name and ID cannot be both specified at the same time"),
			},
			{
				Config:      invalidGetByIdNonExistentTemplate,
				ExpectError: regexp.MustCompile("no DBaaS Mongo Template found with the specified criteria"),
			},
			{
				Config:      invalidGetByNameNonExistentTemplate,
				ExpectError: regexp.MustCompile("no DBaaS Mongo Template found with the specified criteria"),
			},
			{
				Config:      invalidGetByNameMultipleTemplates,
				ExpectError: regexp.MustCompile("more than one DBaaS Mongo Template found for the specified search criteria"),
			},
			{
				Config: validGetById,
				Check:  checkFunction,
			},
			{
				Config: validGetByName,
				Check:  checkFunction,
			},
			{
				Config: validGetByNamePartialMatch,
				Check:  checkFunction,
			},
		},
	})
}

const invalidMissingBothIdAndName = `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
}`

// We are looking for an UUID that doesn't exist. Usually, simply generating a value wouldn't be
// enough because we would have to check that the value doesn't really exist, but given the fact
// that the UUID consists of several characters, the probability of generating an UUID that already
// exists is very small, so we don't need any additional check.
const invalidProvidingBothIdAndName = resourceRandomUUID + resourceRandomString + `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	id = random_uuid.uuid.result
	name = random_string.simple_string.result
}`

const invalidGetByIdNonExistentTemplate = resourceRandomUUID + `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	id = random_uuid.uuid.result
}`

const invalidGetByNameNonExistentTemplate = resourceRandomString + `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	name = random_string.simple_string.result
}`

const invalidGetByNameMultipleTemplates = `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	name = "MongoDB"
	partial_match = true
}`

const dataSourceAccess = constant.DataSource + "." + constant.DBaaSMongoTemplateResource + "." + constant.DBaaSMongoTemplateTestDataSource

const validGetById = `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	id = "ea320e28-b973-457a-86c5-68c19dd06d3d"
}`

const validGetByName = `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	name = "MongoDB Business 4XL_S"
}`

const validGetByNamePartialMatch = `
data ` + constant.DBaaSMongoTemplateResource + ` ` + constant.DBaaSMongoTemplateTestDataSource + ` {
	name = "Business 4XL_S"
	partial_match = true
}`

// We are testing using the same template, so we can just define the expected values as constants
// and re-use them.
const expectedId = "ea320e28-b973-457a-86c5-68c19dd06d3d"
const expectedName = "MongoDB Business 4XL_S"
const expectedEdition = "business"
const expectedCores = "32"
const expectedRam = "131072"
const expectedStorageSize = "2048"
