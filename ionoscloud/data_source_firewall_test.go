package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFirewall(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckFirewallConfigBasic,
			},
			{
				Config: testAccDataSourceFirewallMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "name", FirewallResource+"."+FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "protocol", FirewallResource+"."+FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "source_mac", FirewallResource+"."+FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "source_ip", FirewallResource+"."+FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "target_ip", FirewallResource+"."+FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "icmp_type", FirewallResource+"."+FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceById, "icmp_code", FirewallResource+"."+FirewallTestResource, "icmp_code"),
				),
			},
			{
				Config: testAccDataSourceFirewallMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "name", FirewallResource+"."+FirewallTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "protocol", FirewallResource+"."+FirewallTestResource, "protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "source_mac", FirewallResource+"."+FirewallTestResource, "source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "source_ip", FirewallResource+"."+FirewallTestResource, "source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "target_ip", FirewallResource+"."+FirewallTestResource, "target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "icmp_type", FirewallResource+"."+FirewallTestResource, "icmp_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+FirewallResource+"."+FirewallDataSourceByName, "icmp_code", FirewallResource+"."+FirewallTestResource, "icmp_code"),
				),
			},
		},
	})
}

const testAccDataSourceFirewallMatchId = testAccCheckFirewallConfigBasic + `
data ` + FirewallResource + ` ` + FirewallDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  id = ` + FirewallResource + `.` + FirewallTestResource + `.id
}
`

const testAccDataSourceFirewallMatchName = testAccCheckFirewallConfigBasic + `
data ` + FirewallResource + ` ` + FirewallDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  server_id = ` + ServerResource + `.` + ServerTestResource + `.id
  nic_id = ionoscloud_nic.database_nic.id
  name	= "` + FirewallTestResource + `"
}
`
