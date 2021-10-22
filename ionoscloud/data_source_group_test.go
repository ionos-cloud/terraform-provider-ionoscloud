package ionoscloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGroupConfigBasic,
			},
			{
				Config: testAccDataSourceGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "name", GroupResource+"."+GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_datacenter", GroupResource+"."+GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_snapshot", GroupResource+"."+GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "reserve_ip", GroupResource+"."+GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "access_activity_log", GroupResource+"."+GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_pcc", GroupResource+"."+GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "s3_privilege", GroupResource+"."+GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_backup_unit", GroupResource+"."+GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_internet_access", GroupResource+"."+GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceById, "create_k8s_cluster", GroupResource+"."+GroupTestResource, "create_k8s_cluster"),
					testNotEmptySlice(DataSource+"."+GroupResource, "users"),
				),
			},
			{
				Config: testAccDataSourceGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "name", GroupResource+"."+GroupTestResource, "name"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_datacenter", GroupResource+"."+GroupTestResource, "create_datacenter"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_snapshot", GroupResource+"."+GroupTestResource, "create_snapshot"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "reserve_ip", GroupResource+"."+GroupTestResource, "reserve_ip"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "access_activity_log", GroupResource+"."+GroupTestResource, "access_activity_log"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_pcc", GroupResource+"."+GroupTestResource, "create_pcc"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "s3_privilege", GroupResource+"."+GroupTestResource, "s3_privilege"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_backup_unit", GroupResource+"."+GroupTestResource, "create_backup_unit"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_internet_access", GroupResource+"."+GroupTestResource, "create_internet_access"),
					resource.TestCheckResourceAttrPair(DataSource+"."+GroupResource+"."+GroupDataSourceByName, "create_k8s_cluster", GroupResource+"."+GroupTestResource, "create_k8s_cluster"),
					testNotEmptySlice(DataSource+"."+GroupResource, "users"),
				),
			},
		},
	})
}

var testAccDataSourceGroupMatchId = testAccCheckGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceById + ` {
  id			= ` + GroupResource + `.` + GroupTestResource + `.id
}
`

var testAccDataSourceGroupMatchName = testAccCheckGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceByName + ` {
  name			= "` + GroupTestResource + `"
}
`
