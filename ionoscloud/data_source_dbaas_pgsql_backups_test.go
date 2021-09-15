package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDbaasPgSqlBackups_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlAllBackups,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_dbaas_pgsql_backups.test_ds_dbaas_backups", "cluster_backups.0.cluster_id", "f841e240-149d-11ec-953a-ea197349f9c0"),
				),
			},
		},
	})
}

const testAccDataSourceDbaasPgSqlAllBackups = `
data "ionoscloud_dbaas_pgsql_backups" "test_ds_dbaas_backups" {
	cluster_id = "f841e240-149d-11ec-953a-ea197349f9c0"
}
`
