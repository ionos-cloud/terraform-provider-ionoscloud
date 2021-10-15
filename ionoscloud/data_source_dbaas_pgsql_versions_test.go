package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDbaasPgSqlVersions_All(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDbaasPgSqlAllVersions,
				Check: resource.ComposeTestCheckFunc(
					testNotEmptySlice(DBaaSVersionsResource, "postgres_versions.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasPgSqlVersions_ClusterId(t *testing.T) {
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
				Config: testAccDataSourceDbaasPgSqlVersionsByClusterId,
				Check: resource.ComposeTestCheckFunc(
					testNotEmptySlice(DBaaSVersionsResource, "postgres_versions.#"),
				),
			},
		},
	})

}

const testAccDataSourceDbaasPgSqlAllVersions = `
data ` + DBaaSVersionsResource + ` ` + DBaaSVersionsTest + ` {
}
`

const testAccDataSourceDbaasPgSqlVersionsByClusterId = testAccCheckDbaasPgSqlClusterConfigBasic + `
data ` + DBaaSVersionsResource + ` ` + DBaaSVersionsTest + ` {
	cluster_id = ` + DBaaSClusterResource + `.` + DBaaSClusterTestResource + `.id
}
`
