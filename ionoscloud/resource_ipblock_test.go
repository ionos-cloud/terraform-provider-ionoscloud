//go:build compute || all || ipblock

package ionoscloud

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/cloud/v2"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/constant"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const fullIpBlockResourceName = constant.IpBlockResource + "." + constant.IpBlockTestResource

var dataSourceIpBlockNameById = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceById)
var dataSourceIpBlockNameMatching = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceMatching)
var dataSourceIpBlockNameMatchName = fmt.Sprintf("%s.%s.%s", constant.DataSource, constant.IpBlockResource, constant.IpBlockDataSourceByName)

const location = "us/las"

func TestAccIPBlockBasic(t *testing.T) {
	var ipblock ionoscloud.IpBlock

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesInternal(t, &testAccProvider),
		CheckDestroy:             testAccCheckIPBlockDestroyCheck,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIPBlockConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(fullIpBlockResourceName, &ipblock),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "location", location),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "name", constant.IpBlockTestResource),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "size", "1"),
				),
			}, {
				Config: testAccDataSourceIpBlockMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameById, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameById, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameById, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameById, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameById, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatching, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "name", fullIpBlockResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "location", fullIpBlockResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "size", fullIpBlockResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "ips", fullIpBlockResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceIpBlockNameMatchName, "ip_consumers", fullIpBlockResourceName, "ip_consumers"),
				),
			},
			{
				Config:      testAccDataSourceIpBlockMultipleResultsError,
				ExpectError: regexp.MustCompile(`more than one ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockNameError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockMatchNameLocationError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockLocationError,
				ExpectError: regexp.MustCompile(`no ip block found with the specified criteria`),
			},
			{
				Config:      testIpBlockGoodIdLocationError,
				ExpectError: regexp.MustCompile(`location of ip block`),
			},
			{
				Config: testAccCheckIPBlockConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPBlockExists(fullIpBlockResourceName, &ipblock),
					testAccCheckIPBlockAttributes(fullIpBlockResourceName, location),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "name", constant.UpdatedResources),
					resource.TestCheckResourceAttr(fullIpBlockResourceName, "size", "2"),
				),
			},
		},
	})
}

func testAccCheckIPBlockDestroyCheck(s *terraform.State) error {
	client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != constant.IpBlockResource {
			continue
		}

		_, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			if !httpNotFound(apiResponse) {
				return fmt.Errorf("an error occurred while checking deletion of IPBlock %s %s", rs.Primary.ID, err)
			}
		} else {
			return fmt.Errorf("IPBlock still exists %s %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIPBlockAttributes(n string, location string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckLanAttributes: Not found: %s", n)
		}
		if rs.Primary.Attributes["location"] != location {
			return fmt.Errorf("bad name: %s", rs.Primary.Attributes["location"])
		}

		return nil
	}
}

func testAccCheckIPBlockExists(n string, ipblock *ionoscloud.IpBlock) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(services.SdkBundle).CloudApiClient

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("testAccCheckIPBlockExists: Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
		if cancel != nil {
			defer cancel()
		}
		foundIP, apiResponse, err := client.IPBlocksApi.IpblocksFindById(ctx, rs.Primary.ID).Execute()
		logApiRequestTime(apiResponse)

		if err != nil {
			return fmt.Errorf("error occurred while fetching IP Block: %s", rs.Primary.ID)
		}
		if *foundIP.Id != rs.Primary.ID {
			return fmt.Errorf("record not found")
		}

		ipblock = &foundIP

		return nil
	}
}

const testAccCheckIPBlockConfigBasic = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + ` {
  location = "` + location + `"
  size = 1
  name = "` + constant.IpBlockTestResource + `"
}`

const testAccCheckIPBlockConfigUpdate = `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + `{
  location = "` + location + `"
  size = 2
  name = "` + constant.UpdatedResources + `"
}`

const testAccDataSourceIpBlockMatchId = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + `  ` + constant.IpBlockDataSourceById + ` {
	id = ` + fullIpBlockResourceName + `.id 
}
`

const testAccDataSourceIpBlockMatching = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceMatching + ` { 
	name = ` + fullIpBlockResourceName + `.name
	location = ` + fullIpBlockResourceName + `.location 
}`

const testAccDataSourceIpBlockMatchName = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
}`

const testAccDataSourceIpBlockMultipleResultsError = testAccCheckIPBlockConfigBasic + `
resource ` + constant.IpBlockResource + ` ` + constant.IpBlockTestResource + `_same_name{
  location = "` + location + `"
  size = 2
  name = ` + fullIpBlockResourceName + `.name
}

data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
}`

const testAccDataSourceIpBlockNameError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.size
}`
const testAccDataSourceIpBlockMatchNameLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` { 
	name = ` + fullIpBlockResourceName + `.name
	location = "none"
}`
const testAccDataSourceIpBlockLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
	location = "none"
}`

const testIpBlockGoodIdLocationError = testAccCheckIPBlockConfigBasic + `
data ` + constant.IpBlockResource + ` ` + constant.IpBlockDataSourceByName + ` {
    id = ` + fullIpBlockResourceName + `.id
	location = "none"
}`
