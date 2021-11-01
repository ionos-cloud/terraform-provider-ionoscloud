package ionoscloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceById = "test_ipblock_data"
const dataSourceMatching = "test_ipblock_data_name_loc"
const dataSourceMatchName = "test_ipblock_data_name"
const fullResourceName = IpBLockResource + ".test_ipblock"

var dataSourceNameById = fmt.Sprintf("%s.%s.%s", DataSource, IpBLockResource, dataSourceById)
var dataSourceNameMatching = fmt.Sprintf("%s.%s.%s", DataSource, IpBLockResource, dataSourceMatching)
var dataSourceNameMatchName = fmt.Sprintf("%s.%s.%s", DataSource, IpBLockResource, dataSourceMatchName)

const location = "us/las"
const name = "test_ipblock_name"

func TestAccDataSourceIpBlock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIpBLockResource,
			},
			{
				Config: testAccDataSourceIpBlockMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameById, "name", fullResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "location", fullResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "size", fullResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "ips", fullResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceNameById, "ip_consumers", fullResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatching,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameMatching, "name", fullResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatching, "location", fullResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatching, "size", fullResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatching, "ips", fullResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatching, "ip_consumers", fullResourceName, "ip_consumers"),
				),
			},
			{
				Config: testAccDataSourceIpBlockMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceNameMatchName, "name", fullResourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatchName, "location", fullResourceName, "location"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatchName, "size", fullResourceName, "size"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatchName, "ips", fullResourceName, "ips"),
					resource.TestCheckResourceAttrPair(dataSourceNameMatchName, "ip_consumers", fullResourceName, "ip_consumers"),
				),
			},
			{
				Config:      testAccDataSourceIpBlockNameError,
				ExpectError: regexp.MustCompile(`could not find an ip block with name 1`),
			},
			{
				Config:      testAccDataSourceIpBlockMatchNameLocationError,
				ExpectError: regexp.MustCompile(`there are no ip blocks that match the search criteria`),
			},
			{
				Config:      testAccDataSourceIpBlockLocationError,
				ExpectError: regexp.MustCompile(`there are no ip blocks that match the search criteria`),
			},
			{
				Config:      testIpBlockGoodIdLocationError,
				ExpectError: regexp.MustCompile(`does not match expected location`),
			},
		},
	})
}

const testAccDataSourceIpBLockResource = `resource ` + IpBLockResource + ` "test_ipblock" {
  location = "` + location + `"
  size     = 1
  name     = "test_ipblock_name"
}
`

const testAccDataSourceIpBlockMatchId = testAccDataSourceIpBLockResource +
	"data " + IpBLockResource + " test_ipblock_data " +
	"{ id = " + fullResourceName + ".id }"

const testAccDataSourceIpBlockMatching = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name_loc { 
	name = ` + fullResourceName + `.name
	location = ` + fullResourceName + `.location 
}`

const testAccDataSourceIpBlockMatchName = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name { 
	name = ` + fullResourceName + `.name
}`

const testAccDataSourceIpBlockNameError = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name { 
	name = ` + fullResourceName + `.size
}`
const testAccDataSourceIpBlockMatchNameLocationError = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name { 
	name = ` + fullResourceName + `.name
	location = "none"
}`
const testAccDataSourceIpBlockLocationError = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name {
	location = "none"
}`

const testIpBlockGoodIdLocationError = testAccDataSourceIpBLockResource +
	`data ` + IpBLockResource + ` test_ipblock_data_name {
    id = ` + fullResourceName + `.id
	location = "none"
}`
