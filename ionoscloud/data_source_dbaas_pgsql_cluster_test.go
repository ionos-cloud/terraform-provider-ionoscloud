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
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "replicas", DBaaSClusterResource+"."+DBaaSClusterTestResource, "replicas"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "cpu_core_count", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cpu_core_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "ram_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "vdc_connections.vdc_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.vdc_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "vdc_connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "vdc_connections.ip_address", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.ip_address"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.weekday", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.weekday"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "maintenance_window.time", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.time"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.username", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.username"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceById, "credentials.password", DBaaSClusterResource+"."+DBaaSClusterTestResource, "credentials.password"),
				),
			},
			{
				Config: testAccDataSourceDBaaSPgSqlClusterMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "replicas", DBaaSClusterResource+"."+DBaaSClusterTestResource, "replicas"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "cpu_core_count", DBaaSClusterResource+"."+DBaaSClusterTestResource, "cpu_core_count"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "ram_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "ram_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_size", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_size"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "storage_type", DBaaSClusterResource+"."+DBaaSClusterTestResource, "storage_type"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "vdc_connections.vdc_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.vdc_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "vdc_connections.lan_id", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.lan_id"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "vdc_connections.ip_address", DBaaSClusterResource+"."+DBaaSClusterTestResource, "vdc_connections.ip_address"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "location", DBaaSClusterResource+"."+DBaaSClusterTestResource, "location"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "display_name", DBaaSClusterResource+"."+DBaaSClusterTestResource, "display_name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSClusterResource+"."+DBaaSClusterTestDataSourceByName, "maintenance_window.weekday", DBaaSClusterResource+"."+DBaaSClusterTestResource, "maintenance_window.weekday"),
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
