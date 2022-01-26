//go:build compute || all || group

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
				Config: testAccDataSourceGroupConfigBasic,
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

var testAccDataSourceGroupConfigBasic = `
resource ` + UserResource + ` ` + UserTestResource + ` {
  first_name = "user"
  last_name = "test"
  email = "` + GenerateEmail() + `"
  password = "abc123-321CBA"
  administrator = false
  force_sec_auth= false
  active = false
}

resource ` + GroupResource + ` ` + GroupTestResource + ` {
  name = "` + GroupTestResource + `"
  create_datacenter = true
  create_snapshot = true
  reserve_ip = true
  access_activity_log = true
  create_pcc = true
  s3_privilege = true
  create_backup_unit = true
  create_internet_access = true
  create_k8s_cluster = true
  user_id = ` + UserResource + `.` + UserTestResource + `.id
}
`

var testAccDataSourceGroupMatchId = testAccDataSourceGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceById + ` {
  id			= ` + GroupResource + `.` + GroupTestResource + `.id
}
`

var testAccDataSourceGroupMatchName = testAccDataSourceGroupConfigBasic + `
data ` + GroupResource + ` ` + GroupDataSourceByName + ` {
  name			= "` + GroupTestResource + `"
}
`
