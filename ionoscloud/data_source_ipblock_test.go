package ionoscloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const dataSourceById = "test_ipblock_data"
const dataSourceMatching = "test_ipblock_data_name"
const fullResourceName = IpBLockResource + ".test_ipblock"

var dataSourceNameById = fmt.Sprintf("%s.%s.%s", DataSource, IpBLockResource, dataSourceById)
var dataSourceNameMatching = fmt.Sprintf("%s.%s.%s", DataSource, IpBLockResource, dataSourceMatching)

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
	`data ` + IpBLockResource + ` test_ipblock_data_name { 
	name = ` + fullResourceName + `.name
	location = ` + fullResourceName + `.location 
}`
