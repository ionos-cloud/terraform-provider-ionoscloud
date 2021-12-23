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
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", FullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", FullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", FullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", FullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", FullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", FullNicResourceName, "ips"),
				),
			},
			{
				Config: testAccDataSourceNicMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "name", FullNicResourceName, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "dhcp", FullNicResourceName, "dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "firewall_active", FullNicResourceName, "firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "mac", FullNicResourceName, "mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "lan", FullNicResourceName, "lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+dataSourceNicById, "ips", FullNicResourceName, "ips"),
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

const dataSourceNicById = NicResource + ".test_nic_data"

const testAccDataSourceNicMatchId = testAccCheckNicConfigBasic + `
data ` + NicResource + ` test_nic_data {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  id = ` + FullNicResourceName + `.id
}
`

const testAccDataSourceNicMatchName = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = ` + FullNicResourceName + `.name 
}`

const testAccDataSourceNicMatchNameError = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	name = "DoesNotExist"
}`

const testAccDataSourceNicMatchIdAndNameError = testAccCheckNicConfigBasic +
	`data ` + NicResource + ` test_nic_data {
  	datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
	server_id = ` + ServerResource + `.` + ServerTestResource + `.id
	id = ` + FullNicResourceName + `.id
	name = "doesNotExist"
}`
