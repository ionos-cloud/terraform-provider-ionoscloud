package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDbaasPgSqlClusterBackups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlClusterBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(DataSource+"."+DBaaSBackupsResource+"."+DBaaSBackupsTest, "cluster_backups.0.cluster_id", DataSource+"."+DBaaSBackupsResource+"."+DBaaSBackupsTest, "cluster_id"),
					testNotEmptySlice(DBaaSBackupsResource, "cluster_backups.#"),
				),
			},
		},
	})
}

const testAccDataSourceDbaasPgSqlClusterBackups = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSBackupsResource + ` ` + DBaaSBackupsTest + ` {
	cluster_id = ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`
