package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccDataSourceNic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCreateDataCenterAndServer,
			},
			{
				Config: testAccDataSourceNicMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
				),
			},
			{
				Config: testAccDataSourceNicMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", fullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", fullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", fullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", fullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", fullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", fullNicResourceName, "ips"),
				),
			},
			{
				Config:      testAccDataSourceNicMatchNameError,
				ExpectError: regexp.MustCompile(`there are no nics that match the search criteria`),
			},
			{
				Config:      testAccDataSourceNicMatchIdAndNameError,
				ExpectError: regexp.MustCompile(`does not match expected name`),
			},
		},
	})
}

const dataSourceNicById = nicResource + ".test_nic_data"

const testAccDataSourceNicMatchId = testAccCheckNicConfigBasic + `
data ` + nicResource + ` test_nic_data {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  id = ` + fullNicResourceName + `.id
}
`

const testAccDataSourceNicMatchName = testAccCheckNicConfigBasic +
	`data ` + nicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = ` + fullNicResourceName + `.name 
}`

const testAccDataSourceNicMatchNameError = testAccCheckNicConfigBasic +
	`data ` + nicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = "DoesNotExist"
}`

const testAccDataSourceNicMatchIdAndNameError = testAccCheckNicConfigBasic +
	`data ` + nicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	id = ` + fullNicResourceName + `.id
	name = "doesNotExist"
}`
