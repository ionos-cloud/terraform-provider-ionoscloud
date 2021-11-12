package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceServer(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckServerConfigBasic,
			},
			{
				Config: testAccDataSourceServerMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "name", ServerResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "cores", ServerResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "ram", ServerResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "availability_zone", ServerResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "cpu_family", ServerResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "type", ServerResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.name", ServerResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.size", ServerResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.type", ServerResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.bus", ServerResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "volumes.0.availability_zone", ServerResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.lan", ServerResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.name", ServerResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.dhcp", ServerResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_active", ServerResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_type", ServerResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.ips.0", ServerResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.ips.1", ServerResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.protocol", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.name", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_start", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.port_range_end", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_mac", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.source_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.target_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceById, "nics.0.firewall_rules.0.type", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: testAccDataSourceServerMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "name", ServerResource+"."+ServerTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "cores", ServerResource+"."+ServerTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "ram", ServerResource+"."+ServerTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "availability_zone", ServerResource+"."+ServerTestResource, "availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "cpu_family", ServerResource+"."+ServerTestResource, "cpu_family"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "type", ServerResource+"."+ServerTestResource, "type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.name", ServerResource+"."+ServerTestResource, "volume.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.size", ServerResource+"."+ServerTestResource, "volume.0.size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.type", ServerResource+"."+ServerTestResource, "volume.0.disk_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.bus", ServerResource+"."+ServerTestResource, "volume.0.bus"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "volumes.0.availability_zone", ServerResource+"."+ServerTestResource, "volume.0.availability_zone"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.lan", ServerResource+"."+ServerTestResource, "nic.0.lan"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.name", ServerResource+"."+ServerTestResource, "nic.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.dhcp", ServerResource+"."+ServerTestResource, "nic.0.dhcp"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_active", ServerResource+"."+ServerTestResource, "nic.0.firewall_active"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_type", ServerResource+"."+ServerTestResource, "nic.0.firewall_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.ips.0", ServerResource+"."+ServerTestResource, "nic.0.ips.0"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.ips.1", ServerResource+"."+ServerTestResource, "nic.0.ips.1"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.protocol", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.protocol"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.name", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_start", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_start"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.port_range_end", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.port_range_end"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_mac", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_mac"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.source_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.source_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.target_ip", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.target_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+ServerResource+"."+ServerDataSourceByName, "nics.0.firewall_rules.0.type", ServerResource+"."+ServerTestResource, "nic.0.firewall.0.type"),
				),
			},
			{
				Config: "/* intentionally left blank to ensure proper datasource removal from the plan */",
			},
		},
	})
}

const testAccDataSourceServerMatchId = testAccCheckServerConfigBasic + `
data ` + ServerResource + ` ` + ServerDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + ServerResource + `.` + ServerTestResource + `.id
}
`

const testAccDataSourceServerMatchName = testAccCheckServerConfigBasic + `
data ` + ServerResource + ` ` + ServerDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + ServerTestResource + `"
}
`
