package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceLan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLanConfigBasic,
			},
			{
				Config: testAccDataSourceLanMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "name", LanResource+"."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "ip_failover.nic_uuid", LanResource+"."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "ip_failover.ip", LanResource+"."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "pcc", LanResource+"."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceById, "public", LanResource+"."+LanTestResource, "public"),
				),
			},
			{
				Config: testAccDataSourceLanMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "name", LanResource+"."+LanTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "ip_failover.nic_uuid", LanResource+"."+LanTestResource, "ip_failover.nic_uuid"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "ip_failover.ip", LanResource+"."+LanTestResource, "ip_failover.ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "pcc", LanResource+"."+LanTestResource, "pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+LanResource+"."+LanDataSourceByName, "public", LanResource+"."+LanTestResource, "public"),
				),
			},
		},
	})
}

const testAccDataSourceLanMatchId = testAccCheckLanConfigBasic + `
data ` + LanResource + ` ` + LanDataSourceById + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  id			= ` + LanResource + `.` + LanTestResource + `.id
}
`

const testAccDataSourceLanMatchName = testAccCheckLanConfigBasic + `
data ` + LanResource + ` ` + LanDataSourceByName + ` {
  datacenter_id = ` + DatacenterResource + `.` + DatacenterTestResource + `.id
  name			= "` + LanTestResource + `"
}
`
