//go:build all || dsaas
// +build all dsaas

package ionoscloud

import (
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDSaaSVersions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDSaaSVersions,
				Check: resource.ComposeTestCheckFunc(
					utils.TestNotEmptySlice(DataSource+"."+DSaaSVersionsDataSource+"."+DSaaSVersionsTestDataSource, "versions.#"),
				),
			},
		},
	})
}

const testAccDataSourceDSaaSVersions = `
data ` + DSaaSVersionsDataSource + ` ` + DSaaSVersionsTestDataSource + ` {
}
`
