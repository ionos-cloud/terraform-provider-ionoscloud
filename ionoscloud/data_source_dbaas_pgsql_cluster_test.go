//go:build all || dbaas
// +build all dbaas

package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDBaaSPgSqlCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDbaasPgSqlClusterConfigBasic,
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "instances", DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "cores", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "ram", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.datacenter_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "connections.cidr", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.day_of_the_week", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.time", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.username", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.password", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "instances", DBaaSClusterResource+"."+DBaaSClusterTestResource, "instances"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "cores", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cores"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "ram", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.datacenter_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.datacenter_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "connections.cidr", DBaaSClusterResource+"."+DBaaSClusterTestResource, "connections.cidr"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.day_of_the_week", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.day_of_the_week"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.time", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.username", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "credentials.password", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
		},
	})
}

const testAccDataSourceDBaaSPgSqlClusterMatchId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSClusterResource + ` ` + DBaaSClusterTestDataSourceById + ` {
  id	= ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`

const testAccDataSourceDBaaSPgSqlClusterMatchName = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSClusterResource + ` ` + DBaaSClusterTestDataSourceByName + ` {
  display_name	= "` + DBaaSClusterTestResource + `"
}
`
